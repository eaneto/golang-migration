package reader

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

// SQLScript Represents a SQL script with the filename and content.
type SQLScript struct {
	content string
	name    string
}

type ByName []os.FileInfo

func (by ByName) Len() int           { return len(by) }
func (by ByName) Less(i, j int) bool { return by[i].Name() < by[j].Name() }
func (by ByName) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }

// ReadScriptFiles Read all found SQL scripts and return a structure with
// all its content.
func ReadScriptFiles() []SQLScript {
	files := getAllScriptFiles()

	file_content := make([]SQLScript, len(files))
	for _, file := range files {
		content, err := ioutil.ReadFile("migration/" + file.Name())
		if err != nil {
			logrus.Fatal("File not found.\n", err)
		}
		data := SQLScript{
			name:    file.Name(),
			content: string(content),
		}
		file_content = append(file_content, data)
	}
	return file_content
}

// Get all files with .sql extension
func filterSqlFiles(files []os.FileInfo) []os.FileInfo {
	scripts := []os.FileInfo{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			scripts = append(scripts, file)
		}
	}
	return scripts
}

// getAllScriptFiles Get all the SQL scripts inside the migration directory.
func getAllScriptFiles() []os.FileInfo {
	files, err := ioutil.ReadDir("./migration")
	if err != nil {
		log.Fatal("Error reading migration directory.\n", err)
	}

	if len(files) == 0 {
		logrus.Info("Empty directory.\n")
		return nil
	}

	scripts := filterSqlFiles(files)

	// Sort by file name so scripts are executed on order.
	sort.Sort(ByName(scripts))
	return scripts
}

// ExecuteScript Executes a given SQL script.
// Every script must be executed inside a transaction.
func ExecuteScript(db *sql.DB, script SQLScript) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction\n", err)
	}
	logrus.Info("Executing script: ", script.name)
	_, err = db.Exec(script.content)
	if err != nil {
		tx.Rollback()
		log.Fatal("Error executing script.\n", err)
	} else {
		tx.Commit()
	}
}
