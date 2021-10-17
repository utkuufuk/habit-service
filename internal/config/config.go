package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Telegram struct {
	Token  string `yaml:"token"`
	ChatId int64  `yaml:"chat_id"`
}

type Config struct {
	HttpPort         int      `yaml:"http_port"`
	SpreadsheetId    string   `yaml:"spreadsheet_id"`
	TimezoneLocation string   `yaml:"timezone_location"`
	Telegram         Telegram `yaml:"telegram"`
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
