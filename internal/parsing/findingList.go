package parsing

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func ParseFindList(urlStr string, ctx context.Context) string {
	var text string

	err := chromedp.Run(ctx,
		chromedp.Navigate(urlStr),
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

	// Рекурсивно, ждёт пока стрница прогрузится, проверяет существует ли кнопка ещё на стрнице, кликает при существовании, заканчивает при отсутствии.
	log.Println("Начинаем прогружать страницу прогружаем страницу")
	count := 0

	siteIsParse := true
	fmt.Println("прогружаем: ", urlStr)
	for {
		err = chromedp.Run(ctx,
			RunWithTimeOut(3,
				chromedp.Tasks{
					chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),

				}),
		)

		if err != nil {
			log.Println("прогрузили страницу", urlStr)
			break
		}
		count++

		if count > 70 {
			log.Println("перезагружаем страницу: ", urlStr)
			err = chromedp.Run(ctx,
				chromedp.Sleep(time.Millisecond*500),
			)

			text = ParseFindList(urlStr, ctx)
			siteIsParse = false
			break
		}
	}
	//если распарсили страницу (т.е нашли и нажимали кнопки ещё пока их не стало)
	if siteIsParse == false {
		return text
	} else {

		if err != nil && err.Error() != "context deadline exceeded" {
			log.Println("Waiting and Clicking button error:", err)
		}

		err = chromedp.Run(ctx,
			chromedp.Text(`.search`, &text, chromedp.ByQuery),
		)

		if err != nil {
			log.Println(err)
		}
		return text
	}
}
