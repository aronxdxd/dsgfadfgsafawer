package global

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Token  string `json:"token"`
	ChatID string `json:"chat_id"`
	EnergyRecoveryInterval int `json:"energyRecoveryInterval"`
	EnergyRecoveryAmount int `json:"energyRecoveryAmount"`
}

var AppConfig Config
func LoadConfig() (Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %w", err)
	}

	log.Println("Config loaded successfully")
	return config, nil
}
