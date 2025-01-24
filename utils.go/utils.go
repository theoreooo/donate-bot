package utils

import (
	"donate-bot/config"

	"log"
	"strconv"
)

func GetAdminChatID() int64 {
	adminChatID, err := strconv.Atoi(config.Config("adminChatID"))
	if err != nil {
		log.Print(err)
	}
	return int64(adminChatID)
}
