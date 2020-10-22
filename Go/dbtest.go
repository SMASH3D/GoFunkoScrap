package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cnn, err := sql.Open("mysql", "crawler:popopop@tcp(db:3307)/funkoscrap")
	if err != nil {
		log.Fatal(err)
	}

	id := 1
	var name string

	if err := cnn.QueryRow("SELECT name FROM test_tb WHERE id = ? LIMIT 1", id).Scan(&name); err != nil {
		log.Fatal(err)
	}

	fmt.Println(id, name)
}
