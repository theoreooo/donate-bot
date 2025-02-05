package keyboards

import (
	"donate-bot/config"

	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MainKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var cmdKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Каталог"),
			tgbotapi.NewKeyboardButton("Корзина"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Партнерская программа"),
			tgbotapi.NewKeyboardButton("Профиль"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отзывы"),
			tgbotapi.NewKeyboardButton("Помощь"),
		),
	)
	return cmdKeyboard
}

func ChangeGameIDKeyboard() tgbotapi.InlineKeyboardMarkup {
	var gameID = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить игровой ID", "change_game_id"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "home"),
		),
	)
	return gameID
}

func WithdrawBonusKeyboard() tgbotapi.InlineKeyboardMarkup {
	var gameID = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вывести бонусные алмазы", "withdraw_bonus"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "home"),
		),
	)
	return gameID
}

func ConfirmBonusKeyboard(donationID string) tgbotapi.InlineKeyboardMarkup {
	var gameID = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить", fmt.Sprintf("confirm_bonus:%s", donationID)),
		),
	)
	return gameID
}

func SendGameIDKeyboard(gameID string) tgbotapi.ReplyKeyboardMarkup {
	var cmdKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(gameID),
		),
	)
	return cmdKeyboard
}

func ChooseIDKeyboard() tgbotapi.InlineKeyboardMarkup {
	var chooseIDKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ZONE ID из 5 цифр", "zone_id_5"),
			tgbotapi.NewInlineKeyboardButtonData("ZONE ID из 4 цифр", "zone_id_4"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "home"),
		),
	)
	return chooseIDKeyboard
}

func CatalogKeyboard(zone int) tgbotapi.InlineKeyboardMarkup {
	catalog, err := config.LoadCatalog(zone)
	if err != nil {
		log.Print(err)
		return tgbotapi.InlineKeyboardMarkup{}
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < len(catalog); i++ {
		buttonText := fmt.Sprintf("%s - %d", catalog[i].Name, catalog[i].Price)
		button := tgbotapi.NewInlineKeyboardButtonData(buttonText, catalog[i].CallbackData)

		if len(rows) == 0 || len(rows[len(rows)-1]) == 2 {
			rows = append(rows, []tgbotapi.InlineKeyboardButton{button})
		} else {
			rows[len(rows)-1] = append(rows[len(rows)-1], button)
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func ConfirmPurchaseKeyboard() tgbotapi.InlineKeyboardMarkup {
	var confirmPurchase = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Я оплатил йоу", "confirm_purchase"),
		),
	)
	return confirmPurchase
}

func ConfirmKeyboard(donationID string) tgbotapi.InlineKeyboardMarkup {
	var gameID = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Подтвердить", fmt.Sprintf("confirm_donate:%s", donationID)),
			tgbotapi.NewInlineKeyboardButtonData("Отменить", fmt.Sprintf("cancel_donate:%s", donationID)),
		),
	)
	return gameID
}

func ReviewsKeyboard() tgbotapi.InlineKeyboardMarkup {
	var reviews = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Отзывы", "t.me/donaterich"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "home"),
		),
	)
	return reviews
}

func GetCartKeyboard() tgbotapi.InlineKeyboardMarkup {
	var getCart = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Корзина", "view_cart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "home"),
		),
	)
	return getCart
}

func HomeKeyboard() tgbotapi.InlineKeyboardMarkup {
	var getCart = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Главное меню", "home"),
		),
	)
	return getCart
}
