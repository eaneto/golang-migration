package registry

import (
	"database/sql"
	"fmt"

	"github.com/eaneto/grotto/pkg/reader"
	"github.com/sirupsen/logrus"
)

type MigrationRegister struct {
	Tx *sql.Tx
}

// CreateMigrationTable Executes the SQL script that creates the migration table.
func (r MigrationRegister) CreateMigrationTable() error {
	script := getDefaultMigrationTableScript()
	_, err := r.Tx.Exec(string(script))
	if err != nil {
		logrus.Error("Error creating basic migration table.\n", err)
		return err
	}
	return nil
}

// getDefaultMigrationTableScript Gets the basic script for the migration table.
func getDefaultMigrationTableScript() string {
	// This table is responsible to store the scripts executed so they won't be
	// executed multiple times.
	// The table stores a sequencial id, the name of the script which was executed
	// and the date it was created. The script name is a unique field.
	return `create table if not exists grotto_migration(
id bigint generated always as identity primary key,
script_name varchar constraint uk_script_name unique not null,
created_at timestamp not null default now()
);`
}

// IsScriptAlreadyExecuted Check if the script was alreayd executed by counting the rows
// in the migration table with the script name.
func (r MigrationRegister) IsScriptAlreadyExecuted(script reader.SQLScript) (bool, error) {
	query := fmt.Sprintf("SELECT count(id) FROM grotto_migration WHERE script_name = '%s'", script.Name)
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
	query := fmt.Sprintf("INSERT INTO grotto_migration (script_name) VALUES ('%s')", script.Name)
	_, err := r.Tx.Exec(query)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Error("Error registering the executed script.\n", err)
		return err
	}
	return nil
}
