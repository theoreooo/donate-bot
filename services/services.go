package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"donate-bot/keyboards"
	"donate-bot/repositories"
	"donate-bot/utils.go"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var referrerID *int64
	args := strings.Split(update.Message.Text, " ")
	if len(args) > 1 {
		refID, err := strconv.ParseInt(args[1], 10, 64)
		if err == nil {
			referrerID = &refID
		}
	}

	exists, err := repositories.IsUserExists(update.Message.Chat.ID)
	if err != nil {
		log.Printf("Error checking user existense: %v", err)
		return
	}

	if exists {
		text := fmt.Sprintf("С возвращением, %s!", update.Message.From.FirstName)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = keyboards.MainKeyboard()
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
	} else {
		err := repositories.CreateUser(update, referrerID)
		if err != nil {
			log.Print(err)
			return
		}

		var text string
		if referrerID != nil {
			text = fmt.Sprintf("Привет, %s! Здесь ты можешь купить алмазы. Вы были приглашены по реферальной ссылке. \nВведите ваш игровой ID и сервер", update.Message.From.FirstName)
		} else {
			text = fmt.Sprintf("Привет, %s! Здесь ты можешь купить алмазы. \nВведите ваш игровой ID и сервер", update.Message.From.FirstName)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
		log.Print("set")
		SetUserState(update.Message.Chat.ID, "awaiting_game_id")
		log.Print("seted")
	}
}

func SetGameId(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	if err := repositories.SetGameId(chatID, update.Message.Text); err != nil {
		log.Print(err)
	}

	msg := tgbotapi.NewMessage(chatID, "Ваш игровой ID успешно сохранен!")
	msg.ReplyMarkup = keyboards.MainKeyboard()
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}

	SetUserState(chatID, "")
}

func Profile(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var text string
	user, err := repositories.GetUserData(update.Message.Chat.ID)
	if err != nil {
		sendMessageGetError(bot, update)
		log.Print(err)
		return
	}
	text = fmt.Sprintf("Профиль: \nID: %v \nДата регистрации: %v\nЗадоначено алмазов: %v\nАлмазы полученные по реферальной системе: %v\n",
		user.GameId, user.RegistrationDate.Format("02.01.2006"), user.TotalDiamonds, user.TotalReferredDiamonds,
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.ChangeGameIDKeyboard()

	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func ChangeGameIDCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Введите новый игровой ID"

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}

	SetUserState(update.CallbackQuery.Message.Chat.ID, "awaiting_game_id")
}

func Help(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := `
	Доступные команды:
	- /help - показать эту справку
	- Профиль - посмотреть информацию о вашем профиле
	- Каталог - посмотреть доступные товары
	- Партнерская программа - узнать о партнерской программе
	- /home - открыть клавиатуру с командами

	По всем вопросам обращаться к админу - @smog_kotoryi_smog
	`
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func Home(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var chatID int64
	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	} else {
		chatID = update.Message.Chat.ID
	}
	msg := tgbotapi.NewMessage(chatID, "Вы вернулись в главное меню")
	msg.ReplyMarkup = keyboards.MainKeyboard()

	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func Partnership(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var text string
	user, err := repositories.GetUserData(update.Message.Chat.ID)
	if err != nil {
		sendMessageGetError(bot, update)
		log.Print(err)
		return
	}

	referralLink := fmt.Sprintf("https://t.me/luckaanbot?start=%v", user.TelegramID)

	referralCount, err := repositories.GetReferralCount(user.TelegramID)
	if err != nil {
		log.Print(err)
	} else {
		text = fmt.Sprintf("Ваша реферальная ссылка: %s \nКоличество рефералов: %d \nДоступно для вывода: %v", referralLink, referralCount, user.ReferredDiamonds)
		text += "\nРекомендуйте наш сервис своим друзьям и зарабатывайте алмазы! Каждый раз, когда ваши друзья пополняют баланс, вы получаете 1% от их пополнений."
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.WithdrawBonusKeyboard()
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func WithdrawBonusCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var text string
	user, err := repositories.GetUserData(update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		sendMessageGetError(bot, update)
		log.Print(err)
		return
	}

	if user.ReferredDiamonds < 50 {
		text = "Минимальная сумма вывода реферальных бонусов — 50 алмазов."
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
		return
	}
	text = "Укажите игровое ID, получателя. Вы можете отправить свое id нажав на кнопку или же ввести самостоятельно id получателя"

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.SendGameIDKeyboard(user.GameId)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}

	SetUserState(update.CallbackQuery.Message.Chat.ID, "awaiting_bonus_donate_id")
}

func ConfirmBonusDonate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	referredDiamonds, err := repositories.GetReferredDiamonds(update.Message.Chat.ID)
	if err != nil {
		log.Print(err)
	}

	donationID, err := repositories.CreateDonateBonus(update, int64(referredDiamonds))
	if err != nil {
		log.Print(err)
		sendMessageCreateDonateError(bot, update)
		return
	}

	text := fmt.Sprintf("Donate ID(номер заказа):%s \nПользователь %s запросил вывод бонусных алмазов.\nКоличество %v.\nID: %s",
		donationID, update.Message.From.UserName, referredDiamonds, update.Message.Text,
	)

	adminChatID := utils.GetAdminChatID()

	msg := tgbotapi.NewMessage(adminChatID, text)
	msg.ReplyMarkup = keyboards.ConfirmBonusKeyboard(donationID)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Скоро админ подтвердит получение бонусных алмазов. Ожидайте. Количество: %v", referredDiamonds))
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
	SetUserState(update.Message.Chat.ID, "")
}

