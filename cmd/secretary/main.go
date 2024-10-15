package main

import (
	"context"
	"log"

	"github.com/ffff00-korj/secretary/internal/bot_app"
)

func main() {
	ctx := context.Background()
	app := bot_app.NewApp()
	if err := app.Init(); err != nil {
		log.Fatal(err)
	}
	defer app.Close()
	upds, err := app.GetUpdates()
	if err != nil {
		log.Fatal(err)
	}
	for {
		app.ProcessAnUpdateWithContext(ctx, <-upds)
	}
}
