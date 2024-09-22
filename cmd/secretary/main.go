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
	errorsChan := make(chan error)
	go func() {
		for upd := range upds {
			app.ProcessAnUpdate(upd, errorsChan)
		}
	}()
	for err := range errorsChan {
		log.Println(err)
	}
}
