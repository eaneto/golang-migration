package writer

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/eaneto/golang-migration/pkg/reader"
	"github.com/sirupsen/logrus"
)

// CreateMigrationTable Executes the SQL script that creates the migration table.
func CreateMigrationTable(tx *sql.Tx) error {
	script, err := readDefaultMigrationTableScript()
	if err != nil {
		return err
	}
	_, err = tx.Exec(string(script))
	if err != nil {
		logrus.Error("Error creating basic migration table.\n", err)
		return err
	}
	return nil
}

// readDefaultMigrationTableScript Reads the basic script for the migration table.
func readDefaultMigrationTableScript() ([]byte, error) {
	scriptContent, err := ioutil.ReadFile("create_migration_table.sql")
	if err != nil {
		logrus.Error("Error reading base migration file.\n", err)
		return nil, err
	}
	return scriptContent, nil
}

// processScripts Process all given scripts inside a single transaction.
func ProcessScripts(tx *sql.Tx, db *sql.DB, scripts []reader.SQLScript) error {
	for _, script := range scripts {
		err := processScript(tx, script)
		if err != nil {
			return err
		}
	}
	return nil
}

// processScript Process a given script inside the given transaction.
func processScript(tx *sql.Tx, script reader.SQLScript) error {
	isAlreadyProcessed, err := isScriptAlreadyExecuted(tx, script)
	if err != nil {
		return err
	}

	// If already processed ignore script and just log.
	if isAlreadyProcessed {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Info("Script already executed.")
	} else {
		err = executeScriptAndMarkAsExecuted(tx, script)
		if err != nil {
			return err
		}
	}
	return nil
}

// isScriptAlreadyExecuted Check if the script was alreayd executed by counting the rows
// in the migration table with the script name.
func isScriptAlreadyExecuted(tx *sql.Tx, script reader.SQLScript) (bool, error) {
	query := fmt.Sprintf("SELECT count(id) FROM golang_migration WHERE script_name = '%s'", script.Name)
	var count int
	err := tx.QueryRow(query).Scan(&count)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Error("Error checking if the script was alreayd executed.\n", err)
		return false, err
	}
	return count > 0, nil
}

// executeScriptAndMarkAsExecuted Executes the given script and mark it as executed.
func executeScriptAndMarkAsExecuted(tx *sql.Tx, script reader.SQLScript) error {
	err := executeScript(tx, script)
	if err != nil {
		return err
	}
	err = markScriptAsExecuted(tx, script)
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

// markScriptAsExecuted Insert the script name on the migration table.
func markScriptAsExecuted(tx *sql.Tx, script reader.SQLScript) error {
	query := fmt.Sprintf("INSERT INTO golang_migration (script_name) VALUES ('%s')", script.Name)
	_, err := tx.Exec(query)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Error("Error registering the executed script.\n", err)
		return err
	}
	return nil
}

// RollbackTransaction Rollback the given transaction.
func RollbackTransaction(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		logrus.Fatal("Error rollbacking transaction.\n", err)
	}
	logrus.Error("Migration executed unsuccessfully!")
}

// CommitTransaction Commit to the given transaction.
func CommitTransaction(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		logrus.Fatal("Error commiting transaction.\n", err)
	}
	logrus.Info("Migration executed successfully!")
}
