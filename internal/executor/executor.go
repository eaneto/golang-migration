package executor

import (
	"database/sql"

	"github.com/eaneto/grotto/internal/reader"
	"github.com/eaneto/grotto/internal/registry"
	"github.com/sirupsen/logrus"
)

// ScriptExecutor Basic structure to control script execution.
type ScriptExecutor struct {
	Tx                *sql.Tx
	MigrationRegister registry.MigrationRegister
}

// ProcessScripts Process all given scripts inside a single transaction.
func (executor ScriptExecutor) ProcessScripts(scripts []reader.SQLScript) error {
	for _, script := range scripts {
		err := executor.processScript(script)
		if err != nil {
			return err
		}
	}
	return nil
}

// processScript Process a given script inside the given transaction.
func (executor ScriptExecutor) processScript(script reader.SQLScript) error {
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
func (executor ScriptExecutor) executeScriptAndMarkAsExecuted(script reader.SQLScript) error {
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
func executeScript(tx *sql.Tx, script reader.SQLScript) error {
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
func (executor ScriptExecutor) RollbackTransaction() {
	err := executor.Tx.Rollback()
	if err != nil {
		logrus.Fatal("Error rollbacking transaction.\n", err)
	}
	logrus.Error("Migration executed unsuccessfully!")
}

// CommitTransaction Commit to the given transaction.
func (executor ScriptExecutor) CommitTransaction() {
	err := executor.Tx.Commit()
	if err != nil {
		logrus.Fatal("Error commiting transaction.\n", err)
	}
	logrus.Info("Migration executed successfully!")
}
