package main

import (
	"context"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/chromedp/chromedp"
	internal "github.com/testSpace/internal/data/send"
	"github.com/testSpace/internal/db"
	"github.com/testSpace/internal/fullscreen"
	"github.com/testSpace/internal/parsing"
	"github.com/testSpace/model"
	"io/ioutil"
	"log"
	_ "net/http/pprof"
	"strings"
	"time"
)

func main() {
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	data, err := ioutil.ReadFile("cmd/name.txt")
	internal.UsersToSend = make(chan model.UserData, 6)

	if err != nil {
		fmt.Println(err)
	}

	//подключение к clickhouse
	db.DBconnect()
	db.DBinit()
	defer db.Connect.Close()

	urlStr := strings.Split(string(data), "\n")

	countOfUsers := int32(0)
	countOfUsers = db.InitUsers()
	startCountOfUsers := countOfUsers

	internal.UsersSendHanler(internal.UsersToSend, &countOfUsers)
	internal.SendDataHandlers(2)

	start := time.Now()
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancel()

	sendingUsers := 0
	//парсинг аккаунтов
	var text string
	var lines []string
	var ok bool

	for i := 0; i < len(urlStr); i++ {
		//инициализация карты пользователей

		if urlStr[i] == "" || urlStr[i] == " " {
			continue
		}

		nameUser := strings.TrimSpace(urlStr[i])

		//по не понятной причине причине после 12 переходов по страницам поиска tiktok, отображается отсутсвие пользователей на любой запрос
		//поэтому требуется перезагрузить браузер для продолжения поисков
		for {
			text = parsing.ParseFindList("https://www.tiktok.com/search?lang=ru-RU&q="+nameUser, ctx)

			lines = strings.Split(text, "\n\n")
			fmt.Println("число аккаунтов: ", len(lines)/3)
			if len(lines)/3 < 8 {

				var buf []byte
				if err = chromedp.Run(ctx, fullscreen.FullScreenshot(90, &buf)); err != nil {
					log.Fatal(err)
				}
				if err = ioutil.WriteFile("fullScreenshot.png", buf, 0o644); err != nil {
					log.Fatal(err)
				}

				cancel()
				ctx, cancel = chromedp.NewContext(
					context.Background(),
					chromedp.WithLogf(log.Printf),
				)

				defer cancel()
				fmt.Println(lines)

			} else {
				break
			}
		}

		for num, value := range lines {
			if strings.HasSuffix(value, "Подписчики") {
				ok = db.FindUserDB(strings.TrimSpace(lines[num-1]))
				if !ok {
					internal.UsersChan <- strings.TrimSpace(lines[num-1])
					sendingUsers++
				}
			}

		}

		fmt.Println("закончил отправлять ", urlStr[i])
		time.Sleep(time.Millisecond * 200)

	}

	fmt.Println("Отправил пользователей для дб: ", sendingUsers)
	fmt.Println("Отправилось в итоге в DB: ", countOfUsers-startCountOfUsers)
	fmt.Println("Время: ", time.Since(start), " \nКол-во юзеров: ", countOfUsers)

}
