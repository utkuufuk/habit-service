package glados

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/utkuufuk/habit-service/internal/config"
	"github.com/utkuufuk/habit-service/internal/syslog"
)

type Listener struct {
	reader *kafka.Reader
	logger syslog.Logger
}

type CommandHandler func(args []string)

func NewListener(cfg config.Glados, logger syslog.Logger) Listener {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.Broker},
		GroupID:  cfg.GroupId,
		Topic:    cfg.Topic,
		MaxBytes: 1e6, // 1MB
	})

	return Listener{reader, logger}
}

func (l Listener) Close() {
	l.reader.Close()
}

func (l Listener) Listen(ctx context.Context, callback CommandHandler) {
	for {
		message, err := l.reader.ReadMessage(ctx)
		if err != nil {
			topic := l.reader.Config().Topic
			if err.Error() == "EOF" {
				l.logger.Info("Closing %s listener, EOF received.", topic)
				return
			}
			l.logger.Error("Could not read message from topic %s: %v", topic, err)
		}

		var cmd struct {
			Args []string `json:"args"`
		}

		json.Unmarshal(message.Value, &cmd)
		callback(cmd.Args)
	}
}
