package handlers

import (
	"donate-bot/services"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Commands(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Command() {
	case "start":
		services.Start(bot, update)
	case "help":
		services.Help(bot, update)
	case "home":
		services.Home(bot, update)
	}
}