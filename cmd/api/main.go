package main

import (
	"github.com/escoutdoor/social/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}
