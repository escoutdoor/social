package main

import (
	"log/slog"
	"os"

	"github.com/escoutdoor/social/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
