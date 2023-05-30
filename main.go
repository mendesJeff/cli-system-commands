package main

import (
	"commands/src/config"
	"commands/src/controllers"
	"commands/src/database"
	"fmt"
	"log"
	"time"
)

func main() {
	// Load environment variables
	config.EnveriormentsVariable()

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		log.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	fmt.Println("Login DB OK")

	// Start the loop
	for {

		controllers.SearchCommand()

		time.Sleep(5 * time.Second)
	}
}