func ConfirmedDonateBonus(bot *tgbotapi.BotAPI, update tgbotapi.Update, donationID string) {
	text := fmt.Sprintf("Вы подтвердили Donate ID: %s", donationID)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
	donation, err := repositories.GetDataDonation(donationID)
	if err != nil {
		log.Print(err)
	}
	if err := repositories.ResetReferredDiamonds(donation.TelegramID); err != nil {
		log.Print(err)
	}

	text = "Админ подтвердил вашу заявку! Алмазы скоро будут у вас на аккаунте!"

	msg = tgbotapi.NewMessage(donation.TelegramID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
	if err := repositories.SetDonateStatusCompleted(donationID); err != nil {
		log.Print(err)
	}
}

func Catalog(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Каталог товаров. \nВыберите что хотите приобрести ниже"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.CatalogKeyboard()

	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func ProcessPurchaseDiamonds(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	cart := GetCart(update.CallbackQuery.Message.Chat.ID)
	if !CartExists(cart) {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Ваша корзина пуста")
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
		return
	}
	text := "Укажите игровое ID, получателя. Вы можете отправить свое id нажав на кнопку или же ввести самостоятельно id получателя"
	user, err := repositories.GetUserData(update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		log.Print(err)
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.SendGameIDKeyboard(user.GameId)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}

	SetUserState(update.CallbackQuery.Message.Chat.ID, "awaiting_donate_id")
}

func ConfirmDonate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	SetUserGameID(update.Message.Chat.ID, "", update.Message.Text)
	finalPrice, _ := FinalPriceAndDiamonds(update.Message.Chat.ID)

	text := fmt.Sprintf("ID получателя: %s\nРеквизиты для оплаты:\nСумма: %d сом\nМбанк:+996 000 000 000\nПолучатель: Виоле\nПосле оплаты обязательно нажмите кнопку ниже для подтверждения",
		update.Message.Text, finalPrice,
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.ConfirmPurchaseKeyboard()
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func ConfirmedPurchase(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.CallbackQuery.Message.Chat.ID

	user, err := repositories.GetUserData(chatID)
	if err != nil {
		log.Print(err)
	}

	_, finalDiamonds := FinalPriceAndDiamonds(chatID)
	cart := GetCart(chatID)
	cartText := FormatCartAdmin(cart)

	donationID, err := repositories.CreateDonate(update, int64(finalDiamonds))
	if err != nil {
		log.Print(err)
	}

	gameID := GetUserGameID(chatID)
	text := fmt.Sprintf("Donate ID:%s \nПользователь @%s подтвердил оплату заказа: %s.\nИгровое ID: %s",
		donationID, user.TelegramUsername, cartText, gameID,
	)

	adminChatID := utils.GetAdminChatID()

	msg := tgbotapi.NewMessage(adminChatID, text)
	msg.ReplyMarkup = keyboards.ConfirmKeyboard(donationID)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}

	msg = tgbotapi.NewMessage(chatID, "Спасибо за покупку! Скоро админ подтвердит оплату и алмазы будут у вас. Ожидайте.")
	msg.ReplyMarkup = keyboards.MainKeyboard()
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func ConfirmedDonate(bot *tgbotapi.BotAPI, update tgbotapi.Update, donationID string) {
	data, err := repositories.GetDataDonation(donationID)
	if err != nil {
		fmt.Print(err)
	}

	switch data.Status {
	case "failed":
		text := fmt.Sprintf("Вы уже отменили Donate ID: %s", donationID)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
		return
	case "completed":
		text := fmt.Sprintf("Вы уже подтвердили Donate ID: %s", donationID)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
		return
	}

	ClearCart(bot, update)

	text := fmt.Sprintf("Вы подтвердили Donate ID: %s", donationID)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
	donation, err := repositories.GetDataDonation(donationID)
	if err != nil {
		log.Print(err)
	}

	text = "Админ подтвердил вашу оплату! Алмазы скоро будут у вас на аккаунте! Оставьте отзыв!"
	msg = tgbotapi.NewMessage(donation.TelegramID, text)
	msg.ReplyMarkup = keyboards.ReviewsKeyboard()
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
	if err := repositories.SetDonateStatusCompleted(donationID); err != nil {
		log.Print(err)
	}
}

func CanceledDonate(bot *tgbotapi.BotAPI, update tgbotapi.Update, donationID string) {
	data, err := repositories.GetDataDonation(donationID)
	if err != nil {
		fmt.Print(err)
	}
	switch data.Status {
	case "failed":
		text := fmt.Sprintf("Вы уже отменили Donate ID: %s", donationID)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
		return
	case "completed":
		text := fmt.Sprintf("Вы уже подтвердили Donate ID: %s", donationID)
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
		if _, err := bot.Send(msg); err != nil {
			log.Print(err)
		}
		return
	}
	text := fmt.Sprintf("Вы отменили Donate ID: %s", donationID)
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
	donation, err := repositories.GetDataDonation(donationID)
	if err != nil {
		log.Print(err)
	}

	text = "Админ отклонил вашу оплату! Если произошла ошибка пишите админу @smog_kotoryi_smog\n"

	msg = tgbotapi.NewMessage(donation.TelegramID, text)
	msg.ReplyMarkup = keyboards.MainKeyboard()
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
	if err := repositories.SetDonateStatusCanceled(donationID); err != nil {
		log.Print(err)
	}
}

func Reviews(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Вы можете посмотреть отзывы нажав кнопку ниже либо перейдя по ссылке t.me/donaterich"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.ReviewsKeyboard()
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}

}

func UnknownCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Я не знаю этой команды :(. Введите /help, чтобы узнать команды и получить помощь)"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func sendMessageGetError(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Не удалось получить данные о пользователе"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func sendMessageCreateDonateError(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Не удалось сохранить запись о донате"
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}
