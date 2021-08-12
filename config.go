package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

// GetConfig - получение конфигурации
func GetConfig(path string) (*Config, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Ошибка загрузки конфигурации: " + err.Error())
	}

	config := new(Config)

	err = json.Unmarshal(raw, &config)
	if err != nil {
		return nil, errors.New("Ошибка парсинга конфигурации: " + err.Error())
	}

	log.Println()
	return config, nil
}
