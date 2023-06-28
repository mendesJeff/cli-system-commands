package repositories

import (
	"commands/src/config"
	"commands/src/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type commands struct {
	db *sql.DB
}

func NewCommandsRepository(db *sql.DB) *commands {
	return &commands{db}
}

func (repository commands) SearchCommand() error {
	lines, error := repository.db.Query(
		"SELECT TC.id, t_clients.id AS client_id, TC.name, TC.parameters, COALESCE(TC.error, '') AS error, t_onus.onu_serial, olt_name, TC.last_update FROM t_commands TC JOIN t_olts OLT ON OLT.id = TC.olt_id JOIN t_clients ON t_clients.id = TC.client_id JOIN t_onus ON t_clients.onu_id = t_onus.id WHERE TC.error IS NULL OR TC.error = '' AND TC.last_update < NOW() ORDER BY TC.last_update;",
	)
	if error != nil {
		return error
	}
	defer lines.Close()

	var command models.Command
	var parameters string

	if lines.Next() {
		if error = lines.Scan(
			&command.ID,
			&command.Client_id,
			&command.Name,
			&parameters,
			&command.Error,
			&command.Onu_serial,
			&command.Olt_name,
			&command.Last_update,
		); error != nil {
			log.Println(error)
			return error
		}
	}

	// Define a slice to hold the parsed JSON data (array)
	var tl1_command []interface{}

	// Unmarshal the JSON array into the slice
	err := json.Unmarshal([]byte(parameters), &tl1_command)
	if err != nil {
		log.Println("tl1 has no command to apply")
		return nil
	}

	connection, erroConnTL1 := net.Dial("tcp", config.StrTl1Server)
	if erroConnTL1 != nil {
		log.Println(erroConnTL1)
		os.Exit(3)
		return fmt.Errorf("fail to connect to TL1 server")
	}

	login := fmt.Sprintf("LOGIN:::CTAG::UN=%s,PWD=%s;", os.Getenv("TL1_USER"), os.Getenv("TL1_PASS"))
	logout := "LOGOUT:::CTAG::;"

	// escrevendo a mensagem na conexão (socket)
	fmt.Fprintf(connection, login+"\n")
	log.Println(login)

	// ouvindo a resposta do servidor
	responseBytes := make([]byte, 1024) // Adjust the buffer size as per your needs
	n, readErr := connection.Read(responseBytes)
	if readErr != nil {
		log.Println(readErr)
		return fmt.Errorf("failed to read response from TL1 server")
	}

	response := string(responseBytes[:n])
	log.Println(response)

	// verificando retorno com erro
	index := strings.Index(response, "EADD=")
	if index != -1 {
		value := response[index+len("EADD="):] // Extrai a substring após "EADD="
		return fmt.Errorf(value)
	}

	for i := 0; i < len(tl1_command); i++ {

		// escrevendo a mensagem na conexão (socket)
		fmt.Fprintf(connection, "%v", tl1_command[i])
		log.Println(tl1_command[i])

		// ouvindo a resposta do servidor
		responseBytes := make([]byte, 1024) // Adjust the buffer size as per your needs
		n, readErr := connection.Read(responseBytes)
		if readErr != nil {
			log.Println(readErr)
			return fmt.Errorf("failed to read response from TL1 server")
		}

		response := string(responseBytes[:n])
		log.Println(response)

		// verificando retorno com erro
		index := strings.Index(response, "EADD=")
		if index != -1 {
			resultError := response[index+len("EADD="):] // Extrai a substring após "EADD="

			statement, err := repository.db.Prepare(
				"UPDATE t_commands SET error = ? WHERE id = ?")
			if err != nil {
				return err
			}
			defer statement.Close()

			// Execute the query with the JSON string
			_, err = statement.Exec(resultError, command.ID)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// fechando a sessão
			fmt.Fprintf(connection, logout+"\n")

			log.Println("logout successfully")

			return fmt.Errorf(resultError)
		}

	}

	statement, error := repository.db.Prepare(
		"DELETE FROM t_commands WHERE id = ?")

	if error != nil {
		return error
	}

	if _, error = statement.Exec(command.ID); error != nil {
		return error
	}

	return nil
}
