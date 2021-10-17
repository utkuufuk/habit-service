package glados

import (
	"context"
	"fmt"
	"time"

	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/habit"
)

func RunCommand(
	ctx context.Context,
	client habit.Client,
	location *time.Location,
	args []string,
	cfg config.Config,
) (string, error) {
	if len(args) == 0 {
		return reportProgress(ctx, client, location, cfg.Telegram)
	}

	if args[0] == "mark" && len(args) == 3 {
		return "", markHabit(client, args[1], args[2])
	}

	return "", fmt.Errorf("could not parse glados command from args: '%v'", args)
}
