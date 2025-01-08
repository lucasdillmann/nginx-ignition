package main

import (
	"dillmann.com.br/nginx-ignition/application/boot"
	"log"
)

func main() {
	log.Println("Welcome to nginx ignition")

	if err := boot.StartApplication(); err != nil {
		log.Fatal("Application failed to start: ", err)
	}
}
