package main

import (
	"log"

	"myapp/bootstrap"
)

func main() {
	app, err := bootstrap.Boot()
	if err != nil {
		log.Fatalf("boot failed: %v", err)
	}

	if err := app.Run(); err != nil {
		app.Logger.Fatalf("application error: %v", err)
	}
}
