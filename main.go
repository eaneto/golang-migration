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
const DATABASE_URL = "postgres://user:123@localhost:5432/"

// TODO Create a table to registered processed scripts.
func main() {
	databaseName := os.Args[1]
	db, err := sql.Open("pgx", DATABASE_URL+databaseName)
	if err != nil {
		logrus.Fatal("Connection failed\n", err)
	}
	defer db.Close()
	scripts := reader.ReadScriptFiles()
	// TODO execute all migrations inside one transaction.
	createMigrationTable(db)
	for _, script := range scripts {
		// TODO Check why this is happenning, ReadScriptFiles should never
		// return and empty script.
		if script.Name == "" {
			logrus.Warn("Empty file\n")
			continue
		}
		if !isScriptAlreadyExecuted(db, script) {
			executeScript(db, script)
			markScriptAsExecuted(db, script)
		} else {
			logrus.WithField("script_name", script.Name).Info("Script alreayd executed.\n")
		}
	}
	logrus.Info("Migration executed successfully!")
}

// createMigrationTable Executes the SQL script that creates the migration table.
func createMigrationTable(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		logrus.Fatal("Error starting transacation.\n", err)
	}
	scriptContent, err := ioutil.ReadFile("create_migration_table.sql")
	if err != nil {
		logrus.Fatal("Error reading base migration file.\n", err)
	}
	_, err = db.Exec(string(scriptContent))
	if err != nil {
		tx.Rollback()
		logrus.Fatal("Error creating basic migration table.\n", err)
	} else {
		tx.Commit()
	}
}

// isScriptAlreadyExecuted Check if the script was alreayd executed by counting the rows
// in the migration table with the script name.
func isScriptAlreadyExecuted(db *sql.DB, script reader.SQLScript) bool {
	query := fmt.Sprintf("SELECT count(id) FROM golang_migration WHERE script_name = '%s'", script.Name)
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Fatal("Error checking if the script was alreayd executed.\n", err)
	}
	return count > 0
}

// executeScript Executes a given SQL script.
// Every script must be executed inside a transaction.
func executeScript(db *sql.DB, script reader.SQLScript) {
	tx, err := db.Begin()
	if err != nil {
		logrus.Fatal("Error starting transaction\n", err)
	}
	logrus.Info("Executing script: ", script.Name)
	_, err = db.Exec(script.Content)
	if err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Fatal("Error executing script.\n", err)
	} else {
		tx.Commit()
	}
}

// markScriptAsExecuted Insert the script name on the migration table.
func markScriptAsExecuted(db *sql.DB, script reader.SQLScript) {
	tx, err := db.Begin()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Fatal("Error starting transaction to register script\n", err)
	}
	query := fmt.Sprintf("INSERT INTO golang_migration (script_name) VALUES ('%s')", script.Name)
	_, err = db.Exec(query)
	if err != nil {
		tx.Rollback()
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Fatal("Error registering the executed script.\n", err)
	} else {
		tx.Commit()
	}
}
