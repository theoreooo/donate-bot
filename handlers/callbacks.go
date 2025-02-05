package handlers

import (
	"donate-bot/services"

	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CallbackHandler func(bot *tgbotapi.BotAPI, update tgbotapi.Update)

var callbackRoutes = map[string]CallbackHandler{
	"change_game_id":   services.ChangeGameIDCallback,
	"withdraw_bonus":   services.WithdrawBonusCallback,
	"view_cart":        services.SendCart,
	"clear_cart":       services.ClearCart,
	"confirm_cart":     services.ProcessPurchaseDiamonds,
	"confirm_purchase": services.ConfirmedPurchase,
	"home":             services.Home,
}

func Callbacks(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	data := update.CallbackQuery.Data

	switch {
	case strings.HasPrefix(data, "confirm_donate:"):
		donationID := strings.TrimPrefix(data, "confirm_donate:")
		services.ConfirmedDonate(bot, update, donationID)
		return

	case strings.HasPrefix(data, "cancel_donate:"):
		donationID := strings.TrimPrefix(data, "cancel_donate:")
		services.CanceledDonate(bot, update, donationID)
		return

	case strings.HasPrefix(data, "confirm_bonus:"):
		donationID := strings.TrimPrefix(data, "confirm_bonus:")
		services.ConfirmedDonateBonus(bot, update, donationID)
		return

	case strings.HasPrefix(data, "diamonds:"):
		itemID := strings.TrimPrefix(data, "diamonds:")
		services.AddToCart(bot, update, itemID)
		return

	case strings.HasPrefix(data, "delete:"):
		itemID := strings.TrimPrefix(data, "delete:")
		services.RemoveFromCart(update.CallbackQuery.Message.Chat.ID, itemID)
		services.SendCart(bot, update)
		return
	}

	if data == "zone_id_5" {
		services.Catalog(bot, update, 5)
		return
	}
	if data == "zone_id_4" {
		services.Catalog(bot, update, 4)
		return
	}

	if handler, exists := callbackRoutes[data]; exists {
		handler(bot, update)
	}
}
