package handlers

import (
	"donate-bot/clients"

	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2"
)

func WebhookHandler(c *fiber.Ctx) error {
	bot := clients.Init()
	var update tgbotapi.Update

	if err := c.BodyParser(&update); err != nil {
		log.Println("Error parsing update: ", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if update.CallbackQuery != nil {
		Callbacks(bot, update)
	} else if update.Message.IsCommand() {
		Commands(bot, update)
	} else {
		Messages(bot, update)
	}

	return c.SendStatus(fiber.StatusOK)
}
