package controllers

import (
	"commands/src/database"
	"commands/src/repositories"
)

func SearchCommand() {

	db, error := database.Connect()
	if error != nil {

		return
	}
	defer db.Close()

	// Acessa o repositorio para consultar os comandos
	repository := repositories.NewCommandsRepository(db)
	error = repository.SearchCommand()
	if error != nil {
		return
	}

}
