package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	arr := []int{-1, -2, -3}

	fmt.Println(findValue(arr, 3))

}

func findValue(arr []int, res int) (int, int, error) {
	if len(arr) < 2 {
		return -1, -1, fmt.Errorf("min len")
	}

	for i := 0; i < len(arr); i++ {
		for j := 1; j < len(arr); j++ {
			if (arr[i] + arr[j]) == res {
				return arr[i], arr[j], nil
			}
		}

	}
	return -1, -1, fmt.Errorf("coudnt find")
}

//	db, err := sql.Open("mysql", "")
//
//	if err != nil {
//		panic(err)
//	}
//	defer db.Close()
//
//	res, err := db.Exec("INSERT INTO Users (LinkAccount,Title,SubTitle,Comment,Mail,Telegram,Instagram,Links,LanguageAccount,Phone,Following,Followers,Likes,LastPostShowTotal,AverageShows,MedianShows,TotalPosts,LastActionTime,ParsingTime) VAlUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
//		"https://www.tiktok.com/@oleg_c", "oleg_c", "Oleg🌐", "Licensed adrenaline addict Snap: gelo737", "", "", "", "", "", "", 70, 422000, 10699999, 9676, 0, 0, 0, "2021-03-31 19:00:00", "2021-04-21 10:29:43")
//	if err != nil {
//		panic(err)
//	}
//
//	//res, err := db.Exec("INSERT INTO users (name, secondname) VALUES(?,?)", "Гена", "Малышев")
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	fmt.Println(res)
//}
