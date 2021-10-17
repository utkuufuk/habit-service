package config

import (
	"log"
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

var Values Config

func init() {
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf("could not open config file: %v", err)
	}
	defer f.Close()

	var cfg config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("could not decode config: %v", err)
	}

	location, err := time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		log.Fatalf("Invalid timezone location: '%s': %v", cfg.TimezoneLocation, err)
	}

	Values = Config{
		cfg.HttpPort,
		cfg.SpreadsheetId,
		location,
		cfg.Telegram.Token,
		cfg.Telegram.ChatId,
	}
}
