package repositories

import (
	"donate-bot/models"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateDonate(update tgbotapi.Update, Diamonds int64) (string, error) {
	donation := models.DonationHistory{
		TelegramID: update.CallbackQuery.Message.Chat.ID,
		GameID:     update.CallbackQuery.Message.Text,
		Diamonds:   Diamonds,
		Status:     "pending",
	}

	if result := DB.Create(&donation); result.Error != nil {
		return "", result.Error
	}

	return donation.ID.String(), nil
}

func CreateDonateBonus(update tgbotapi.Update, referredDiamonds int64) (string, error) {
	donation := models.DonationHistory{
		TelegramID: update.Message.Chat.ID,
		GameID:     update.Message.Text,
		Diamonds:   referredDiamonds,
		Status:     "pending",
		Bonus:      true,
	}

	if result := DB.Create(&donation); result.Error != nil {
		return "", result.Error
	}

	return donation.ID.String(), nil
}

func SetDonateStatusCompleted(donationID string) error {
	var donation models.DonationHistory
	if err := DB.Where("id", donationID).First(&donation).Error; err != nil {
		log.Print(err)
	}
	var user models.User
	if err := DB.Where("telegram_id", donation.TelegramID).First(&user).Error; err != nil {
		log.Print(err)
	}
	user.TotalDiamonds += donation.Diamonds

	if err := DB.Save(&user).Error; err != nil {
		return err
	}
	var referrals []models.Referral

	if err := DB.Where("referred_user_id = ?", donation.TelegramID).Find(&referrals).Error; err != nil {
		return err
	}
	if len(referrals) > 0 {
		for _, referral := range referrals {
			var referrer models.User
			if err := DB.Where("telegram_id", referral.ReferrerID).First(&referrer).Error; err != nil {
				log.Print(err)
			}
			referrer.ReferredDiamonds += float64(donation.Diamonds) / 100
			referrer.TotalReferredDiamonds += donation.Diamonds / 100
			if err := DB.Save(&referrer).Error; err != nil {
				return err
			}
		}
	}

	donation.Status = "completed"
	if err := DB.Save(&donation).Error; err != nil {
		return err
	}

	return nil
}

func SetDonateStatusCanceled(donationID string) error {
	var donation models.DonationHistory
	if err := DB.Where("id", donationID).First(&donation).Error; err != nil {
		log.Print(err)
	}

	donation.Status = "canceled"
	if err := DB.Save(&donation).Error; err != nil {
		return err
	}

	return nil
}

func GetDataDonation(donationID string) (models.DonationHistory, error) {
	var donation models.DonationHistory

	err := DB.Where("id = ?", donationID).First(&donation).Error
	if err != nil {
		return donation, err
	}

	return donation, nil
}
