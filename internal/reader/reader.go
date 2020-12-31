package reader

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/eaneto/grotto/pkg/database"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

// MigrationReader Basic interface for the migration reader.
type MigrationReader interface {
	ReadScriptFiles() []database.SQLScript
}

// MigrationReaderFS Basic structure for the migration script file system reader.
type MigrationReaderFS struct {
	MigrationDirectory string
}

type ByName []os.FileInfo

func (by ByName) Len() int           { return len(by) }
func (by ByName) Less(i, j int) bool { return by[i].Name() < by[j].Name() }
func (by ByName) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }

// ReadScriptFiles Read all found SQL scripts and return a structure
// with all its content.
func (r MigrationReaderFS) ReadScriptFiles() []database.SQLScript {
	files := getAllScriptFiles(r.MigrationDirectory)

	scripts := make([]database.SQLScript, len(files))
	for index, file := range files {
		content := getFileContent(r.MigrationDirectory, file)
		scripts[index] = database.SQLScript{
			Name:    file.Name(),
			Content: content,
		}
	}
	return scripts
}

// getAllScriptFiles Get all the SQL scripts inside the migration
// directory.
func getAllScriptFiles(migrationDirectory string) []os.FileInfo {
	files, err := ioutil.ReadDir(migrationDirectory)
	if err != nil {
		logrus.Fatal("Error reading migration directory.\n", err)
	}

	if len(files) == 0 {
		logrus.Info("Empty directory, no migrations executed.")
		return nil
	}

	scripts := filterSqlFiles(files)

	// Sort by file name so scripts are executed on order.
	sort.Sort(ByName(scripts))
	return scripts
}

// filterSqlFiles Get all files with .sql extension
func filterSqlFiles(files []os.FileInfo) []os.FileInfo {
	scripts := []os.FileInfo{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			scripts = append(scripts, file)
		}
	}
	return scripts
}

// getFileContent Reads the content from a given file and logs fatal
// if the file is unreadable.
func getFileContent(directory string, file os.FileInfo) string {
	content, err := ioutil.ReadFile(directory + "/" + file.Name())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"file_name": file.Name(),
		}).Fatal("File not found.\n", err)
	}
	return string(content)
}
