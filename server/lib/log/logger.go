package log

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

var logger *slog.Logger

func New() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func Infof(ctx context.Context, title string, args ...any) {
	logger.InfoCtx(ctx, title, args...)
}
