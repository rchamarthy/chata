package chata

import (
	"context"
	"io"

	"log/slog"
)

type LOG_KEY_TYPE string

const LOG_KEY = LOG_KEY_TYPE("slog")

func Log(ctx context.Context) *slog.Logger {
	log := ctx.Value(LOG_KEY)
	if log != nil {
		return log.(*slog.Logger)
	}

	return NilLogger()
}

func NilLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{}))
}
