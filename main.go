package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const DATABASE_URL = "postgres://user:123@localhost:5432/todo"

func main() {
	db, err := sql.Open("pgx", DATABASE_URL)
	if err != nil {
		log.Fatal("Connection failed\n", err)
	}
	defer db.Close()
	var greeting string
	err = db.QueryRow("select 'Hello, World!'").Scan(&greeting)
	if err != nil {
		log.Fatal("QueryRow failed\n", err)
	}
	fmt.Println(greeting)
}
