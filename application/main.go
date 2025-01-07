package main

import (
	"dillmann.com.br/nginx-ignition/application/startup"
	"log"
)

func main() {
	log.Println("Welcome to nginx ignition")

	if err := startup.StartApplication(); err != nil {
		log.Fatal("Application failed to start: ", err)
	}
}
