package handlers

import (
	"donate-bot/services"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Callbacks(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := update.CallbackQuery.Data

	switch {
	case strings.HasPrefix(data, "confirm_donate:"):
		donationID := strings.TrimPrefix(data, "confirm_donate:")
		services.ConfirmedDonate(bot, update, donationID)

	case strings.HasPrefix(data, "cancel_donate:"):
		donationID := strings.TrimPrefix(data, "cancel_donate:")
		services.CanceledDonate(bot, update, donationID)

	case strings.HasPrefix(data, "confirm_bonus:"):
		donationID := strings.TrimPrefix(data, "confirm_bonus:")
		services.ConfirmedDonateBonus(bot, update, donationID)

	case strings.HasPrefix(data, "diamonds:"):
		itemID := strings.TrimPrefix(data, "diamonds:")
		services.AddToCart(bot, update, itemID)

	case strings.HasPrefix(data, "delete:"):
		itemID := strings.TrimPrefix(data, "delete:")
		services.RemoveFromCart(update.CallbackQuery.Message.Chat.ID, itemID)
		services.SendCart(bot, update)

	default:
		switch data {
		case "change_game_id":
			services.ChangeGameIDCallback(bot, update)
		case "withdraw_bonus":
			services.WithdrawBonusCallback(bot, update)
		case "view_cart":
			services.SendCart(bot, update)
		case "clear_cart":
			services.ClearCart(bot, update)
		case "confirm_cart":
			services.ProcessPurchaseDiamonds(bot, update)
		case "confirm_purchase":
			services.ConfirmedPurchase(bot, update)
		case "home":
			services.Home(bot, update)
		}
	}
}
