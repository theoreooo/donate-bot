package services

import (
	"donate-bot/config"
	"donate-bot/keyboards"

	"fmt"
	"log"
	"strconv"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CartItem struct {
	ItemID           string
	Name             string
	Price            int
	Quantity         int
	QuantityDiamonds int
}

type Cart struct {
	Items map[string]*CartItem
}

var userCarts = make(map[int64]*Cart)
var cartMu sync.Mutex

func AddToCart(bot *tgbotapi.BotAPI, update tgbotapi.Update, itemID string) {
	var err error
	cartMu.Lock()
	defer cartMu.Unlock()

	chatID := update.CallbackQuery.Message.Chat.ID

	itemIDint, err := strconv.Atoi(itemID)
	if err != nil {
		log.Print(err)
	}

	var catalog []config.Item

	catalog, err = config.LoadALlCatalog()
	if err != nil {
		log.Print(err)
	}

	if _, exists := userCarts[chatID]; !exists {
		userCarts[chatID] = &Cart{Items: make(map[string]*CartItem)}
	}

	cart := userCarts[chatID]

	if item, exists := cart.Items[itemID]; exists {
		item.Quantity++
	} else {
		cart.Items[itemID] = &CartItem{
			ItemID:           itemID,
			Name:             catalog[itemIDint].Name,
			Price:            catalog[itemIDint].Price,
			Quantity:         1,
			QuantityDiamonds: catalog[itemIDint].Quantity,
		}
	}

	text := fmt.Sprintf("Вы добавили в корзину: %s", catalog[itemIDint].Name)

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, text)
	msg.ReplyMarkup = keyboards.GetCartKeyboard()
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func SendCart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var chatID int64
	if update.CallbackQuery != nil {
		chatID = update.CallbackQuery.Message.Chat.ID
	} else if update.Message != nil {
		chatID = update.Message.Chat.ID
	}

	cart := GetCart(chatID)
	text, exists := FormatCart(cart)
	msg := tgbotapi.NewMessage(chatID, text)

	if exists {
		msg.ReplyMarkup = cartKeyboard(cart)
	}

	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func GetCart(chatID int64) *Cart {
	cartMu.Lock()
	defer cartMu.Unlock()

	if cart, exists := userCarts[chatID]; exists {
		return cart
	}
	return &Cart{Items: make(map[string]*CartItem)}
}

func FormatCart(cart *Cart) (string, bool) {
	exists := CartExists(cart)
	if !exists {
		return "Ваша корзина пуста", exists
	}
	text := "Ваша корзина:\n"
	total := 0
	totalDiamonds := 0

	for _, item := range cart.Items {
		text += fmt.Sprintf("%s (x%d) - %d сом\n", item.Name, item.Quantity, item.Price*item.Quantity)
		total += item.Price * item.Quantity
		totalDiamonds += item.QuantityDiamonds
	}

	text += fmt.Sprintf("\nИтого: %d Алмазов на сумму %d сом", totalDiamonds, total)
	return text, exists
}

func CartExists(cart *Cart) bool {
	return len(cart.Items) >= 1
}

func FormatCartAdmin(cart *Cart) string {
	text := "Корзина:\n"
	total := 0
	totalDiamonds := 0

	for _, item := range cart.Items {
		text += fmt.Sprintf("%s (x%d) - %d сом\n", item.Name, item.Quantity, item.Price*item.Quantity)
		total += item.Price * item.Quantity
		totalDiamonds += item.QuantityDiamonds
	}

	text += fmt.Sprintf("\nИтого: %d Алмазов на сумму %d сом", totalDiamonds, total)
	return text
}

func RemoveFromCart(chatID int64, itemID string) {
	cartMu.Lock()
	defer cartMu.Unlock()

	if cart, exists := userCarts[chatID]; exists {
		if cart.Items[itemID].Quantity > 1 {
			cart.Items[itemID].Quantity--
		} else {
			delete(cart.Items, itemID)
		}
	}
}

func ClearCart(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	cartMu.Lock()
	defer cartMu.Unlock()

	chatID := update.CallbackQuery.Message.Chat.ID
	delete(userCarts, chatID)

	text := "Вы очистили корзину."

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboards.MainKeyboard()
	if _, err := bot.Send(msg); err != nil {
		log.Print(err)
	}
}

func PurchaseClearCart(chatID int64) {
	cartMu.Lock()
	defer cartMu.Unlock()

	delete(userCarts, chatID)
}

func FinalPriceAndDiamonds(chatID int64) (int, int) {
	cart := GetCart(chatID)
	total := 0
	totalDiamonds := 0

	for _, item := range cart.Items {
		total += item.Price
		totalDiamonds += item.QuantityDiamonds
	}

	return total, totalDiamonds
}

func cartKeyboard(cart *Cart) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, item := range cart.Items {
		button := tgbotapi.NewInlineKeyboardButtonData(item.Name, "delete:"+item.ItemID)
		deleteButton := tgbotapi.NewInlineKeyboardButtonData("Удалить", "delete:"+item.ItemID)
		rows = append(rows, []tgbotapi.InlineKeyboardButton{button, deleteButton})
	}

	clearButton := tgbotapi.NewInlineKeyboardButtonData("Очистить корзину", "clear_cart")
	rows = append(rows, []tgbotapi.InlineKeyboardButton{clearButton})
	confirmButton := tgbotapi.NewInlineKeyboardButtonData("Подтвердить оплату", "confirm_cart")
	rows = append(rows, []tgbotapi.InlineKeyboardButton{confirmButton})
	homeButton := tgbotapi.NewInlineKeyboardButtonData("Главное меню", "home")
	rows = append(rows, []tgbotapi.InlineKeyboardButton{homeButton})

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
