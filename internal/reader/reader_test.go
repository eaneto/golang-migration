package reader

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestReadNonExistentDirectoryShouldLogFatal(t *testing.T) {
	reader := MigrationReaderFS{
		MigrationDirectory: "non_existent_directory",
	}

	defer func() { logrus.StandardLogger().ExitFunc = nil }()
	var fatal bool
	logrus.StandardLogger().ExitFunc = func(int) { fatal = true }

	reader.ReadScriptFiles()

	assert.True(t, fatal)
}

func TestReadEmptyDirectoryShouldReturnEmptySlice(t *testing.T) {
	dir := "reader_test"
	os.RemoveAll(dir)
	os.Mkdir(dir, os.ModePerm)

	reader := MigrationReaderFS{
		MigrationDirectory: dir,
	}

	scripts := reader.ReadScriptFiles()

	assert.Empty(t, scripts)
}

func TestReadDirectoryWithFilesButNoSqlFilesShouldReturnEmptySlice(t *testing.T) {
	dir := "reader_test"
	os.RemoveAll(dir)
	os.Mkdir(dir, os.ModePerm)
	data := []byte("data")
	ioutil.WriteFile(dir+"/file.txt", data, os.ModePerm)

	reader := MigrationReaderFS{
		MigrationDirectory: dir,
	}

	scripts := reader.ReadScriptFiles()

	assert.Empty(t, scripts)
}

func TestReadDirectoryWithOneSqlFileShouldReturnListWithFileContentAndName(t *testing.T) {
	dir := "reader_test"
	os.RemoveAll(dir)
	os.Mkdir(dir, os.ModePerm)
	filename := "file_name.sql"
	content := []byte("data")
	ioutil.WriteFile(dir+"/"+filename, content, os.ModePerm)

	reader := MigrationReaderFS{
		MigrationDirectory: dir,
	}

	scripts := reader.ReadScriptFiles()

	assert.NotEmpty(t, scripts)
	assert.Equal(t, filename, scripts[0].Name)
	assert.Equal(t, string(content), scripts[0].Content)
}

func TestReadDirectoryWithOneSqlFileWithOutPermissionShouldLogFatal(t *testing.T) {
	dir := "reader_test"
	os.RemoveAll(dir)
	os.Mkdir(dir, os.ModePerm)
	filename := "file_name.sql"
	content := []byte("data")
	ioutil.WriteFile(dir+"/"+filename, content, 0000)

	reader := MigrationReaderFS{
		MigrationDirectory: dir,
	}

	defer func() { logrus.StandardLogger().ExitFunc = nil }()
	var fatal bool
	logrus.StandardLogger().ExitFunc = func(int) { fatal = true }

	reader.ReadScriptFiles()

	assert.True(t, fatal)
}

func TestReadDirectoryWithMultipleSqlFilesShouldReturnListWithAllFilesButNoneThatAreNotSql(t *testing.T) {
	dir := "reader_test"
	os.RemoveAll(dir)
	os.Mkdir(dir, os.ModePerm)
	files := []string{
		"file_name.sql",
		"file_2.sql",
		"not_sql.txt",
	}
	contents := [][]byte{
		[]byte("data1"),
		[]byte("data2"),
		[]byte("data3"),
	}
	expectedScriptsSize := 2

	for index := range files {
		ioutil.WriteFile(dir+"/"+files[index], contents[index], os.ModePerm)
	}

	reader := MigrationReaderFS{
		MigrationDirectory: dir,
	}

	scripts := reader.ReadScriptFiles()

	assert.NotEmpty(t, scripts)
	assert.Equal(t, expectedScriptsSize, len(scripts))

	// one of the created scripts has a number on it, so it will come
	// first, it's important to validate that this script is the first
	// one returned.
	assert.Equal(t, files[1], scripts[0].Name)
	assert.Equal(t, string(contents[1]), scripts[0].Content)

	assert.Equal(t, files[0], scripts[1].Name)
	assert.Equal(t, string(contents[0]), scripts[1].Content)
}
