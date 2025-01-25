package handlers

import (
	"donate-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageHandler func(bot *tgbotapi.BotAPI, update tgbotapi.Update)

var messageRoutes = map[string]MessageHandler{
	"Каталог": services.Catalog,
	"Корзина": services.SendCart,
	"Профиль": services.Profile,
	"Партнерская программа": services.Partnership,
	"Помощь": services.Help,
	"Отзывы": services.Reviews,
}

func Messages(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	text := update.Message.Text

	state := services.GetUserState(chatID)
	switch state {
	case "awaiting_game_id":
		services.SetGameId(bot, update)
		return
	case "awaiting_donate_id":
		services.ConfirmDonate(bot, update)
		return
	case "awaiting_bonus_donate_id":
		services.ConfirmBonusDonate(bot, update)
		return
	}

	if handler, exists := messageRoutes[text]; exists {
		handler(bot, update)
	} else {
		services.UnknownCommand(bot, update)
	}
}
