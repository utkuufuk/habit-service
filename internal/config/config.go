package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Glados struct {
	Broker  string `yaml:"broker"`
	GroupId string `yaml:"group_id"`
	Topic   string `yaml:"topic"`
}

type Habit struct {
	SpreadsheetId string `yaml:"spreadsheet_id"`
}

type Telegram struct {
	Token  string `yaml:"token"`
	ChatId int64  `yaml:"chat_id"`
}

type Config struct {
	Glados           Glados   `yaml:"glados"`
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
