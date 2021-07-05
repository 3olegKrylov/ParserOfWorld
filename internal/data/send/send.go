package internal

import (
	"context"
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
func UsersSendHanler(ch chan model.UserData, totalUsers *int32) {
	go func() {
		user := model.UserData{}
		for {
			user = <-ch
			*totalUsers = *totalUsers + int32(1)
			user.Id = *totalUsers

			if user.Title != "" {
				db.DBAddUser(user)
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
			chromedp.ProxyServer("176.9.119.170:3128")
			defer cancel()

			user := model.UserData{}

			for {

				userNick := <-UsersChan

				user = parsing.ParsingAccountData(userNick, user, ctx)
				if user.Title != "" {
					UsersToSend <- user
				}

			}

		}()
	}
}
