package config

import (
	"encoding/json"
	"os"
)

type Item struct {
	Name         string `json:"name"`
	Price        int    `json:"price"`
	CallbackData string `json:"callbackdata"`
	Quantity     int    `json:"quantity,omitempty"`
}

func LoadCatalog() ([]Item, error) {
	data, err := os.ReadFile("prices.json")
	if err != nil {
		return nil, err
	}

	var catalog []Item
	err = json.Unmarshal(data, &catalog)
	if err != nil {
		return nil, err
	}

	return catalog, nil
}
