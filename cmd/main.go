package main

import (
	"context"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/chromedp/chromedp"
	internal "github.com/testSpace/internal/data/send"
	"github.com/testSpace/internal/db"
	"github.com/testSpace/internal/parsing"
	"github.com/testSpace/model"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	data, err := ioutil.ReadFile("cmd/name.txt")
	internal.UsersToSend = make(chan model.UserData, 6)

	if err != nil {
		fmt.Println(err)
	}

	//подключение к clickhouse
	dbConnect := db.DBconnect()
	db.DBinit(dbConnect)
	defer dbConnect.Close()

	urlStr := strings.Split(string(data), "\n")


	countOfUsers := int32(0)
	countOfUsers = db.InitUsers(dbConnect)
	startCountOfUsers := countOfUsers

	internal.UsersSendHanler(internal.UsersToSend, &countOfUsers, dbConnect)
	internal.SendDataHandlers(2)

	start := time.Now()

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancel()

	sendingUsers:= 0
	//парсинг аккаунтов
	for i := 0; i < len(urlStr); i++ {
		//инициализация карты пользователей

		if urlStr[i] == "" || urlStr[i] == " " {
			continue
		}

		nameUser := strings.TrimSpace(urlStr[i])

		text := parsing.ParseFindList("https://www.tiktok.com/search?q="+nameUser+"&lang=ru-RU", ctx)
		lines := strings.Split(text, "\n\n")


		for num, value := range lines {
			if strings.HasSuffix(value, "Подписчики") {
				ok := db.FindUserDB(dbConnect, strings.TrimSpace(lines[num-1]))
				if !ok {
					internal.UsersChan <- strings.TrimSpace(lines[num-1])
					sendingUsers++
				}
			}
		}



		fmt.Println("закончил отправлять ", urlStr[i])
		time.Sleep(time.Millisecond * 200)


	}

	fmt.Println("Отправил пользователей для дб: ",sendingUsers)
	fmt.Println("Отправилось в итоге в DB: ", countOfUsers - startCountOfUsers)
	fmt.Println("Время: ", time.Since(start), " \nКол-во юзеров: ", countOfUsers)

}
