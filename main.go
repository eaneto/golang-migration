package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/eaneto/golang-migration/reader"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

const DATABASE_URL = "postgres://user:123@localhost:5432/"

// TODO Create a table to registered processed scripts.
func main() {
	databaseName := os.Args[1]
	db, err := sql.Open("pgx", DATABASE_URL+databaseName)
	if err != nil {
		log.Fatal("Connection failed\n", err)
	}
	defer db.Close()
	scripts := reader.ReadScriptFiles()
	for _, script := range scripts {
		reader.ExecuteScript(db, script)
	}
	logrus.Info("Migration executed successfully!")
}
