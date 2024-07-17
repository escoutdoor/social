package logger

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

func SetupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag := strings.ToLower(os.Getenv("LEVEL"))
	switch flag {
	case "debug", "d":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "error", "e":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "info", "i":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}
}
