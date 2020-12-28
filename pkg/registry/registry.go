package registry

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/eaneto/grotto/pkg/reader"
	"github.com/sirupsen/logrus"
)

type MigrationRegister struct {
	Tx *sql.Tx
}

// CreateMigrationTable Executes the SQL script that creates the migration table.
func (r MigrationRegister) CreateMigrationTable() error {
	script, err := readDefaultMigrationTableScript()
	if err != nil {
		return err
	}
	_, err = r.Tx.Exec(string(script))
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

// IsScriptAlreadyExecuted Check if the script was alreayd executed by counting the rows
// in the migration table with the script name.
func (r MigrationRegister) IsScriptAlreadyExecuted(script reader.SQLScript) (bool, error) {
	query := fmt.Sprintf("SELECT count(id) FROM golang_migration WHERE script_name = '%s'", script.Name)
	var count int
	err := r.Tx.QueryRow(query).Scan(&count)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Error("Error checking if the script was alreayd executed.\n", err)
		return false, err
	}
	return count > 0, nil
}

// MarkScriptAsExecuted Insert the script name on the migration table.
func (r MigrationRegister) MarkScriptAsExecuted(script reader.SQLScript) error {
	query := fmt.Sprintf("INSERT INTO golang_migration (script_name) VALUES ('%s')", script.Name)
	_, err := r.Tx.Exec(query)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Error("Error registering the executed script.\n", err)
		return err
	}
	return nil
}
