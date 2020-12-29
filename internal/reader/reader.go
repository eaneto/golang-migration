package reader

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

// MigrationReader Basic structure for the migration script reader.
type MigrationReader struct {
	MigrationDirectory string
}

// SQLScript Represents a SQL script with the filename and content.
type SQLScript struct {
	// The actual SQL script content.
	Content string
	// The script filename with the .sql extension.
	Name string
}

type ByName []os.FileInfo

func (by ByName) Len() int           { return len(by) }
func (by ByName) Less(i, j int) bool { return by[i].Name() < by[j].Name() }
func (by ByName) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }

// ReadScriptFiles Read all found SQL scripts and return a structure with
// all its content.
func (r MigrationReader) ReadScriptFiles() []SQLScript {
	files := getAllScriptFiles(r.MigrationDirectory)

	scripts := make([]SQLScript, len(files))
	for index, file := range files {
		content, err := ioutil.ReadFile(r.MigrationDirectory + "/" + file.Name())
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"file_name": file.Name(),
			}).Fatal("File not found.\n", err)
		}
		script := SQLScript{
			Name:    file.Name(),
			Content: string(content),
		}
		scripts[index] = script
	}
	return scripts
}

// getAllScriptFiles Get all the SQL scripts inside the migration directory.
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
