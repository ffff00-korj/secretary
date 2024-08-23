package main

import (
	"log"

	"github.com/ffff00-korj/secretary/internal/bot_app"
)

func main() {
	app := bot_app.NewApp()
	if err := app.Init(); err != nil {
		log.Fatal(err)
	}
	defer app.Close()
	upds, err := app.GetUpdates()
	if err != nil {
		log.Fatal(err)
	}
	for upd := range upds {
		if upd.Message == nil {
			log.Print("No messages are consumed!")
		}
		if err := app.ProcessAnUpdate(upd); err != nil {
			log.Print(err)
		}
	}
}
