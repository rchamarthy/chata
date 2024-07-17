package chata

import (
	"context"
	"io"
	"log/slog"
)

type LogKeyType string

const LogKey = LogKeyType("slog")

func Log(ctx context.Context) *slog.Logger {
	log := ctx.Value(LogKey)
	if log != nil {
		logger, ok := log.(*slog.Logger)

		if !ok {
			return NilLogger()
		}

		return logger
	}

	return NilLogger()
}

func NilLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}
