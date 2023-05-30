package controllers

import (
	"commands/src/database"
	"commands/src/repositories"
	"fmt"
)

func SearchCommand() {

	db, error := database.Connect()
	if error != nil {

		return
	}
	defer db.Close()

	repository := repositories.NewCommandsRepository(db)
	command, error := repository.SearchCommand()
	if error != nil {

		return
	}

	fmt.Println(command)

}
