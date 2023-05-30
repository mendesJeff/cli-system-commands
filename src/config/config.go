package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	StrDatabaseConnection = ""
	StrTl1Server          = ""
	StrTl1Login           = ""
	Tl1Port               = 0
	DBPort                = 0
	SecretKey             []byte
)

func EnveriormentsVariable() {
	var error error

	if error = godotenv.Load(); error != nil {
		log.Fatal(error)
	}

	DBPort, error = strconv.Atoi(os.Getenv("API_PORT"))
	if error != nil {
		DBPort = 9000
	}

	Tl1Port, error = strconv.Atoi(os.Getenv("TL1_PORT"))
	if error != nil {
		Tl1Port = 5000
	}

	StrDatabaseConnection = fmt.Sprintf(
		"%s:%s@/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)

	StrTl1Server = fmt.Sprintf(
		"%s:%s",
		os.Getenv("TL1_SERVER"),
		os.Getenv("TL1_PORT"),
	)

	StrTl1Login = fmt.Sprintf(
		"LOGIN:::CTAG::UN=%s,PWD=%s;",
		os.Getenv("TL1_USER"),
		os.Getenv("TL1_PASS"),
	)

	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
