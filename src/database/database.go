package database

import (
	"commands/src/config"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Driver
)

func Connect() (*sql.DB, error) {
	db, error := sql.Open("mysql", config.StrDatabaseConnection)
	if error != nil {
		return nil, error
	}

	if error = db.Ping(); error != nil {
		db.Close()
		return nil, error
	}

	return db, nil
}
