package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/testSpace/model"
	"log"
)

var Connect *sql.DB

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
func DBconnect() {
	connect, err := sql.Open("mysql", "oleg:1@tcp(65.21.53.188:3306)/tiktok")
	if err != nil {
		log.Fatal(err)
	}

	if err := connect.Ping(); err != nil {
		log.Fatal(err)
	}
	Connect = connect
}

func DBinit() {
	//if _, err := Connect.Exec(`
	//	CREATE TABLE IF NOT EXISTS tiktok.Users(
	//		'Id' 			  	  Nullable(Int32),
	//		LinkAccount	      Nullable(String),
	//		Title 		      Nullable(String),
	//		SubTitle 		  Nullable(String),
	//		Comment 		  Nullable(String),
	//		Mail   			  Nullable(String),
	//		Telegram          Nullable(String),
	//		Instagram         Nullable(String),
	//		Links             Nullable(String),
	//		LanguageAccount   Nullable(String),
	//		Phone			  Nullable(String),
	//		Following         Nullable(Int32),
	//		Followers         Nullable(Int32),
	//		Likes             Nullable(Int32),
	//		LastPostShowTotal Nullable(Int32),
	//		AverageShows      Nullable(Int32),
	//		MedianShows       Nullable(Int32),
	//		TotalPosts        Nullable(Int32),
	//		LastActionTime    Nullable(DateTime),
	//		ParsingTime       Nullable(DateTime)
	//	) engine=Memory`); err != nil {
	//	log.Print("DB is no Iinit", err)
	//
	if _, err := Connect.Exec(`
		CREATE TABLE IF NOT EXISTS Users
(
    Id                INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    LinkAccount       TEXT               NOT NULL,
    Title             CHAR(50)           NOT NULL UNIQUE,
    SubTitle          TEXT               NULL,
    Comment           TEXT               NULL,
    Mail              TEXT               NULL,
    Telegram          TEXT               NULL,
    Instagram         TEXT               NULL,
    Links             TEXT               NULL,
    LanguageAccount   TEXT               NULL,
    Phone             TEXT               NULL,
    Following         INT                NULL,
    Followers         INT                NULL,
    Likes             INT                NULL,
    LastPostShowTotal INT                NULL,
    AverageShows      INT                NULL,
    MedianShows       INT                NULL,
    TotalPosts        INT                NULL,
    LastActionTime    DATETIME           NULL,
    ParsingTime       DATETIME           NOT NULL
)`); err != nil {
		log.Print("DB is no Iinit", err)
		return
	}

}

func DBAddUser(user model.UserData) {

	var (
		Tx, _   = Connect.Begin()
		stmt, _ = Tx.Prepare("INSERT INTO tiktok.Users (LinkAccount,Title,SubTitle,Comment,Mail,Telegram,Instagram,Links,LanguageAccount,Phone,Following,Followers,Likes,LastPostShowTotal,AverageShows,MedianShows,TotalPosts,LastActionTime,ParsingTime) VAlUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	)

	defer stmt.Close()

	if user.LastActionTime.IsZero() {
		if _, err := stmt.Exec(
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
			nil,
			user.ParsingTime,
		); err != nil {
			log.Println(err)
			return
		}
	} else if _, err := stmt.Exec(
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
		log.Println(err)
		return
	}

	err := Tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func InitUsers() int32 {
	usersMap := make(map[string]int32)

	rows, err := Connect.Query("SELECT Id, Title FROM tiktok.Users ")
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

	return int32(len(usersMap))
}

//возвращает true если существует данный user иначе false
func FindUserDB(nick string) bool {

	rows, err := Connect.Query("SELECT Title FROM tiktok.Users WHERE Title = ?", nick)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			comment string
		)
		if err := rows.Scan(&comment); err != nil {
			log.Fatal(err)
		}
		if comment == nick {
			return true
		}
	}

	return false

}
