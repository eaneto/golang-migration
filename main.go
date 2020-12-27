package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/eaneto/golang-migration/reader"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

// TODO get user, password and host from CLI
const DATABASE_URL = "postgres://%s:%s@localhost:5432/%s"

func main() {
	user := os.Args[1]
	password := os.Args[2]
	database := os.Args[3]

	db, err := sql.Open("pgx", fmt.Sprintf(DATABASE_URL, user, password, database))
	if err != nil {
		logrus.Fatal("Connection failed\n", err)
	}
	defer db.Close()
	scripts := reader.ReadScriptFiles()
	// TODO execute all migrations inside one transaction.
	tx, err := db.Begin()
	if err != nil {
		logrus.Fatal("Error starting transaction.\n", err)
	}

	// Creates the basic migration table.
	err = createMigrationTable(tx)
	if err != nil {
		logrus.Info("Rollbacking transacation.")
		tx.Rollback()
		panic(-1)
	}

	// Process all read scripts
	err = processScripts(tx, db, scripts)

	// Only commits if all operations were succesful.
	if err != nil {
		tx.Rollback()
		logrus.Error("Migration executed unsuccessfully!")
	} else {
		tx.Commit()
		logrus.Info("Migration executed successfully!")
	}
}

// processScripts Process all given scripts inside a single transaction.
func processScripts(tx *sql.Tx, db *sql.DB, scripts []reader.SQLScript) error {
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
	// TODO Check why this is happenning, ReadScriptFiles should never
	// return and empty script.
	if script.Name == "" {
		logrus.Warn("Empty file\n")
		return nil
	}
	isAlreadyProcessed, err := isScriptAlreadyExecuted(tx, script)
	if err != nil {
		logrus.Info("Rollbacking transacation.")
		return err
	}
	if !isAlreadyProcessed {
		err = executeScript(tx, script)
		if err != nil {
			logrus.Info("Rollbacking transacation.")
			return err
		}
		err = markScriptAsExecuted(tx, script)
		if err != nil {
			logrus.Info("Rollbacking transacation.")
			return err
		}
	} else {
		logrus.WithField("script_name", script.Name).Info("Script alreayd executed.\n")
	}
	return nil
}

// createMigrationTable Executes the SQL script that creates the migration table.
func createMigrationTable(tx *sql.Tx) error {
	scriptContent, err := ioutil.ReadFile("create_migration_table.sql")
	if err != nil {
		logrus.Error("Error reading base migration file.\n", err)
		return err
	}
	_, err = tx.Exec(string(scriptContent))
	if err != nil {
		logrus.Error("Error creating basic migration table.\n", err)
		return err
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

// executeScript Executes a given SQL script.
// Every script must be executed inside a transaction.
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
