package repositories

import (
	"donate-bot/database"

	"gorm.io/gorm"
)

var DB *gorm.DB = database.Init()
