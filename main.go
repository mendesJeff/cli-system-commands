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

	fmt.Printf("\nRodando servidor de comandos\n")

	// Connect to the database
	db, err := database.Connect()
	if err != nil {
		log.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Start the loop
	for {

		controllers.SearchCommand()

		time.Sleep(5 * time.Second)
	}

}
