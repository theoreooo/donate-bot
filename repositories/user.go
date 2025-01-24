package repositories

import (
	"log"

	"donate-bot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func IsUserExists(tgID int64) (bool, error) {
	var count int64

	err := DB.Model(&models.User{}).Where("telegram_id = ?", tgID).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func CreateUser(update tgbotapi.Update, referrerID *int64) error {
	user := models.User{
		TelegramID:       update.Message.Chat.ID,
		TelegramUsername: update.Message.Chat.UserName,
	}

	result := DB.Create(&user)

	if referrerID != nil {
		user.ReferredID = referrerID
		referral := models.Referral{
			ReferrerID:     *referrerID,
			ReferredUserID: update.Message.Chat.ID,
		}
		if result := DB.Create(&referral); result.Error != nil {
			return result.Error
		}
	}

	return result.Error
}

func SetGameId(chatID int64, gameID string) error {
	var user models.User
	if err := DB.Where("telegram_id", chatID).First(&user).Error; err != nil {
		log.Print(err)
	}

	user.GameId = gameID
	if err := DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func GetUserData(chatID int64) (models.User, error) {
	var user models.User
	if result := DB.Where("telegram_id = ?", chatID).Find(&user); result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func ResetReferredDiamonds(telegramID int64) error {
	var user models.User
	if err := DB.Where("telegram_id", telegramID).First(&user).Error; err != nil {
		log.Print(err)
	}

	user.ReferredDiamonds = 0
	if err := DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
