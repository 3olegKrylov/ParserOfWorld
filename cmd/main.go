package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	internal "github.com/testSpace/internal/data/send"
	"github.com/testSpace/internal/db"
	"github.com/testSpace/internal/fullscreen"
	"io/ioutil"
	"log"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

func main() {
	urlStr := "https://www.tiktok.com/search?q=ksenii_or&lang=ru-RU"

	dbConnect := db.DBconnect()
	db.DBinit(dbConnect)

	start := time.Now()

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	//var url string
	var buf []byte

	var text string
	//var val string

	err := chromedp.Run(ctx,
		chromedp.Navigate(urlStr),
	)
	if err != nil {
		log.Println("Navigate to ", urlStr, " error: ", err)
	}

	// Рекурсивно, ждёт пока стрница прогрузится, проверяет существует ли кнопка ещё на стрнице, кликает при существовании, заканчивает при отсутствии.
	flagButton := true
	for flagButton {
		err = chromedp.Run(ctx,
			RunWithTimeOut(&ctx, 3, chromedp.Tasks{
				chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
			}),
		)
		if err != nil {
			flagButton = false
		}
	}

	if err != nil && err.Error() != "context deadline exceeded" {
		log.Println("Waiting and Clicking button error:", err)
	}

	err = chromedp.Run(ctx,
		chromedp.Text(`.search`, &text, chromedp.ByQuery),
		fullscreen.FullScreenshot(urlStr, 1, &buf),
	)

	if err != nil {
		log.Println(err)
	}

	// записываем данные в фотографию
	if err := ioutil.WriteFile("elementScreenshot.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	lines := strings.Split(text, "\n\n")
	DataOfUsers := make(map[string]string)

	for num, value := range lines {
		if strings.HasSuffix(value, "Подписчики") {
			DataOfUsers[lines[num-1]] = value
		}
	}

	internal.SendData(text, dbConnect)

	fmt.Println("Количество пользователей: ", len(DataOfUsers))
	fmt.Println("Время работы: ", time.Since(start))

}

//ждём появления кнопки
func RunWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}

//TODO: принимает @data текст с аккаунтами, создаёт map c данными юзеров
func PasrsingAccounts(data *string) string {

	fmt.Println("реализовать PasrsingAccounts")
	return ""
}
