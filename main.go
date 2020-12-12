package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "user:123@tcp(localhost:5432)/todo")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM todo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	fmt.Println(rows)
}
