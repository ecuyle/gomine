package db

import (
	"database/sql"

	"github.com/ecuyle/gomine/internal/servermanager"
	_ "github.com/mattn/go-sqlite3"
)

type ServerOptions struct {
	Name           string                 `json:"name"`
	UserID         string                 `json:"userId"`
	Runtime        string                 `json:"runtime"`
	IsEulaAccepted bool                   `json:"isEulaAccepted"`
	Config         map[string]interface{} `json:"config"`
}

func InsertServer(server *servermanager.MCServer, options *ServerOptions) error {
	db, err := sql.Open("sqlite3", "./gomine.db")

	if err != nil {
		return err
	}

	defer db.Close()
	transaction, err := db.Begin()

	if err != nil {
		return err
	}

	statement, err := transaction.Prepare("insert into servers(id, name, runtime, path, pid, status, user_id) values(?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(server.ID, server.Name, options.Runtime, server.Path, nil, false, options.UserID)

	if err != nil {
		return err
	}

	err = transaction.Commit()

	if err != nil {
		return err
	}

	return nil
}
