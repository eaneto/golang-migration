package writer

import (
	"database/sql"

	"github.com/eaneto/grotto/pkg/reader"
	"github.com/eaneto/grotto/pkg/registry"
	"github.com/sirupsen/logrus"
)

type Writer struct {
	Tx       *sql.Tx
	Registry registry.Registry
}

// processScripts Process all given scripts inside a single transaction.
func (w Writer) ProcessScripts(scripts []reader.SQLScript) error {
	for _, script := range scripts {
		err := w.processScript(script)
		if err != nil {
			return err
		}
	}
	return nil
}

// processScript Process a given script inside the given transaction.
func (w Writer) processScript(script reader.SQLScript) error {
	isAlreadyProcessed, err := w.Registry.IsScriptAlreadyExecuted(script)
	if err != nil {
		return err
	}

	// If already processed ignore script and just log.
	if isAlreadyProcessed {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Info("Script already executed.")
	} else {
		err = w.executeScriptAndMarkAsExecuted(script)
		if err != nil {
			return err
		}
	}
	return nil
}

// executeScriptAndMarkAsExecuted Executes the given script and mark it as executed.
func (w Writer) executeScriptAndMarkAsExecuted(script reader.SQLScript) error {
	err := executeScript(w.Tx, script)
	if err != nil {
		return err
	}
	err = w.Registry.MarkScriptAsExecuted(script)
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
func (w Writer) RollbackTransaction() {
	err := w.Tx.Rollback()
	if err != nil {
		logrus.Fatal("Error rollbacking transaction.\n", err)
	}
	logrus.Error("Migration executed unsuccessfully!")
}

// CommitTransaction Commit to the given transaction.
func (w Writer) CommitTransaction() {
	err := w.Tx.Commit()
	if err != nil {
		logrus.Fatal("Error commiting transaction.\n", err)
	}
	logrus.Info("Migration executed successfully!")
}
