package main

import (
	"github.com/VikaPaz/time_tracker/internal/app"
	"log"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
