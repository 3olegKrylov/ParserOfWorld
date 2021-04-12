package db

import (
	"database/sql"
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

//Открывает базу данных и проверяет работоспособность
func DBconnect() *sql.DB{
	connect, err := sql.Open("clickhouse", "tcp://127.0.0.1:9000?debug=true")

	if err != nil {
		log.Fatal(err)
	}

	if err := connect.Ping(); err != nil {
		log.Fatal(err)
	}

	return connect
}


func DBinit(connect *sql.DB) {
	if _, err := connect.Exec(`
		CREATE TABLE IF NOT EXISTS default.Users(
			id Nullable(Int32),
			Title Nullable(String),
			SubTitle Nullable(String),
			Comment Nullable(String)
		) engine=Memory
	`); err != nil{
		log.Fatal(err)
	}
}

func DBAddUser(ID int32, Title string, SubTitle string, Comment string, connect *sql.DB){
	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO default.Users (id, Title, SubTitle, Comment) VAlUES (?, ?, ?, ?)")
	)

	defer stmt.Close()

	if _, err := stmt.Exec(
		ID,
		Title,
		SubTitle,
		Comment,
	); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}




}
