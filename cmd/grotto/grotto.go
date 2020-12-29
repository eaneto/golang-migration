package grotto

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/eaneto/grotto/pkg/executor"
	"github.com/eaneto/grotto/pkg/reader"
	"github.com/eaneto/grotto/pkg/registry"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

const DATABASE_URL = "postgres://%s:%s@%s:5432/%s"

func Run() {
	user := flag.String("user", "", "The user's name")
	password := flag.String("password", "", "The user's password")
	database := flag.String("database", "", "The database name")
	address := flag.String("addresss", "localhost", "Database server address")
	migrationDirectory := flag.String("dir", "", "The migration directory containing the scripts to be executed")

	flag.Parse()

	db, err := sql.Open("pgx", fmt.Sprintf(DATABASE_URL, *user, *password, *address, *database))
	if err != nil {
		logrus.Fatal("Failure stablishing database connection.\n", err)
	}
	defer db.Close()
	migrationReader := reader.MigrationReader{MigrationDirectory: *migrationDirectory}
	scripts := migrationReader.ReadScriptFiles()
	tx, err := db.Begin()
	if err != nil {
		logrus.Fatal("Error starting transaction.\n", err)
	}

	executor := executor.ScriptExecutor{
		Tx: tx,
		MigrationRegister: registry.MigrationRegisterSQL{
			Tx: tx,
		},
	}
	createMigrationTable(executor)

	// Process all read scripts
	err = executor.ProcessScripts(scripts)

	// Only commits if all operations were succesful.
	if err != nil {
		executor.RollbackTransaction()
	} else {
		executor.CommitTransaction()
	}
}

// createMigrationTable Creates the basic migration table.
func createMigrationTable(scriptExecutor executor.ScriptExecutor) {
	err := scriptExecutor.MigrationRegister.CreateMigrationTable()
	if err != nil {
		logrus.Error("Rollbacking transacation.")
		err = scriptExecutor.Tx.Rollback()
		if err != nil {
			logrus.Fatal("Error rollbacking transaction.\n", err)
		}
		panic(-1)
	}
}
