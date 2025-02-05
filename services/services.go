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
		sendBotMessage(bot, update.Message.Chat.ID, text, keyboards.MainKeyboard())
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

		sendBotMessage(bot, update.Message.Chat.ID, text, nil)

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

	sendBotMessage(bot, chatID, "Ваш игровой ID успешно сохранен!", keyboards.MainKeyboard())

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

	sendBotMessage(bot, update.Message.Chat.ID, text, keyboards.ChangeGameIDKeyboard())
}

func ChangeGameIDCallback(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Введите новый игровой ID"
	sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, tgbotapi.NewRemoveKeyboard(true))

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

	По всем вопросам обращаться к админу - @admin
	`
	sendBotMessage(bot, update.Message.Chat.ID, text, keyboards.HomeKeyboard())
}

func Home(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var chatID int64
	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	} else {
		chatID = update.Message.Chat.ID
	}
	sendBotMessage(bot, chatID, "Вы вернулись в главное меню", keyboards.MainKeyboard())
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

	sendBotMessage(bot, update.Message.Chat.ID, text, keyboards.WithdrawBonusKeyboard())
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
		sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, nil)
		return
	}

	text = "Укажите игровое ID, получателя. Вы можете отправить свое id нажав на кнопку или же ввести самостоятельно id получателя"
	sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, keyboards.SendGameIDKeyboard(user.GameId))

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

	sendBotMessage(bot, adminChatID, text, keyboards.ConfirmBonusKeyboard(donationID))

	text = fmt.Sprintf("Скоро админ подтвердит получение бонусных алмазов. Ожидайте. Количество: %v", referredDiamonds)
	sendBotMessage(bot, update.Message.Chat.ID, text, nil)

	SetUserState(update.Message.Chat.ID, "")
}

func ConfirmedDonateBonus(bot *tgbotapi.BotAPI, update tgbotapi.Update, donationID string) {
	text := fmt.Sprintf("Вы подтвердили Donate ID: %s", donationID)
	sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, nil)

	donation, err := repositories.GetDataDonation(donationID)
	if err != nil {
		log.Print(err)
	}
	if err := repositories.ResetReferredDiamonds(donation.TelegramID); err != nil {
		log.Print(err)
	}

	text = "Админ подтвердил вашу заявку! Алмазы скоро будут у вас на аккаунте!"

	sendBotMessage(bot, donation.TelegramID, text, nil)

	if err := repositories.SetDonateStatusCompleted(donationID); err != nil {
		log.Print(err)
	}
}

func ChooseServerID(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Выберите количество цифр в ZONE ID"

	sendBotMessage(bot, update.Message.Chat.ID, text, keyboards.ChooseIDKeyboard())
}

func Catalog(bot *tgbotapi.BotAPI, update tgbotapi.Update, zone int) {
	var text string
	if zone == 4 {
		text = "Каталог товаров для ZONE ID 4. \nВыберите что хотите приобрести ниже"
	} else {
		text = "Каталог товаров для ZONE ID 5. \nВыберите что хотите приобрести ниже"
	}

	sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, keyboards.CatalogKeyboard(zone))
}

func ProcessPurchaseDiamonds(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	cart := GetCart(update.CallbackQuery.Message.Chat.ID)
	if !CartExists(cart) {
		sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, "Ваша корзина пуста", nil)
		return
	}
	user, err := repositories.GetUserData(update.CallbackQuery.Message.Chat.ID)
	if err != nil {
		log.Print(err)
	}

	text := "Укажите игровое ID, получателя. Вы можете отправить свое id нажав на кнопку или же ввести самостоятельно id получателя"
	sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, keyboards.SendGameIDKeyboard(user.GameId))

	SetUserState(update.CallbackQuery.Message.Chat.ID, "awaiting_donate_id")
}

func ConfirmDonate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	SetUserGameID(update.Message.Chat.ID, "", update.Message.Text)
	finalPrice, _ := FinalPriceAndDiamonds(update.Message.Chat.ID)

	text := fmt.Sprintf("ID получателя: %s\nРеквизиты для оплаты:\nСумма: %d сом\nМбанк:+996 000 000 000\nПолучатель: Виоле\nПосле оплаты обязательно нажмите кнопку ниже для подтверждения\nПомощь - /help",
		update.Message.Text, finalPrice,
	)
	sendBotMessage(bot, update.Message.Chat.ID, text, keyboards.ConfirmPurchaseKeyboard())
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
	sendBotMessage(bot, adminChatID, text, keyboards.ConfirmKeyboard(donationID))

	sendBotMessage(bot, chatID, "Спасибо за покупку! Скоро админ подтвердит оплату и алмазы будут у вас. Ожидайте.", keyboards.MainKeyboard())
}

func ConfirmedDonate(bot *tgbotapi.BotAPI, update tgbotapi.Update, donationID string) {
	data, err := repositories.GetDataDonation(donationID)
	if err != nil {
		fmt.Print(err)
	}

	switch data.Status {
	case "failed":
		text := fmt.Sprintf("Вы уже отменили Donate ID: %s", donationID)
		sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, nil)
		return
	case "completed":
		text := fmt.Sprintf("Вы уже подтвердили Donate ID: %s", donationID)
		sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, nil)
		return
	}

	ClearCart(bot, update)

	text := fmt.Sprintf("Вы подтвердили Donate ID: %s", donationID)
	sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, nil)

	donation, err := repositories.GetDataDonation(donationID)
	if err != nil {
		log.Print(err)
	}

	text = "Админ подтвердил вашу оплату! Алмазы скоро будут у вас на аккаунте! Оставьте отзыв!"
	sendBotMessage(bot, donation.TelegramID, text, keyboards.ReviewsKeyboard())

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
		sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, nil)
		return
	case "completed":
		text := fmt.Sprintf("Вы уже подтвердили Donate ID: %s", donationID)
		sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, nil)
		return
	}
	text := fmt.Sprintf("Вы отменили Donate ID: %s", donationID)
	sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, nil)

	donation, err := repositories.GetDataDonation(donationID)
	if err != nil {
		log.Print(err)
	}

	text = "Админ отклонил вашу оплату! Если произошла ошибка пишите админу @admin\n"

	sendBotMessage(bot, donation.TelegramID, text, keyboards.MainKeyboard())

	if err := repositories.SetDonateStatusCanceled(donationID); err != nil {
		log.Print(err)
	}
}

func Reviews(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Вы можете посмотреть отзывы нажав кнопку ниже либо перейдя по ссылке t.me/donaterich"
	sendBotMessage(bot, update.Message.Chat.ID, text, keyboards.ReviewsKeyboard())
}

func UnknownCommand(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Я не знаю этой команды :(. Введите /help, чтобы узнать команды и получить помощь)"
	sendBotMessage(bot, update.Message.Chat.ID, text, nil)
}

func sendMessageGetError(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Не удалось получить данные о пользователе"
	sendBotMessage(bot, update.Message.Chat.ID, text, nil)
}

func sendMessageCreateDonateError(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	text := "Не удалось сохранить запись о донате"
	sendBotMessage(bot, update.CallbackQuery.Message.Chat.ID, text, nil)
}

func sendBotMessage(bot *tgbotapi.BotAPI, chatID int64, text string, replyMarkup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = replyMarkup
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}
