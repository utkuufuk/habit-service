package glados

import (
	"fmt"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/service"
)

func ParseCommand(args []string, cfg config.Config) (service.Action, error) {
	if len(args) == 0 {
		return service.ReportProgressAction{
			TelegramChatId:   cfg.TelegramChatId,
			TelegramToken:    cfg.TelegramToken,
			TimezoneLocation: cfg.TimezoneLocation,
		}, nil
	}

	if args[0] == "mark" && len(args) == 3 {
		return service.MarkHabitAction{
			Cell:   args[1],
			Symbol: args[2],
		}, nil
	}

	return nil, fmt.Errorf("could not parse glados command from args: '%v'", args)
}
