package internal

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/testSpace/internal/db"
	"github.com/testSpace/internal/parsing"
	"github.com/testSpace/model"
	"log"
	"sync"
	"time"
)

func SendData(users []string, dbConnect *sql.DB, totalUsers *int32) []string {
	start := time.Now()
	zeroUser := *totalUsers
	var newUsers []string
	sendWg := sync.WaitGroup{}

	usersToSend := make(chan model.UserData, 6)

	go func(ch chan model.UserData, totalUsers *int32) {
		for {
			user := <-ch
			*totalUsers = *totalUsers + int32(1)
			user.Id = *totalUsers

			if user.Title != "" {
				sendWg.Add(1)
				go func(user model.UserData) {
					db.DBAddUser(user, dbConnect)
					fmt.Println("Отправил - Name:", user.Title, "\nID:", user.Id)
					sendWg.Done()
				}(user)
			}

		}

	}(usersToSend, totalUsers)

	usersChan := make(chan string, 6)

	for i := 0; i < 6; i++ {
		go func(chan string) {
			ctx, cancel := chromedp.NewContext(
				context.Background(),
				chromedp.WithLogf(log.Printf),
			)
			defer cancel()

			for {
				userNick := <-usersChan

				user := parsing.ParsingAccountData(userNick, ctx)
				if user.Title != "" {
					usersToSend <- user
					fmt.Println("Спарсил - Name:", user.Title)
				}

			}

		}(usersChan)
	}


	for num, value := range users {
		usersChan <- value
		fmt.Println(num, "аккаунт из", len(users))
	}

	fmt.Println("Пользоавтелей парсить закончил")
	sendWg.Wait()

	fmt.Println("Обработал: ", *totalUsers-zeroUser, " userов || За Время:", time.Since(start))

	return newUsers
}
