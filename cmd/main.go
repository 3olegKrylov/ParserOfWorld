package main

import (
	"context"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/chromedp/chromedp"
	internal "github.com/testSpace/internal/data/send"
	"github.com/testSpace/internal/db"
	"github.com/testSpace/internal/parsing"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func main() {

	data, err := ioutil.ReadFile("name.txt")

	if err != nil {
		fmt.Println(err)
	}

	urlStr := strings.Split(string(data), "\n")

	//подключение к clickhouse
	dbConnect := db.DBconnect()
	db.DBinit(dbConnect)
	defer dbConnect.Close()

	start := time.Now()
	countOfUsers := int32(0)
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancel()

	//парсинг аккаунтов
	for i := 0; i < len(urlStr); i++ {
		//инициализация карты пользователей
		userMap := db.InitUsers(dbConnect)
		countOfUsers = int32(len(userMap))

		if urlStr[i]=="" || urlStr[i]==" "{
			continue
		}



		nameUser:=strings.TrimSpace(urlStr[i])

		text := parsing.ParseFindList("https://www.tiktok.com/search?q=" + nameUser + "&lang=ru-RU", ctx)

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
