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
}
