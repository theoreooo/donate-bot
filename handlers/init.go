package handlers

import (
	"donate-bot/config"

	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Init(bot *tgbotapi.BotAPI) {
	webhookURL := config.Config("WEBHOOK_URL") + "/webhook"

	webhook, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Print(err)
	}
	_, err = bot.Request(webhook)
	if err != nil {
		log.Print(err)
	}

	// u := tgbotapi.NewUpdate(0)
	// u.Timeout = 60
	// updates := bot.GetUpdatesChan(u)

	// for update := range updates {
	// 	if update.CallbackQuery != nil {
	// 		Callbacks(bot, update)
	// 	} else if update.Message.IsCommand() {
	// 		Commands(bot, update)
	// 	} else {
	// 		Messages(bot, update)
	// 	}
	// }
}
