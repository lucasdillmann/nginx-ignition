package main

import (
	"log"
)

func main() {
	log.Println("Welcome to nginx ignition")

	container, err := startContainer()
	if err != nil {
		log.Fatal("Application startup failed: ", err)
	}

	if err = container.Invoke(runApplication); err != nil {
		log.Fatal("Application startup failed: ", err)
	}
}
