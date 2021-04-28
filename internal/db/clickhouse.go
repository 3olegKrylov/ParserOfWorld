package db

import (
	"database/sql"
	"github.com/testSpace/model"
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
func DBconnect() *sql.DB {
	connect, err := sql.Open("clickhouse", "tcp://localhost:9000?debug=true")
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
			Id 			  	  Nullable(Int32),
			LinkAccount	      Nullable(String),    
			Title 		      Nullable(String),
			SubTitle 		  Nullable(String),
			Comment 		  Nullable(String), 
			Mail   			  Nullable(String),
			Telegram          Nullable(String),
			Instagram         Nullable(String),
			Links             Nullable(String),
			LanguageAccount   Nullable(String),
			Phone			  Nullable(String),
			Following         Nullable(Int32),
			Followers         Nullable(Int32),
			Likes             Nullable(Int32),
			LastPostShowTotal Nullable(Int32),
			AverageShows      Nullable(Int32),
			MedianShows       Nullable(Int32),
			TotalPosts        Nullable(Int32),
			LastActionTime    Nullable(DateTime),
			ParsingTime       Nullable(DateTime)
		) engine=Memory`); err != nil {
		log.Print("DB is no Iinit", err)
		return
	}

}

func DBAddUser(user model.UserData, connect *sql.DB) {
	var (
		tx, _   = connect.Begin()
		stmt, _ = tx.Prepare("INSERT INTO default.Users (Id,LinkAccount,Title,SubTitle,Comment,Mail,Telegram,Instagram,Links,LanguageAccount,Phone,Following,Followers,Likes,LastPostShowTotal,AverageShows,MedianShows,TotalPosts,LastActionTime,ParsingTime) VAlUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	)

	defer stmt.Close()

	if _, err := stmt.Exec(
		user.Id,
		user.LinkAccount,
		user.Title,
		user.SubTitle,
		user.Comment,
		user.Mail,
		user.Telegram,
		user.Instagram,
		user.Links,
		user.LanguageAccount,
		user.Phone,
		user.Following,
		user.Followers,
		user.Likes,
		user.LastPostShowTotal,
		user.AverageShows,
		user.MedianShows,
		user.TotalPosts,
		user.LastActionTime,
		user.ParsingTime,
	); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

}

func InitUsers(connect *sql.DB) map[string]int32 {
	usersMap := make(map[string]int32)

	rows, err := connect.Query("SELECT Id, Title FROM default.Users")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			id    int32
			title string
		)
		if err := rows.Scan(&id, &title); err != nil {
			log.Fatal(err)
		}

		_, ok := usersMap[title]

		if ok {
			log.Println("user повторяется")
		} else {
			usersMap[title] = id
		}

	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return usersMap
}
