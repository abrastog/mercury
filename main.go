package main

import (
	"log"
	"os"
)

func main() {
	err := loadAppConfiguration()
	if err != nil {
		log.Fatalf("Failed to Application Configuration. %+v", err)
		os.Exit(1)
	}

	setupLogging()
	setupDatabase()
	setupTickHandler()

	setupKiteConnection()

	waitForCtrlC()
	log.Println("Exit")
}
