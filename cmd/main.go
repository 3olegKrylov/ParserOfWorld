package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	internal "github.com/testSpace/internal/data/send"
	"github.com/testSpace/internal/db"
	"github.com/testSpace/internal/fullscreen"
	"github.com/testSpace/internal/parsing"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

func main() {
	urlStr := []string{
		"https://www.tiktok.com/search?q=ksenii_or&lang=ru-RU",
		"https://www.tiktok.com/search?q=ааааа&lang=ru-RU",
		"https://www.tiktok.com/search?q=ббббб&lang=ru-RU",
		"https://www.tiktok.com/search?q=ссссс&lang=ru-RU",
		"https://www.tiktok.com/search?q=ддддд&lang=ru-RU",
		"https://www.tiktok.com/search?q=ййййй&lang=ru-RU",
		"https://www.tiktok.com/search?q=ццццц&lang=ru-RU",
		"https://www.tiktok.com/search?q=ууууу&lang=ru-RU",
		"https://www.tiktok.com/search?q=ккккк&lang=ru-RU",
		"https://www.tiktok.com/search?q=еееее&lang=ru-RU",
		"https://www.tiktok.com/search?q=ннннн&lang=ru-RU",
		"https://www.tiktok.com/search?q=ггггг&lang=ru-RU",
		"https://www.tiktok.com/search?q=жжжжж&lang=ru-RU",
		"https://www.tiktok.com/search?q=эээээ&lang=ru-RU",
		"https://www.tiktok.com/search?q=иииии&lang=ru-RU",
		"https://www.tiktok.com/search?q=яяяяя&lang=ru-RU",
		"https://www.tiktok.com/search?q=ттттт&lang=ru-RU",
		"https://www.tiktok.com/search?q=ффффф&lang=ru-RU",
		"https://www.tiktok.com/search?q=ююююю&lang=ru-RU",
		"https://www.tiktok.com/search?q=ёёёёё&lang=ru-RU",
		"https://www.tiktok.com/search?q=ззззз&lang=ru-RU",
		"https://www.tiktok.com/search?q=ыыыыы&lang=ru-RU",
		"https://www.tiktok.com/search?q=ччччч&lang=ru-RU",
	}

	//подключение к clickhouse
	dbConnect := db.DBconnect()
	db.DBinit(dbConnect)
	defer dbConnect.Close()

	//инициализация карты пользователей
	userMap := db.InitUsers(dbConnect)

	start := time.Now()
	countOfUsers := int32(len(userMap))
	//парсинг аккаунтов
	for i := 0; i < len(urlStr); i++ {
		text := Parse(urlStr[i])

		lines := strings.Split(text, "\n\n")

		var newUsers []string

		for num, value := range lines {
			if strings.HasSuffix(value, "Подписчики") {
				_, ok := userMap[lines[num-1]]
				if !ok {
					newUsers = append(newUsers, lines[num-1])
				}
			}
		}
		if len(newUsers) > 0 {
			internal.SendData(newUsers, dbConnect, countOfUsers+1)
			countOfUsers = countOfUsers + int32(len(newUsers))
			fmt.Println("Время работы: ", time.Since(start))
		}
	}

	fmt.Println("Время: ", time.Since(start), " \nКол-во юзеров: ", countOfUsers)

}

func Parse(urlStr string) string {
	var buf []byte
	var text string

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(urlStr),
		fullscreen.FullScreenshot(1, &buf),
	)
	if err != nil {
		log.Fatal("Error Navigate Parsing Accounts to ", urlStr, "\nerror: ", err)
	} else {
		log.Println("Navigate Parsing Accounts to ", urlStr)
	}

	checkClearFinding := ""

	err = chromedp.Run(ctx, parsing.RunWithTimeOut(
		1,
		chromedp.Tasks{chromedp.Text(`.error-page`, &checkClearFinding, chromedp.NodeVisible, chromedp.ByQuery),
		},
	))
	if checkClearFinding != "" {
		return ""
	}

	if err := ioutil.WriteFile("elementScreenshotStart.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	// Рекурсивно, ждёт пока стрница прогрузится, проверяет существует ли кнопка ещё на стрнице, кликает при существовании, заканчивает при отсутствии.
	log.Println("Начинаем прогружать страницу прогружаем страницу")
	count := 0

	siteIsParse := true
	for {
		err = chromedp.Run(ctx,
			parsing.RunWithTimeOut(3,
				chromedp.Tasks{
					chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
					fullscreen.FullScreenshot(1, &buf),
				}),
		)

		if err != nil {
			log.Println("ошибка выходим из for: ", err)
			break
		}

		log.Println("Нашёл кнопку ещё и нажал:")
		count++

		if count > 52 {
			if err := ioutil.WriteFile("BANelementScreenshot"+strconv.Itoa(count)+".png", buf, 0o644); err != nil {
				log.Fatal(err)
			}

			log.Println("преывшено кол-во итераций перезагружаем страницу: ", urlStr)
			err = chromedp.Run(ctx,
				chromedp.Sleep(time.Millisecond*500),
			)

			text = Parse(urlStr)
			siteIsParse = false
			break
		}
	}

	if siteIsParse == false {
		return text
	} else {
		fmt.Println("Прогурузили кнопки")

		if err != nil && err.Error() != "context deadline exceeded" {
			log.Println("Waiting and Clicking button error:", err)
		}

		err = chromedp.Run(ctx,
			chromedp.Text(`.search`, &text, chromedp.ByQuery),
		)

		if err != nil {
			log.Println(err)
		}

		// записываем данные в фотографию
		if err := ioutil.WriteFile("elementScreenshot.png", buf, 0o644); err != nil {
			log.Fatal(err)
		}
		return text
	}
}
