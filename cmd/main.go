package main

import (
	"fmt"
	internal "github.com/testSpace/internal/data/send"
	"github.com/testSpace/internal/db"
	"github.com/testSpace/internal/parsing"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

func main() {
	urlStr := []string{
		//"https://www.tiktok.com/search?q=нутик&lang=ru-RU",
		//"https://www.tiktok.com/search?q=анастэйша&lang=ru-RU",
		//"https://www.tiktok.com/search?q=маша&lang=ru-RU",
		//"https://www.tiktok.com/search?q=машуня&lang=ru-RU",
		//"https://www.tiktok.com/search?q=мария&lang=ru-RU",
		//"https://www.tiktok.com/search?q=мила&lang=ru-RU",
		//"https://www.tiktok.com/search?q=милана&lang=ru-RU",
		//"https://www.tiktok.com/search?q=подняла&lang=ru-RU",
		//"https://www.tiktok.com/search?q=насосала&lang=ru-RU",
		//"https://www.tiktok.com/search?q=инстасамка&lang=ru-RU",
		//"https://www.tiktok.com/search?q=блогерша&lang=ru-RU",
		//"https://www.tiktok.com/search?q=милиардерша&lang=ru-RU",
		//"https://www.tiktok.com/search?q=миллионерша&lang=ru-RU",
		//"https://www.tiktok.com/search?q=русская&lang=ru-RU",
		//"https://www.tiktok.com/search?q=ксюша&lang=ru-RU",
		//"https://www.tiktok.com/search?q=ксения&lang=ru-RU",
		//"https://www.tiktok.com/search?q=инст&lang=ru-RU",
		//"https://www.tiktok.com/search?q=инста&lang=ru-RU",
		"https://www.tiktok.com/search?q=катя&lang=ru-RU",
		"https://www.tiktok.com/search?q=катюша&lang=ru-RU",
		"https://www.tiktok.com/search?q=екатерина&lang=ru-RU",
		"https://www.tiktok.com/search?q=люба&lang=ru-RU",
		"https://www.tiktok.com/search?q=вероника&lang=ru-RU",
		"https://www.tiktok.com/search?q=клава&lang=ru-RU",
		"https://www.tiktok.com/search?q=вера&lang=ru-RU",
		"https://www.tiktok.com/search?q=ника&lang=ru-RU",
		"https://www.tiktok.com/search?q=аня&lang=ru-RU",
		"https://www.tiktok.com/search?q=анна&lang=ru-RU",
		"https://www.tiktok.com/search?q=таня&lang=ru-RU",
		"https://www.tiktok.com/search?q=оксана&lang=ru-RU",
		"https://www.tiktok.com/search?q=инста&lang=ru-RU",
		"https://www.tiktok.com/search?q=нюша&lang=ru-RU",
		"https://www.tiktok.com/search?q=дина&lang=ru-RU",
		"https://www.tiktok.com/search?q=лера&lang=ru-RU",
		"https://www.tiktok.com/search?q=нюра&lang=ru-RU",
		"https://www.tiktok.com/search?q=алёна&lang=ru-RU",
		"https://www.tiktok.com/search?q=татьяна&lang=ru-RU",

	}

	//подключение к clickhouse
	dbConnect := db.DBconnect()
	db.DBinit(dbConnect)
	defer dbConnect.Close()

	start := time.Now()
	countOfUsers := int32(0)
	//парсинг аккаунтов
	for i := 0; i < len(urlStr); i++ {
		//инициализация карты пользователей
		userMap := db.InitUsers(dbConnect)
		countOfUsers = int32(len(userMap))
		text := parsing.ParseFindList(urlStr[i])

		lines := strings.Split(text, "\n\n")

		var newUsers []string

		for num, value := range lines {
			if strings.HasSuffix(value, "Подписчики") {
				_, ok := userMap[strings.TrimSpace(lines[num-1])]
				if !ok {
					newUsers = append(newUsers, lines[num-1])
				}
			}
		}
		if len(newUsers) > 0 {
			internal.SendData(newUsers, dbConnect, &countOfUsers)
			countOfUsers = countOfUsers + int32(len(newUsers))
			fmt.Println("Время работы: ", time.Since(start))
		}
	}

	fmt.Println("Время: ", time.Since(start), " \nКол-во юзеров: ", countOfUsers)

}
