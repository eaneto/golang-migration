package processor

import (
	"database/sql"
	"fmt"

	"github.com/eaneto/grotto/internal/connection"
	"github.com/eaneto/grotto/pkg/executor"
	"github.com/eaneto/grotto/pkg/reader"
	"github.com/eaneto/grotto/pkg/registry"
	"github.com/sirupsen/logrus"
)

// MigrationProcessor Interface for the migration processor
type MigrationProcessor interface {
	ProcessMigration(migrationDirectory string)
}

// MigrationProcessorSQL Migration processor for SQL database.
type MigrationProcessorSQL struct {
	DatabaseInformation connection.DatabaseInformation
	DB                  *sql.DB
}

// DATABASE_URL Basic postgres connection string.  All options are
// replaced with command line arguments.
const DATABASE_URL = "postgres://%s:%s@%s:%s/%s"

// CreateProcessor Creates a migration processor with the given database information.
func CreateProcessor(databaseInformation connection.DatabaseInformation) MigrationProcessorSQL {
	return MigrationProcessorSQL{
		DatabaseInformation: databaseInformation,
		DB:                  stablishConnection(databaseInformation),
	}
}

// stablishConnection Stablished a connection with the database.
func stablishConnection(databaseInformation connection.DatabaseInformation) *sql.DB {
	db, err := sql.Open("pgx", fmt.Sprintf(DATABASE_URL, databaseInformation.User, databaseInformation.Password,
		databaseInformation.Address, databaseInformation.Port, databaseInformation.Database))
	if err != nil {
		logrus.Fatal("Failure stablishing database connection.\n", err)
	}
	return db
}

// ProcessMigration Process all migration located on the given directory.
func (m MigrationProcessorSQL) ProcessMigration(migrationDirectory string) {
	migrationReader := reader.MigrationReader{MigrationDirectory: migrationDirectory}
	scripts := migrationReader.ReadScriptFiles()

	executor := initializeExecutor(m.DB)
	createMigrationTable(executor)

	// Process all read scripts
	err := executor.ProcessScripts(scripts)

	// Only commits if all operations were succesful.
	if err != nil {
		executor.RollbackTransaction()
	} else {
		executor.CommitTransaction()
	}
}

// initializeExecutor Initialize the script executor with the database connection.
func initializeExecutor(db *sql.DB) executor.ScriptExecutor {
	tx, err := db.Begin()
	if err != nil {
		logrus.Fatal("Error starting transaction.\n", err)
	}

	return executor.ScriptExecutor{
		Tx: tx,
		MigrationRegister: registry.MigrationRegisterSQL{
			Tx: tx,
		},
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
