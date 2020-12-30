package executor

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eaneto/grotto/pkg/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MigrationRegisterMock struct {
	mock.Mock
}

func (m *MigrationRegisterMock) CreateMigrationTable() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MigrationRegisterMock) IsScriptAlreadyExecuted(script database.SQLScript) (bool, error) {
	args := m.Called(script)
	return args.Bool(0), args.Error(1)
}

func (m *MigrationRegisterMock) MarkScriptAsExecuted(script database.SQLScript) error {
	args := m.Called()
	return args.Error(0)
}

func TestProcessScriptWithNilTransactionShouldPanic(t *testing.T) {
	migrationRegister := new(MigrationRegisterMock)
	scriptExecutor := ScriptExecutorSQL{
		Tx:                nil,
		MigrationRegister: migrationRegister,
	}

	scripts := []database.SQLScript{
		{
			Name:    "script_name.sql",
			Content: "INSERT INTO USERS VALUES ('id')",
		},
	}
	assert.Panics(t, func() {
		scriptExecutor.ProcessScripts(scripts)
	})
	migrationRegister.AssertExpectations(t)
}

func TestProcessScriptWithEmptyListShouldNotReturnErrorAndDoNothing(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	mock.ExpectBegin()
	tx, _ := db.Begin()

	migrationRegister := new(MigrationRegisterMock)
	scriptExecutor := ScriptExecutorSQL{
		Tx:                tx,
		MigrationRegister: migrationRegister,
	}

	scripts := []database.SQLScript{}

	error := scriptExecutor.ProcessScripts(scripts)

	assert.Nil(t, error)
	migrationRegister.AssertExpectations(t)
	assertDatabaseExpectations(t, mock)
}

func TestProcessOneScriptWithErrorCheckingIfTheScriptWasExecutedShouldReturnError(t *testing.T) {
	db, dbMock, _ := sqlmock.New()
	defer db.Close()
	dbMock.ExpectBegin()
	tx, _ := db.Begin()

	migrationRegister := new(MigrationRegisterMock)
	expectedError := errors.New("Error checking if the script was executed")
	migrationRegister.On("IsScriptAlreadyExecuted", mock.Anything).
		Return(false, expectedError)

	scriptExecutor := ScriptExecutorSQL{
		Tx:                tx,
		MigrationRegister: migrationRegister,
	}

	scripts := []database.SQLScript{
		{
			Name:    "script_name.sql",
			Content: "INSERT INTO USERS VALUES ('id')",
		},
	}
	actualError := scriptExecutor.ProcessScripts(scripts)

	assert.NotNil(t, actualError)
	assert.Equal(t, expectedError, actualError)
	migrationRegister.AssertExpectations(t)
	assertDatabaseExpectations(t, dbMock)
}

func TestProcessOneScriptWithAlreadyProcessedScriptShouldNotReturnError(t *testing.T) {
	db, dbMock, _ := sqlmock.New()
	defer db.Close()
	dbMock.ExpectBegin()
	tx, _ := db.Begin()

	migrationRegister := new(MigrationRegisterMock)
	migrationRegister.On("IsScriptAlreadyExecuted", mock.Anything).Return(true, nil)

	scriptExecutor := ScriptExecutorSQL{
		Tx:                tx,
		MigrationRegister: migrationRegister,
	}

	scripts := []database.SQLScript{
		{
			Name:    "script_name.sql",
			Content: "INSERT INTO USERS VALUES ('id')",
		},
	}
	error := scriptExecutor.ProcessScripts(scripts)

	assert.Nil(t, error)
	migrationRegister.AssertExpectations(t)
	assertDatabaseExpectations(t, dbMock)
}

