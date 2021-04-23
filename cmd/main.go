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
		"https://www.tiktok.com/search?q=настюля&lang=ru-RU",
		"https://www.tiktok.com/search?q=настя&lang=ru-RU",
		"https://www.tiktok.com/search?q=анастасия&lang=ru-RU",
		"https://www.tiktok.com/search?q=настюня&lang=ru-RU",
		"https://www.tiktok.com/search?q=нутик&lang=ru-RU",
		"https://www.tiktok.com/search?q=анастэйша&lang=ru-RU",
		"https://www.tiktok.com/search?q=бизнесменша&lang=ru-RU",

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
