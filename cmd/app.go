package main

import (
	"log"
	"prices/internal/app"
)

func main() {
	a, cleanup, err := app.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer cleanup()

	a.Run()
}
