package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/eaneto/golang-migration/reader"
	"github.com/eaneto/golang-migration/writer"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

const DATABASE_URL = "postgres://%s:%s@localhost:5432/%s"

func main() {
	user := os.Args[1]
	password := os.Args[2]
	database := os.Args[3]

	db, err := sql.Open("pgx", fmt.Sprintf(DATABASE_URL, user, password, database))
	if err != nil {
		logrus.Fatal("Failure stablishing database connection.\n", err)
	}
	defer db.Close()
	scripts := reader.ReadScriptFiles()
	tx, err := db.Begin()
	if err != nil {
		logrus.Fatal("Error starting transaction.\n", err)
	}

	createMigrationTable(tx)

	// Process all read scripts
	err = writer.ProcessScripts(tx, db, scripts)

	// Only commits if all operations were succesful.
	if err != nil {
		writer.RollbackTransaction(tx)
	} else {
		writer.CommitTransaction(tx)
	}
}

// createMigrationTable Creates the basic migration table.
func createMigrationTable(tx *sql.Tx) {
	err := writer.CreateMigrationTable(tx)
	if err != nil {
		logrus.Error("Rollbacking transacation.")
		err = tx.Rollback()
		if err != nil {
			logrus.Fatal("Error rollbacking transaction.\n", err)
		}
		panic(-1)
	}
}
