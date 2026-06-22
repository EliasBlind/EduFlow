package logger

import (
	"log/slog"
	"os"

	envutil "github.com/EliasBlind/EduFlow/pkg/env"
)

func MustLoad(env envutil.Env) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envutil.EnvLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)

	case envutil.EnvDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)

	case envutil.EnvProd:
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
