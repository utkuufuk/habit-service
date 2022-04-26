package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/utkuufuk/habit-service/internal/logger"
)

type GoogleSheetsConfig struct {
	SpreadsheetId      string
	GoogleClientId     string
	GoogleClientSecret string
	GoogleAccessToken  string
	GoogleRefreshToken string
}

type ServerConfig struct {
	TimezoneLocation *time.Location
	GoogleSheets     GoogleSheetsConfig
	Port             int
	Secret           string
}

type ProgressReportConfig struct {
	TimezoneLocation *time.Location
	GoogleSheets     GoogleSheetsConfig
	TelegramChatId   int64
	TelegramToken    string
}

func ParseCommonConfig() (loc *time.Location, cfg GoogleSheetsConfig) {
	godotenv.Load()

	cfg.SpreadsheetId = os.Getenv("SPREADSHEET_ID")
	cfg.GoogleClientId = os.Getenv("GSHEETS_CLIENT_ID")
	cfg.GoogleClientSecret = os.Getenv("GSHEETS_CLIENT_SECRET")
	cfg.GoogleAccessToken = os.Getenv("GSHEETS_ACCESS_TOKEN")
	cfg.GoogleRefreshToken = os.Getenv("GSHEETS_REFRESH_TOKEN")

	loc, err := time.LoadLocation(os.Getenv("TIMEZONE_LOCATION"))
	if err != nil {
		logger.Warn(
			"Invalid timezone location: '%s', falling back to UTC: %v\n",
			os.Getenv("TIMEZONE_LOCATION"),
			err,
		)
		loc, _ = time.LoadLocation("UTC")
	}

	return loc, cfg
}

func ParseServerConfig() (cfg ServerConfig, err error) {
	loc, common := ParseCommonConfig()

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return cfg, fmt.Errorf("PORT not set")
	}

	return ServerConfig{
		loc,
		common,
		port,
		os.Getenv("SECRET"),
	}, nil
}

func ParseProgressReportConfig() (cfg ProgressReportConfig, err error) {
	loc, common := ParseCommonConfig()

	chatId, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		return cfg, fmt.Errorf("Invalid Telegram Chat ID")
	}

	return ProgressReportConfig{
		loc,
		common,
		chatId,
		os.Getenv("TELEGRAM_TOKEN"),
	}, nil
}
