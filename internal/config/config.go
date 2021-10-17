package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type telegram struct {
	Token  string `yaml:"token"`
	ChatId int64  `yaml:"chat_id"`
}

type config struct {
	HttpPort         int      `yaml:"http_port"`
	SpreadsheetId    string   `yaml:"spreadsheet_id"`
	TimezoneLocation string   `yaml:"timezone_location"`
	Telegram         telegram `yaml:"telegram"`
}

type Config struct {
	HttpPort         int
	SpreadsheetId    string
	TimezoneLocation *time.Location
	TelegramToken    string
	TelegramChatId   int64
}

// ReadConfig reads the YAML config file & decodes all parameters
func ReadConfig(fileName string) (Config, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return Config{}, fmt.Errorf("could not open config file: %w", err)
	}
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("could not decode config: %w", err)
	}

	location, err := time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		return Config{}, fmt.Errorf("Invalid timezone location: '%s': %v", cfg.TimezoneLocation, err)
	}

	return Config{
		cfg.HttpPort,
		cfg.SpreadsheetId,
		location,
		cfg.Telegram.Token,
		cfg.Telegram.ChatId,
	}, nil
}
