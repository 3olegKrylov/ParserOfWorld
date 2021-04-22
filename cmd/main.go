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
		"https://www.tiktok.com/search?q=олег&lang=ru-RU",
		"https://www.tiktok.com/search?q=лена&lang=ru-RU",
		"https://www.tiktok.com/search?q=егор&lang=ru-RU",
		"https://www.tiktok.com/search?q=влад&lang=ru-RU",
		"https://www.tiktok.com/search?q=вова&lang=ru-RU",
		"https://www.tiktok.com/search?q=настя&lang=ru-RU",
		"https://www.tiktok.com/search?q=анжела&lang=ru-RU",
		"https://www.tiktok.com/search?q=женя&lang=ru-RU",
		"https://www.tiktok.com/search?q=витя&lang=ru-RU",
		"https://www.tiktok.com/search?q=ксюша&lang=ru-RU",
		"https://www.tiktok.com/search?q=полина&lang=ru-RU",
		"https://www.tiktok.com/search?q=лера&lang=ru-RU",
		"https://www.tiktok.com/search?q=аня&lang=ru-RU",
		"https://www.tiktok.com/search?q=лиза&lang=ru-RU",
		"https://www.tiktok.com/search?q=вика&lang=ru-RU",
		"https://www.tiktok.com/search?q=дина&lang=ru-RU",
		"https://www.tiktok.com/search?q=юля&lang=ru-RU",
		"https://www.tiktok.com/search?q=маша&lang=ru-RU",
		"https://www.tiktok.com/search?q=степан&lang=ru-RU",
		"https://www.tiktok.com/search?q=евгений&lang=ru-RU",
		"https://www.tiktok.com/search?q=антон&lang=ru-RU",
		"https://www.tiktok.com/search?q=марк&lang=ru-RU",
		"https://www.tiktok.com/search?q=артём&lang=ru-RU",
		"https://www.tiktok.com/search?q=бизнес&lang=ru-RU",
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
		fmt.Println("userMap: ", userMap)
		text := parsing.ParseFindList(urlStr[i])

		lines := strings.Split(text, "\n\n")

		var newUsers []string

		for num, value := range lines {
			if strings.HasSuffix(value, "Подписчики") {
				_, ok := userMap[strings.TrimSpace(lines[num-1])]
				if !ok {
					fmt.Println("newUser: ", lines[num-1])
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
