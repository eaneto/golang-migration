package processor

import (
	"database/sql"
	"fmt"

	"github.com/eaneto/grotto/internal/executor"
	"github.com/eaneto/grotto/internal/reader"
	"github.com/eaneto/grotto/internal/registry"
	"github.com/eaneto/grotto/pkg/connection"
	"github.com/sirupsen/logrus"
)

// MigrationProcessor Interface for the migration processor
type MigrationProcessor interface {
	ProcessMigration()
}

// MigrationProcessorSQL Migration processor for SQL database.
type MigrationProcessorSQL struct {
	Executor executor.ScriptExecutor
	Reader   reader.MigrationReader
}

// DATABASE_URL Basic postgres connection string.  All options are
// replaced with command line arguments.
const DATABASE_URL = "postgres://%s:%s@%s:%s/%s"

// New Creates a migration processor with the given database information.
func New(databaseInformation connection.DatabaseInformation, migrationDirecetory string) MigrationProcessorSQL {
	return MigrationProcessorSQL{
		Executor: initializeExecutor(stablishConnection(databaseInformation)),
		Reader: reader.MigrationReaderFS{
			MigrationDirectory: migrationDirecetory,
		},
	}
}

// ProcessMigration Process all migration located on the given directory.
func (m MigrationProcessorSQL) ProcessMigration() {
	// Creates migration table
	createMigrationTable(m.Executor)

	// Read all scripts on the migration directory
	scripts := m.Reader.ReadScriptFiles()

	// Process all read scripts
	err := m.Executor.ProcessScripts(scripts)

	// Only commits if all operations were succesful.
	if err != nil {
		m.Executor.RollbackTransaction()
	} else {
		m.Executor.CommitTransaction()
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

// initializeExecutor Initialize the script executor with the database connection.
func initializeExecutor(db *sql.DB) executor.ScriptExecutorSQL {
	tx, err := db.Begin()
	if err != nil {
		logrus.Fatal("Error starting transaction.\n", err)
	}

	return executor.ScriptExecutorSQL{
		Tx: tx,
		MigrationRegister: registry.MigrationRegisterSQL{
			Tx: tx,
		},
	}
}

// createMigrationTable Creates the basic migration table.
func createMigrationTable(scriptExecutor executor.ScriptExecutor) {
	err := scriptExecutor.CreateMigrationTable()
	if err != nil {
		logrus.Error("Rollbacking transacation.")
		scriptExecutor.RollbackTransaction()
	}
}
