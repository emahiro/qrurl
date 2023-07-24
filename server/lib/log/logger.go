package log

import (
	"os"

	"golang.org/x/exp/slog"
)

var logger slog.Logger

func New() slog.Logger {
	return *slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
