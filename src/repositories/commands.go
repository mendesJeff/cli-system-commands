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

func (repository commands) SearchCommand() error {

	lines, error := repository.db.Query(
		"SELECT TC.id, t_clients.id AS client_id, TC.name AS command, TC.parameters, COALESCE(TC.error, '') AS error, t_onus.onu_serial, olt_name, TC.creation_date FROM t_commands TC 	JOIN t_olts OLT ON OLT.id = TC.olt_id JOIN t_clients ON t_clients.id = TC.client_id JOIN t_onus ON t_clients.onu_id = t_onus.id WHERE TC.error IS NULL ORDER BY TC.creation_date;",
	)
	if error != nil {
		return error
	}
	defer lines.Close()

	var command models.Command

	if lines.Next() {
		if error = lines.Scan(
			&command.ID,
			&command.Client_id,
			&command.Name,
			&command.Parameters,
			&command.Onu_serial,
			&command.Olt_name,
			&command.Error,
			&command.Creation_date,
		); error != nil {
			return error
		}
	}

	if command.Parameters == "" {
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

	// ouvindo a resposta do servidor (eco)
	scannerLogin := bufio.NewScanner(connection)
	scannerLogin.Split(bufio.ScanLines)

	for scannerLogin.Scan() {
		resultError := strings.Split(scannerLogin.Text(), "EADD=")
		if len(resultError) > 1 {
			return fmt.Errorf(resultError[1])
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

			statement, err := repository.db.Prepare(
				"UPDATE t_commands SET error = ? WHERE client_id = ?")
			if err != nil {
				return err
			}
			defer statement.Close()

			// Execute the query with the JSON string
			_, err = statement.Exec(resultError[1], command.Client_id)
			if err != nil {
				fmt.Println(err)
				return err
			}

			// fechando a sessão
			fmt.Fprintf(connection, logout+"\n")

			return fmt.Errorf(resultError[1])
		}

		if result := strings.Contains(scanner.Text(), "ENDESC=No error"); result {
			statement, error := repository.db.Prepare(
				"DELETE FROM t_commands WHERE id = ?")

			if error != nil {
				return error
			}

			if _, error = statement.Exec(command.ID); error != nil {
				return error
			}

			// fechando a sessão
			fmt.Fprintf(connection, logout+"\n")
			break
		}
	}

	return nil
}
