package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// ...
func main() {
	db, err := sql.Open("mysql", "root:root@tcp(65.21.53.188:3306)/tiktok")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	insert, err := db.Query("INSERT INTO `users` (`name`, `age`) VALUE ()")
	fmt.Println("подключено к mysql")
}
