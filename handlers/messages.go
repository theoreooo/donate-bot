package handlers

import (
	"donate-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Messages(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	state := services.GetUserState(chatID)
	switch state {
	case "awaiting_game_id":
		services.SetGameId(bot, update)
	case "awaiting_donate_id":
		services.ConfirmDonate(bot, update)
	case "awaiting_bonus_donate_id":
		services.ConfirmBonusDonate(bot, update)
	default:
		switch text {
		case "Каталог":
			services.Catalog(bot, update)
		case "Корзина":
			services.SendCart(bot, update)
		case "Профиль":
			services.Profile(bot, update)
		case "Партнерская программа":
			services.Partnership(bot, update)
		case "Помощь":
			services.Help(bot, update)
		case "Отзывы":
			services.Reviews(bot, update)
		default:
			services.UnknownCommand(bot, update)
		}
	}
}
