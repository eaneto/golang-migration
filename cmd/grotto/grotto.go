package grotto

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/eaneto/grotto/pkg/reader"
	"github.com/eaneto/grotto/pkg/registry"
	"github.com/eaneto/grotto/pkg/writer"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

const DATABASE_URL = "postgres://%s:%s@localhost:5432/%s"

func Run() {
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

	writer := writer.Writer{
		Tx: tx,
		Registry: registry.Registry{
			Tx: tx,
		},
	}
	createMigrationTable(writer)

	// Process all read scripts
	err = writer.ProcessScripts(scripts)

	// Only commits if all operations were succesful.
	if err != nil {
		writer.RollbackTransaction()
	} else {
		writer.CommitTransaction()
	}
}

// createMigrationTable Creates the basic migration table.
func createMigrationTable(writer writer.Writer) {
	err := writer.Registry.CreateMigrationTable()
	if err != nil {
		logrus.Error("Rollbacking transacation.")
		err = writer.Tx.Rollback()
		if err != nil {
			logrus.Fatal("Error rollbacking transaction.\n", err)
		}
		panic(-1)
	}
}
