package main

import (
	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go"
	"log"
)

type User struct {
	Id        int32
	Title     string
	SubTitle  string
	Comment   string
	Following int32
	Followers int32
	Likes     int32
	Linkes    string
}

func main() {
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")

	if err != nil {
		log.Fatal(err)
	}

	if err := connect.Ping(); err != nil {
		log.Fatal(err)
	}

	_, err = connect.Exec(`
		CREATE TABLE IF NOT EXISTS default.Users(
			id Nullable(Int32),
			Title Nullable(String),
			SubTitle Nullable(String),
			Comment Nullable(String)
		) engine=Memory
	`)

	if err != nil {
		log.Fatal(err)
	}

	//newUser := User{
	//	3,
	//	"Иван",
	//	"Ментов",
	//	"комментарий так себе",
	//	5,
	//	8,
	//	3,
	//	"https://clickhouse.tech/docs/ru/sql-reference/statements/insert-into/",
	//}

	//var (
	//	tx, _   = connect.Begin()
	//	stmt, _ = tx.Prepare("INSERT INTO default.Users (id, Title, SubTitle, Comment) VAlUES (?, ?, ?, ?)")
	//)
	//
	//defer stmt.Close()
	//
	//if res, err := stmt.Exec(
	//	newUser.Id,
	//	newUser.Title,
	//	newUser.SubTitle,
	//	newUser.Comment,
	//); err != nil {
	//	log.Fatal(err)
	//} else {
	//	fmt.Println(res)
	//}
	//
	//if err := tx.Commit(); err != nil {
	//	log.Fatal(err)
	//}

	rows, err := connect.Query("SELECT Title FROM default.Users WHERE id = 1")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			title string
		)
		if err := rows.Scan(&title); err != nil {
			log.Fatal(err)
		}
		log.Printf("title: %s", title)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}
