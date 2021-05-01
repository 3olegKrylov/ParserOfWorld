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
)

var UsersToSend chan model.UserData
var UsersChan chan string

//sener users to DB включать в gorutine
func UsersSendHanler(ch chan model.UserData, totalUsers *int32, dbConnect *sql.DB) {
	go func() {
		for {
			user := <-ch
			*totalUsers = *totalUsers + int32(1)
			user.Id = *totalUsers

			if user.Title != "" {
				db.DBAddUser(user, dbConnect)
				fmt.Println("Отправил - Name:", user.Title, "\nID:", user.Id)
			}
		}

	}()
}

func SendDataHandlers(chromedbHandlCount int) {
	UsersChan = make(chan string, chromedbHandlCount)

	for i := 0; i < chromedbHandlCount; i++ {

		go func() {
			ctx, cancel := chromedp.NewContext(
				context.Background(),
				chromedp.WithLogf(log.Printf),
			)
			defer cancel()

			for {
				userNick := <-UsersChan

				user := parsing.ParsingAccountData(userNick, ctx)
				if user.Title != "" {

					UsersToSend <- user
					fmt.Println("Спарсил пользователя: ", user.Title)
				}

			}

		}()
	}
}

//func SendData(users []string, totalUsers *int32, chromedpHandleCount int) []string {
//
//	start := time.Now()
//	zeroUser := *totalUsers
//	var newUsers []string
//	sendWg := sync.WaitGroup{}
//
//	usersChan := make(chan string, 6)
//
//	for i := 0; i < 6; i++ {
//		go func(chan string) {
//			ctx, cancel := chromedp.NewContext(
//				context.Background(),
//				chromedp.WithLogf(log.Printf),
//			)
//			defer cancel()
//
//			for {
//				userNick := <-usersChan
//
//				user := parsing.ParsingAccountData(userNick, ctx)
//				if user.Title != "" {
//					UsersToSend <- user
//					fmt.Println("Спарсил - Name:", user.Title)
//				}
//
//			}
//
//		}(usersChan)
//	}
//
//	for num, value := range users {
//		usersChan <- value
//		fmt.Println(num, "аккаунт из", len(users))
//	}
//
//	fmt.Println("Пользоавтелей парсить закончил")
//	sendWg.Wait()
//
//	fmt.Println("Обработал: ", *totalUsers-zeroUser, " userов || За Время:", time.Since(start))
//
//	return newUsers
//}
