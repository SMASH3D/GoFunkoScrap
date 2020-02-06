package main

import(
        _ "github.com/go-sql-driver/mysql"
        "database/sql"
        "fmt"
        "log"
)

func main() {
        cnn, err := sql.Open("mysql", "crawler:popopop@tcp(db:3306)/funkoscrap")
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
