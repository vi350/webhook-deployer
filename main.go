package main

import (
	"github.com/vi350/webhook-deployer/internal/app"
	"log"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Panic(err)
	}

	err = a.Run()
	if err != nil {
		log.Panic(err)
	}
}
