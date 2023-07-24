package log

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/exp/slog"
)

var logger *slog.Logger

func New() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func Infof(ctx context.Context, format string, args ...any) {
	now := time.Now()
	msg := fmt.Sprintf(format, args...)
	logger.InfoCtx(ctx, msg,
		"severity", slog.LevelInfo,
		"time", now.Format(time.RFC3339Nano),
		"message", msg,
		"insertId", "0",
		"spanId", "0",
		"trace", "0",
	)
}
