package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Habit struct {
	SpreadsheetId   string `yaml:"spreadsheet_id"`
}

type Telegram struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	ChatId  int64  `yaml:"chat_id"`
}

type Config struct {
	Habit            Habit    `yaml:"habit"`
	HttpPort         int      `yaml:"http_port"`
	Telegram         Telegram `yaml:"telegram"`
	TimezoneLocation string   `yaml:"timezone_location"`
}

// ReadConfig reads the YAML config file & decodes all parameters
func ReadConfig(fileName string) (cfg Config, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return cfg, fmt.Errorf("could not open config file: %v", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return cfg, err
}