func TestProcessOneUnexecutedScriptAndErrorMarkingAsExecutedShouldReturnError(t *testing.T) {
	db, dbMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	defer db.Close()
	dbMock.ExpectBegin()
	tx, _ := db.Begin()

	migrationRegister := new(MigrationRegisterMock)
	migrationRegister.On("IsScriptAlreadyExecuted", mock.Anything).Return(false, nil)
	expectedError := errors.New("Error marking as executed")
	migrationRegister.On("MarkScriptAsExecuted", mock.Anything).Return(expectedError)

	scriptExecutor := ScriptExecutorSQL{
		Tx:                tx,
		MigrationRegister: migrationRegister,
	}

	scripts := []database.SQLScript{
		{
			Name:    "script_name.sql",
			Content: "INSERT INTO USERS VALUES ('id')",
		},
	}
	dbMock.ExpectExec(regexp.QuoteMeta(scripts[0].Content)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	actualError := scriptExecutor.ProcessScripts(scripts)

	assert.NotNil(t, actualError)
	assert.Equal(t, expectedError, actualError)
	migrationRegister.AssertExpectations(t)
	assertDatabaseExpectations(t, dbMock)
}

func TestProcessOneUnexecutedScriptShouldExecuteScriptContentAndNotReturnError(t *testing.T) {
	db, dbMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	defer db.Close()
	dbMock.ExpectBegin()
	tx, _ := db.Begin()

	migrationRegister := new(MigrationRegisterMock)
	migrationRegister.On("IsScriptAlreadyExecuted", mock.Anything).Return(false, nil)
	migrationRegister.On("MarkScriptAsExecuted", mock.Anything).Return(nil)

	scriptExecutor := ScriptExecutorSQL{
		Tx:                tx,
		MigrationRegister: migrationRegister,
	}

	scripts := []database.SQLScript{
		{
			Name:    "script_name.sql",
			Content: "INSERT INTO USERS VALUES ('id')",
		},
	}
	dbMock.ExpectExec(regexp.QuoteMeta(scripts[0].Content)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	error := scriptExecutor.ProcessScripts(scripts)

	assert.Nil(t, error)
	migrationRegister.AssertExpectations(t)
	assertDatabaseExpectations(t, dbMock)
}

func TestProcessTwoUnexecutedScriptShouldExecuteScriptContentAndNotReturnError(t *testing.T) {
	db, dbMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	defer db.Close()
	dbMock.ExpectBegin()
	tx, _ := db.Begin()

	migrationRegister := new(MigrationRegisterMock)
	migrationRegister.On("IsScriptAlreadyExecuted", mock.Anything).Return(false, nil)
	migrationRegister.On("MarkScriptAsExecuted", mock.Anything).Return(nil)

	scriptExecutor := ScriptExecutorSQL{
		Tx:                tx,
		MigrationRegister: migrationRegister,
	}

	scripts := []database.SQLScript{
		{
			Name:    "script_name.sql",
			Content: "INSERT INTO USERS VALUES ('id')",
		},
		{
			Name:    "script_name_2.sql",
			Content: "UPDATE TABLE SET NAME = ''",
		},
	}
	for _, script := range scripts {
		dbMock.ExpectExec(regexp.QuoteMeta(script.Content)).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	error := scriptExecutor.ProcessScripts(scripts)

	assert.Nil(t, error)
	migrationRegister.AssertExpectations(t)
	assertDatabaseExpectations(t, dbMock)
}

func TestProcessOneExecutedAndOneUnexecutedScriptShouldExecuteOneScriptContentAndNotReturnError(t *testing.T) {
	db, dbMock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	defer db.Close()
	dbMock.ExpectBegin()
	tx, _ := db.Begin()

	migrationRegister := new(MigrationRegisterMock)
	migrationRegister.On("MarkScriptAsExecuted", mock.Anything).Return(nil)

	scriptExecutor := ScriptExecutorSQL{
		Tx:                tx,
		MigrationRegister: migrationRegister,
	}

	scripts := []database.SQLScript{
		{
			Name:    "script_name.sql",
			Content: "INSERT INTO USERS VALUES ('id')",
		},
		{
			Name:    "script_name_2.sql",
			Content: "UPDATE TABLE SET NAME = 'new_name'",
		},
	}
	migrationRegister.On("IsScriptAlreadyExecuted", scripts[0]).Return(true, nil)
	migrationRegister.On("IsScriptAlreadyExecuted", scripts[1]).Return(false, nil)
	dbMock.ExpectExec(regexp.QuoteMeta(scripts[1].Content)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	error := scriptExecutor.ProcessScripts(scripts)

	assert.Nil(t, error)
	migrationRegister.AssertExpectations(t)
	assertDatabaseExpectations(t, dbMock)
}

func assertDatabaseExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectation were met: %s", err)
	}
}
