package executor

import (
	"database/sql"

	"github.com/eaneto/grotto/internal/registry"
	"github.com/eaneto/grotto/pkg/database"
	"github.com/sirupsen/logrus"
)

// ScriptExecutor Basic interface for the script executor.
type ScriptExecutor interface {
	CreateMigrationTable() error
	ProcessScripts(scripts []database.SQLScript) error
	RollbackTransaction()
	CommitTransaction()
}

// ScriptExecutorSQL Basic structure to control script execution.
type ScriptExecutorSQL struct {
	Tx                *sql.Tx
	MigrationRegister registry.MigrationRegister
}

// CreateMigrationTable Creates the migration table with the migration register.
func (executor ScriptExecutorSQL) CreateMigrationTable() error {
	return executor.MigrationRegister.CreateMigrationTable()
}

// ProcessScripts Process all given scripts inside a single transaction.
func (executor ScriptExecutorSQL) ProcessScripts(scripts []database.SQLScript) error {
	for _, script := range scripts {
		err := executor.processScript(script)
		if err != nil {
			return err
		}
	}
	return nil
}

// processScript Process a given script inside the given transaction.
func (executor ScriptExecutorSQL) processScript(script database.SQLScript) error {
	isAlreadyProcessed, err := executor.MigrationRegister.IsScriptAlreadyExecuted(script)
	if err != nil {
		return err
	}

	// If already processed ignore script and just log.
	if isAlreadyProcessed {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Info("Script already executed.")
	} else {
		err = executor.executeScriptAndMarkAsExecuted(script)
		if err != nil {
			return err
		}
	}
	return nil
}

// executeScriptAndMarkAsExecuted Executes the given script and mark it as executed.
func (executor ScriptExecutorSQL) executeScriptAndMarkAsExecuted(script database.SQLScript) error {
	err := executeScript(executor.Tx, script)
	if err != nil {
		return err
	}
	err = executor.MigrationRegister.MarkScriptAsExecuted(script)
	if err != nil {
		return err
	}
	return nil
}

// executeScript Executes a given SQL script.
func executeScript(tx *sql.Tx, script database.SQLScript) error {
	logrus.Info("Executing script: ", script.Name)
	_, err := tx.Exec(script.Content)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Error("Error executing script.\n", err)
		return err
	}
	return nil
}

// RollbackTransaction Rollback the given transaction.
func (executor ScriptExecutorSQL) RollbackTransaction() {
	err := executor.Tx.Rollback()
	if err != nil {
		logrus.Fatal("Error rollbacking transaction.\n", err)
	}
	logrus.Error("Migration executed unsuccessfully!")
}

// CommitTransaction Commit to the given transaction.
func (executor ScriptExecutorSQL) CommitTransaction() {
	err := executor.Tx.Commit()
	if err != nil {
		logrus.Fatal("Error commiting transaction.\n", err)
	}
	logrus.Info("Migration executed successfully!")
}
