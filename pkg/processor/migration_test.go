package processor

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/eaneto/grotto/pkg/database"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ScriptExecutorMock struct {
	mock.Mock
}

func (m *ScriptExecutorMock) CreateMigrationTable() error {
	args := m.Called()
	return args.Error(0)
}

func (m *ScriptExecutorMock) ProcessScripts(scripts []database.SQLScript) error {
	args := m.Called()
	return args.Error(0)
}

func (m *ScriptExecutorMock) RollbackTransaction() {
	m.Called()
}

func (m *ScriptExecutorMock) CommitTransaction() {
	m.Called()
}

type ReaderMock struct {
	mock.Mock
}

func (m *ReaderMock) ReadScriptFiles() []database.SQLScript {
	args := m.Called()
	return args.Get(0).([]database.SQLScript)
}

func TestProcessingWithNoScriptsReturnedByReaderShouldCommitTransaction(t *testing.T) {
	executorMock := new(ScriptExecutorMock)
	readerMock := new(ReaderMock)

	executorMock.On("CreateMigrationTable").Return(nil)
	executorMock.On("ProcessScripts", mock.Anything).Return(nil)
	executorMock.On("CommitTransaction").Return(nil)
	readerMock.On("ReadScriptFiles").Return([]database.SQLScript{})

	processor := MigrationProcessorSQL{
		Executor: executorMock,
		Reader:   readerMock,
	}

	processor.ProcessMigration()

	executorMock.AssertExpectations(t)
	executorMock.AssertNotCalled(t, "RollbackTransaction")
	readerMock.AssertExpectations(t)
}

func TestProcessingWithErrorShouldRollbackTransaction(t *testing.T) {
	executorMock := new(ScriptExecutorMock)
	readerMock := new(ReaderMock)

	executorMock.On("CreateMigrationTable").Return(nil)
	executorMock.On("ProcessScripts", mock.Anything).Return(errors.New(""))
	executorMock.On("RollbackTransaction").Return(nil)
	readerMock.On("ReadScriptFiles").Return([]database.SQLScript{})

	processor := MigrationProcessorSQL{
		Executor: executorMock,
		Reader:   readerMock,
	}

	processor.ProcessMigration()

	executorMock.AssertExpectations(t)
	executorMock.AssertNotCalled(t, "CommitTransaction")
	readerMock.AssertExpectations(t)
}

func TestProcessingWithErrorCreatingMigrationTableShouldRollbackTransaction(t *testing.T) {
	executorMock := new(ScriptExecutorMock)
	readerMock := new(ReaderMock)

	executorMock.On("CreateMigrationTable").Return(errors.New(""))
	executorMock.On("RollbackTransaction").Panic("Panic")

	processor := MigrationProcessorSQL{
		Executor: executorMock,
		Reader:   readerMock,
	}

	assert.Panics(t, processor.ProcessMigration)

	executorMock.AssertExpectations(t)
	executorMock.AssertNotCalled(t, "CommitTransaction")
	executorMock.AssertNotCalled(t, "ProcessScript", mock.Anything)
	readerMock.AssertNotCalled(t, "ReadScriptFiles")
}

func TestInitializeExecutorWithSuccessShouldBeginDatabaseConnection(t *testing.T) {
	db, dbMock, _ := sqlmock.New()
	defer db.Close()
	dbMock.ExpectBegin()

	initializeExecutor(db)

	assertDatabaseExpectations(t, dbMock)
}

func TestInitializeExecutorWithErrorShouldPanic(t *testing.T) {
	db, dbMock, _ := sqlmock.New()
	defer db.Close()
	dbMock.ExpectBegin().WillReturnError(errors.New("Error"))

	// Changes the logger for a mock that changes the value of fatal
	// to true if the log was fatal.
	defer func() { logrus.StandardLogger().ExitFunc = nil }()
	var fatal bool
	logrus.StandardLogger().ExitFunc = func(int) { fatal = true }

	initializeExecutor(db)

	assert.True(t, fatal)
	assertDatabaseExpectations(t, dbMock)
}

func assertDatabaseExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not all expectation were met: %s", err)
	}
}
