package registry

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eaneto/grotto/pkg/reader"
	"github.com/stretchr/testify/assert"
)

func TestCreatingDefaultMigrationTableWithNilTransactionShouldPanic(t *testing.T) {
	// If the transaction is not initialized the register should not start a new transaction,
	// it should just panic and exit the program.
	registry := MigrationRegister{
		Tx: nil,
	}

	assert.Panics(t, func() {
		registry.CreateMigrationTable()
	})
}

func TestCreatingDefaultMigrationTableWithErrorExecutingScriptShouldReturnError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectBegin()
	tx, _ := db.Begin()

	registry := MigrationRegister{
		Tx: tx,
	}
	expectedError := errors.New("Error creating table")
	mock.ExpectExec(MIGRATION_TABLE_NAME).WillReturnError(expectedError)

	actualError := registry.CreateMigrationTable()

	assert.NotNil(t, actualError)
	assert.Equal(t, expectedError, actualError)
	assertDatabaseExpectations(t, mock)
}

func TestCreatingDefaultMigrationTableWithSuccessShouldReturnNilError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectBegin()
	tx, _ := db.Begin()

	registry := MigrationRegister{
		Tx: tx,
	}
	mock.ExpectExec(MIGRATION_TABLE_NAME).WillReturnResult(sqlmock.NewResult(0, 1))

	actualError := registry.CreateMigrationTable()

	assert.Nil(t, actualError)
	assertDatabaseExpectations(t, mock)
}

func TestSearchForMigrationNotExecutedShouldReturnIsScriptNotExecuted(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	defer db.Close()
	mock.ExpectBegin()
	tx, _ := db.Begin()

	registry := MigrationRegister{
		Tx: tx,
	}
	script := reader.SQLScript{
		Name:    "script_name.sql",
		Content: "Script content",
	}
	mock.ExpectQuery(regexp.QuoteMeta(script.Name)).
		WillReturnRows(sqlmock.NewRows([]string{"count(id)"}).AddRow(0))

	isAlreadyExecuted, err := registry.IsScriptAlreadyExecuted(script)

	assert.Nil(t, err)
	assert.False(t, isAlreadyExecuted)
	assertDatabaseExpectations(t, mock)
}

func TestSearchForMigrationExecutedShouldReturnIsScriptExecuted(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	defer db.Close()
	mock.ExpectBegin()
	tx, _ := db.Begin()

	registry := MigrationRegister{
		Tx: tx,
	}
	script := reader.SQLScript{
		Name:    "script_name.sql",
		Content: "Script content",
	}
	mock.ExpectQuery(regexp.QuoteMeta(script.Name)).
		WillReturnRows(sqlmock.NewRows([]string{"count(id)"}).AddRow(1))

	isAlreadyExecuted, err := registry.IsScriptAlreadyExecuted(script)

	assert.Nil(t, err)
	assert.True(t, isAlreadyExecuted)
	assertDatabaseExpectations(t, mock)
}

func TestSearchForMigrationExecutedWithErrorOnQueryShouldReturnError(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	defer db.Close()
	mock.ExpectBegin()
	tx, _ := db.Begin()

	registry := MigrationRegister{
		Tx: tx,
	}
	script := reader.SQLScript{
		Name:    "script_name.sql",
		Content: "Script content",
	}
	expectedError := errors.New("Error querying table.")
	mock.ExpectQuery(regexp.QuoteMeta(script.Name)).
		WillReturnError(expectedError)

	isAlreadyExecuted, actualError := registry.IsScriptAlreadyExecuted(script)

	assert.NotNil(t, actualError)
	assert.Equal(t, expectedError, actualError)
	assert.False(t, isAlreadyExecuted)
	assertDatabaseExpectations(t, mock)
}

func TestMarkScriptAsExecutedWithSuccessShouldNotReturnError(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	defer db.Close()
	mock.ExpectBegin()
	tx, _ := db.Begin()

	registry := MigrationRegister{tx}
	script := reader.SQLScript{
		Name:    "script_name.sql",
		Content: "Script content",
	}
	regex := fmt.Sprintf(".*(%s).*(%s)*", MIGRATION_TABLE_NAME, script.Name)
	mock.ExpectExec(regex).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := registry.MarkScriptAsExecuted(script)

	assert.Nil(t, err)
	assertDatabaseExpectations(t, mock)
}

func TestMarkScriptAsExecutedWithErrorShouldReturnError(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	defer db.Close()
	mock.ExpectBegin()
	tx, _ := db.Begin()

	registry := MigrationRegister{tx}
	script := reader.SQLScript{
		Name:    "script_name.sql",
		Content: "Script content",
	}
	expectedError := errors.New("Error executing insert")
	regex := fmt.Sprintf(".*(%s).*(%s)*", MIGRATION_TABLE_NAME, script.Name)
	mock.ExpectExec(regex).
		WillReturnError(expectedError)

	actualError := registry.MarkScriptAsExecuted(script)

	assert.NotNil(t, actualError)
	assert.Equal(t, expectedError, actualError)
	assertDatabaseExpectations(t, mock)
}

func assertDatabaseExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectation were met: %s", err)
	}
}
