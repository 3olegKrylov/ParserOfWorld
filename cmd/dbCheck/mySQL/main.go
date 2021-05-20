package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "oleg:1@tcp(65.21.53.188:3306)/tiktok")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO Users (LinkAccount,Title,SubTitle,Comment,Mail,Telegram,Instagram,Links,LanguageAccount,Phone,Following,Followers,Likes,LastPostShowTotal,AverageShows,MedianShows,TotalPosts,LastActionTime,ParsingTime) VAlUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"https://www.tiktok.com/@oleg_c", "oleg_c", "Oleg🌐", "Licensed adrenaline addict Snap: gelo737", "", "", "", "", "", "", 70, 422000, 10699999, 9676, 0, 0, 0, "2021-03-31 19:00:00", "2021-04-21 10:29:43")
	if err != nil {
		panic(err)
	}

	//res, err := db.Exec("INSERT INTO users (name, secondname) VALUES(?,?)", "Гена", "Малышев")
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println(res)
}
