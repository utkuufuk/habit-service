package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/utkuufuk/habit-service/internal/logger"
)

var (
	AppEnv             string
	HttpPort           int
	TimezoneLocation   *time.Location
	SpreadsheetId      string
	GoogleClientId     string
	GoogleClientSecret string
	GoogleAccessToken  string
	GoogleRefreshToken string
	TelegramChatId     int64
	TelegramToken      string
)

func init() {
	godotenv.Load()

	AppEnv = os.Getenv("APP_ENV")
	TelegramToken = os.Getenv("TELEGRAM_TOKEN")
	SpreadsheetId = os.Getenv("SPREADSHEET_ID")
	GoogleClientId = os.Getenv("GSHEETS_CLIENT_ID")
	GoogleClientSecret = os.Getenv("GSHEETS_CLIENT_SECRET")
	GoogleAccessToken = os.Getenv("GSHEETS_ACCESS_TOKEN")
	GoogleRefreshToken = os.Getenv("GSHEETS_REFRESH_TOKEN")

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = os.Getenv("HTTP_PORT")
	}

	port, err := strconv.Atoi(httpPort)
	if err != nil {
		logger.Error("PORT or HTTP_PORT not set")

		if AppEnv == "production" {
			os.Exit(1)
		}
	}

	location, err := time.LoadLocation(os.Getenv("TIMEZONE_LOCATION"))
	if err != nil {
		fmt.Printf(
			"Invalid timezone location: '%s', falling back to UTC: %v\n",
			location,
			err,
		)
		location, _ = time.LoadLocation("UTC")
	}

	chatId, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		logger.Error("Invalid Telegram Chat ID")

		if AppEnv == "production" {
			os.Exit(1)
		}
	}

	HttpPort = port
	TimezoneLocation = location
	TelegramChatId = chatId
}
