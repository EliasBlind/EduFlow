package logger

import (
	"log/slog"
	"os"

	"github.com/EliasBlind/EduFlow/internal/journal_service/config"
)

func MustLoad(env config.Env) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)

	case config.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)

	case config.EnvProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)

	default:
		panic("unknown environment: " + string(env))
	}

	return log
}
