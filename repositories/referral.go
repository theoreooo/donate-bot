package repositories

import (
	"donate-bot/models"
)

func GetReferralCount(referrerID int64) (int64, error) {
	var count int64
	err := DB.Model(&models.Referral{}).Where("referrer_id = ?", referrerID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetReferredDiamonds(tgID int64) (float64, error) {
	var user models.User
	err := DB.Where("telegram_id = ?", tgID).First(&user).Error
	if err != nil {
		return 0, err
	}

	return user.ReferredDiamonds, nil
}
