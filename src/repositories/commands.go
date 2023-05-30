package repositories

import (
	"bufio"
	"commands/src/config"
	"commands/src/models"
	"database/sql"
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

func (repository commands) SearchCommand() (models.Command, error) {

	lines, error := repository.db.Query(
		"SELECT * FROM t_commands LIMIT 1",
	)
	if error != nil {
		return models.Command{}, error
	}
	defer lines.Close()

	var command models.Command

	if lines.Next() {
		if error = lines.Scan(
			&command.ID,
			&command.Name,
			&command.Parameters,
			&command.Olt_id,
			&command.Last_update,
			&command.Creation_date,
		); error != nil {
			return models.Command{}, error
		}
	}

	if command.Parameters == "" {
		fmt.Println("Nao tem comando na fila")
		return command, nil
	}

	fmt.Println(command.Parameters)

	connection, erroConnTL1 := net.Dial("tcp", config.StrTl1Server)
	if erroConnTL1 != nil {
		log.Println(erroConnTL1)
		os.Exit(3)
	}

	fmt.Println("conexao tl1 ok")

	login := fmt.Sprintf("LOGIN:::CTAG::UN=%s,PWD=%s;", os.Getenv("TL1_USER"), os.Getenv("TL1_PASS"))
	logout := "LOGOUT:::CTAG::;"

	// escrevendo a mensagem na conexão (socket)
	fmt.Fprintf(connection, login+"\n")

	// ouvindo a resposta do servidor (eco)
	scannerLogin := bufio.NewScanner(connection)
	scannerLogin.Split(bufio.ScanLines)

	for scannerLogin.Scan() {
		resultError := strings.Split(scannerLogin.Text(), "EADD=")
		if len(resultError) > 1 {
			return command, fmt.Errorf(resultError[1])
		}

		if result := strings.Contains(scannerLogin.Text(), "ENDESC=No error"); result {
			break
		}
	}

	// escrevendo a mensagem na conexão (socket)
	fmt.Fprintf(connection, command.Parameters+"\n")

	// ouvindo a resposta do servidor (eco)
	scanner := bufio.NewScanner(connection)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		resultError := strings.Split(scanner.Text(), "EADD=")

		if len(resultError) > 1 {

			statement, err := repository.db.Prepare("INSERT INTO t_commands_errors (error, command_id) VALUES (?, ?)")
			if err != nil {
				return models.Command{}, err
			}
			defer statement.Close()

			// Execute the query with the JSON string
			_, err = statement.Exec(resultError[1], command.ID)
			if err != nil {
				return models.Command{}, err
			}

			// fechando a sessão
			fmt.Fprintf(connection, logout+"\n")

			fmt.Println(resultError[1])
			return command, fmt.Errorf(resultError[1])
		}

		if result := strings.Contains(scanner.Text(), "ENDESC=No error"); result {
			statementError, error := repository.db.Prepare(
				"DELETE FROM t_commands_errors WHERE command_id = ?")

			if error != nil {
				return command, error
			}

			if _, error = statementError.Exec(command.ID); error != nil {
				return command, error
			}

			statement, error := repository.db.Prepare(
				"DELETE FROM t_commands WHERE id = ?")

			if error != nil {
				return command, error
			}

			if _, error = statement.Exec(command.ID); error != nil {
				return command, error
			}

			// fechando a sessão
			fmt.Fprintf(connection, logout+"\n")
			break
		}
	}

	return command, nil
}
