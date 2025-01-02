package log

import (
	"github.com/joakimen/clone"
	"io"
	"log/slog"
)

func ConfigureLogger(w io.Writer, cfg clone.Config) *slog.Logger {
	var handlerOpts slog.HandlerOptions
	if cfg.Verbose {
		handlerOpts.Level = slog.LevelInfo
	} else {
		handlerOpts.Level = slog.LevelWarn
	}
	logger := slog.New(slog.NewTextHandler(w, &handlerOpts))
	return logger
}
