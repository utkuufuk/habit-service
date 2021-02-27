package syslog

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Logger struct {
	api    *tgbotapi.BotAPI
	chatId int64
}

// NewLogger creates a new system logger instance
func NewLogger(chatId int64, botToken string) (l Logger) {
	api, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("could not create Telegram Logger: %v", err)
		return l
	}
	return Logger{api, chatId}
}

// Debug logs an debug message to stdout, but doesn't send a Telegram notification
func (l Logger) Debug(msg string, v ...interface{}) {
	msg = fmt.Sprintf(msg, v...)
	log.Println(msg)
}

// Info logs an informational message to stdout, and also sends a Telegram notification
func (l Logger) Info(msg string, v ...interface{}) {
	l.logf("Habit Service:", msg, v...)
}

// Warn logs a warning message to stdout, and also sends a Telegram notification
func (l Logger) Warn(msg string, v ...interface{}) {
	l.logf("Habit Service Warning:", msg, v...)
}

// Error logs an error message to stdout, and also sends a Telegram notification
func (l Logger) Error(msg string, v ...interface{}) {
	l.logf("Habit Service Error:", msg, v...)
}

// Fatal works like Error, but it returns with a non-zero exit code after logging
func (l Logger) Fatal(msg string, v ...interface{}) {
	l.Error(msg, v...)
	os.Exit(1)
}

// logf prints the message to stdout, and after prepending the given prefix to the message,
// also sends a Telegram notification
func (l Logger) logf(prefix, msg string, v ...interface{}) {
	msg = fmt.Sprintf(msg, v...)
	log.Println(msg)

	if l.api == nil || l.chatId == 0 {
		return
	}
	msg = fmt.Sprintf("%s %s", prefix, msg)
	m := tgbotapi.NewMessage(l.chatId, msg)
	l.api.Send(m)
}
