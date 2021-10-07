package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
	"github.com/utkuufuk/habit-service/internal/syslog"
)

var (
	cfg      config.Config
	client   habit.Client
	location *time.Location
	log      syslog.Logger
)

func init() {
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		fmt.Printf("Could not read config variables: %v", err)
		os.Exit(1)
	}

	log = syslog.NewLogger(cfg.Telegram.ChatId, cfg.Telegram.Token)

	location, err = time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		log.Warn(
			"Warning: invalid timezone location: '%s', falling back to UTC: %v",
			cfg.TimezoneLocation,
			err,
		)
		location, _ = time.LoadLocation("UTC")
	}

	client, err = habit.GetClient(context.Background(), cfg.Habit.SpreadsheetId, location)
	if err != nil {
		log.Fatal("Could not create gsheets client for Habit Service: %v", err)
	}
}

func main() {
	now := time.Now().In(location)
	currentHabits, err := client.FetchHabits(now)
	if err != nil {
		fmt.Printf("could not fetch this month's habits: %v\n", err)
	}

	fmt.Println("This Month:")
	for _, habit := range currentHabits {
		fmt.Printf("%s: %s%%\n", habit.Name, strconv.FormatFloat(habit.Score*100, 'f', 0, 32))
	}

	y, m, _ := now.Date()
	for _, i := range []int{1, 2, 3} {
		month := time.Month(int(m) + 1 - i)
		lastMonth := time.Date(y, month, 1, 0, 0, 0, 0, location).Add(-time.Nanosecond)
		habits, err := client.FetchHabits(lastMonth)
		if err != nil {
			fmt.Printf("could not fetch habits: %v\n", err)
		}

		fmt.Printf("\n%v:\n", month-1)
		for _, habit := range habits {
			fmt.Printf("%s: %s%%\n", habit.Name, strconv.FormatFloat(habit.Score*100, 'f', 0, 32))
		}
	}
}
