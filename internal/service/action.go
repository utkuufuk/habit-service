package service

import (
	"context"

	"github.com/utkuufuk/habit-service/internal/habit"
)

type Action interface {
	Run(context.Context, habit.Client) (string, error)
}
