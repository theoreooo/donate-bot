package services

import (
	"sync"
)

type UserState struct {
	State  string
	ItemID int
	GameID string
}

var userStates = make(map[int64]*UserState)
var mu sync.Mutex

func SetUserGameID(chatID int64, state, gameID string) {
	mu.Lock()
	defer mu.Unlock()
	if userState, exists := userStates[chatID]; exists {
		userState.State = state
		userState.GameID = gameID
	} else {
		userStates[chatID] = &UserState{State: state, ItemID: 0, GameID: gameID}
	}
}

func GetUserGameID(chatID int64) string {
	mu.Lock()
	defer mu.Unlock()
	if user, exists := userStates[chatID]; exists {
		return user.GameID
	}
	return ""
}

func SetUserState(chatID int64, state string) {
	mu.Lock()
	defer mu.Unlock()
	if userState, exists := userStates[chatID]; exists {
		userState.State = state
	} else {
		userStates[chatID] = &UserState{State: state, ItemID: 0}
	}
}

func GetUserState(chatID int64) string {
	mu.Lock()
	defer mu.Unlock()

	if userState, exists := userStates[chatID]; exists && userState != nil {
		return userState.State
	}

	return ""
}
