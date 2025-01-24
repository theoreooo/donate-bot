package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	TelegramID            int64     `gorm:"primaryKey;not null;unique"`
	TelegramUsername      string    `gorm:"type:varchar(255)"`
	RegistrationDate      time.Time `gorm:"type:timestamp;default:current_timestamp"`
	GameId                string    `gorm:"type:varchar(255)"`
	ReferredID            *int64    `gorm:"type:bigint"`
	ReferredDiamonds      float64   `gorm:"type:bigint;default:0"`
	TotalReferredDiamonds int64     `gorm:"type:bigint;default:0"`
	TotalDiamonds         int64     `gorm:"type:bigint;default:0"`
}

type Referral struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ReferrerID     int64     `gorm:"type:bigint;not null"`
	ReferredUserID int64     `gorm:"type:bigint;not null"`
	Date           time.Time `gorm:"type:timestamp;default:current_timestamp"`

	Referrer     User `gorm:"foreignKey:ReferrerID;references:TelegramID"`
	ReferredUser User `gorm:"foreignKey:ReferredUserID;references:TelegramID"`
}

type DonationHistory struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	TelegramID int64     `gorm:"type:bigint;not null"`
	GameID     string    `gorm:"type:varchar(255);not null"`
	Diamonds   int64     `gorm:"type:bigint;not null"`
	Date       time.Time `gorm:"type:timestamp;default:current_timestamp"`
	Status     string    `gorm:"type:varchar(50);not null"`
	Bonus      bool      `gorm:"type:boolean;default:false"`
}
