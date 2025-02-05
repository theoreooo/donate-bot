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

func LoadCatalog(zone int) ([]Item, error) {
	var data []byte
	var err error
	if zone == 5 {
		data, err = os.ReadFile("prices5.json")
	} else {
		data, err = os.ReadFile("prices4.json")
	}

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

func LoadALlCatalog() ([]Item, error) {
	var data []byte
	var err error

	data, err = os.ReadFile("prices5.json")
	if err != nil {
		return nil, err
	}

	var catalog5 []Item
	err = json.Unmarshal(data, &catalog5)
	if err != nil {
		return nil, err
	}

	data, err = os.ReadFile("prices4.json")
	if err != nil {
		return nil, err
	}

	var catalog4 []Item
	err = json.Unmarshal(data, &catalog4)
	if err != nil {
		return nil, err
	}

	catalog := append(catalog5, catalog4...)
	return catalog, nil
}
