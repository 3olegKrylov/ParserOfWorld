package internal

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/testSpace/internal/db"
	"github.com/testSpace/internal/parsing"
	"log"
)

func SendData(users []string, dbConnect *sql.DB, totalUsers *int32) []string {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var newUsers []string

	for _, value := range users {

		user := parsing.ParsingAccountData(value, ctx, *totalUsers+1)
		if user.Title != "" {
			*totalUsers = *totalUsers + int32(1)
			go func() { db.DBAddUser(user, dbConnect) }()
			newUsers = append(newUsers, user.Title)
			fmt.Println("Name:", user.Title, "\nID:", user.Id)
		}
	}

	return newUsers
}
