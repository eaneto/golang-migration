package registry

import (
	"database/sql"
	"fmt"

	"github.com/eaneto/grotto/pkg/reader"
	"github.com/sirupsen/logrus"
)

// MIGRATION_TABLE_NAME The name of the table that stores all the executed migrations.
const MIGRATION_TABLE_NAME = "grotto_migration"

// DEFAULT_MIGRATION_SCRIPT the basic script for the migration table.
// This table is responsible to store the scripts executed so they won't be
// executed multiple times.
// The table stores a sequencial id, the name of the script which was executed
// and the date it was created. The script name is a unique field.
const DEFAULT_MIGRATION_SCRIPT = `create table if not exists ` +
	MIGRATION_TABLE_NAME +
	`(id bigint generated always as identity primary key,
script_name varchar constraint uk_script_name unique not null,
created_at timestamp not null default now());`

type MigrationRegister struct {
	Tx *sql.Tx
}

// CreateMigrationTable Executes the SQL script that creates the migration table.
func (r MigrationRegister) CreateMigrationTable() error {
	_, err := r.Tx.Exec(DEFAULT_MIGRATION_SCRIPT)
	if err != nil {
		logrus.Error("Error creating basic migration table.\n", err)
		return err
	}
	return nil
}

// IsScriptAlreadyExecuted Check if the script was alreayd executed by counting the rows
// in the migration table with the script name.
func (r MigrationRegister) IsScriptAlreadyExecuted(script reader.SQLScript) (bool, error) {
	query := fmt.Sprintf("SELECT count(id) FROM %s WHERE script_name = '%s'",
		MIGRATION_TABLE_NAME, script.Name)
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
	query := fmt.Sprintf("INSERT INTO %s (script_name) VALUES ('%s')",
		MIGRATION_TABLE_NAME, script.Name)
	_, err := r.Tx.Exec(query)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"script_name": script.Name,
		}).Error("Error registering the executed script.\n", err)
		return err
	}
	return nil
}
