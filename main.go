package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

const DATABASE_URL = "postgres://user:123@localhost:5432/todo"

type ByName []os.FileInfo

func (by ByName) Len() int           { return len(by) }
func (by ByName) Less(i, j int) bool { return by[i].Name() < by[j].Name() }
func (by ByName) Swap(i, j int)      { by[i], by[j] = by[j], by[i] }

func getAllScriptFiles() []os.FileInfo {
	files, err := ioutil.ReadDir("./migration")
	if err != nil {
		log.Fatal("Error reading migration directory.\n", err)
	}

	if len(files) == 0 {
		logrus.Info("Empty directory.\n")
		return nil
	}

	// Get all files with .sql extension
	scripts := []os.FileInfo{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			scripts = append(scripts, file)
		}
	}

	// Sort by file name so scripts are executed on order.
	sort.Sort(ByName(scripts))
	return scripts
}

func readScriptFiles() []string {
	files := getAllScriptFiles()

	file_content := make([]string, len(files))
	for _, file := range files {
		content, err := ioutil.ReadFile("migration/" + file.Name())
		if err != nil {
			logrus.Fatal("File not found.\n", err)
		}
		file_content = append(file_content, string(content))
	}
	return file_content
}

func executeScript(db *sql.DB, script string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Error starting transaction\n", err)
	}
	_, err = db.Exec(script)
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func main() {
	scripts := readScriptFiles()
	for _, script := range scripts {
		fmt.Print(script)
	}
	db, err := sql.Open("pgx", DATABASE_URL)
	if err != nil {
		log.Fatal("Connection failed\n", err)
	}
	defer db.Close()
	for _, script := range scripts {
		executeScript(db, script)
	}
}
