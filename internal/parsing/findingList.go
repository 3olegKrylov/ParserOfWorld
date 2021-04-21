package parsing

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/testSpace/internal/fullscreen"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

func ParseFindList(urlStr string) string {
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

	err = chromedp.Run(ctx, RunWithTimeOut(
		1,
		chromedp.Tasks{chromedp.Text(`.error-page`, &checkClearFinding, chromedp.NodeVisible, chromedp.ByQuery)},
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
			RunWithTimeOut(3,
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

			text = ParseFindList(urlStr)
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
