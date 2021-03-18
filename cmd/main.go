package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"math"
	"strings"
	"time"
)

var buttonExist []*cdp.Node

func main() {
	urlStr := "https://www.tiktok.com/search?q=ааа&lang=ru-RU"
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

		//находит кнопку ещё
		//chromedp.Nodes(`.jsx-3392055781 .more`, &nodes, chromedp.AtLeast(0)),

	)
	if err != nil {
		log.Println( err)
	}

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
	if err != nil {
		log.Println( err)
	}

	err = chromedp.Run(ctx,
		chromedp.Text(`.search`, &text, chromedp.ByQuery),
		fullScreenshot(urlStr, 1, &buf),
	)

	if err != nil {
		log.Println( err)
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


	for name, comment := range DataOfUsers{
		fmt.Println("Name:", name,  "\nValue:", comment,"\n")
	}

	fmt.Println("Количество пользователей: ",  len(DataOfUsers))
	fmt.Println("Время работы: ", time.Since(start))

}

//делает скиншот экрана полноразмерный заптсывавает под названием elementScreenshot.png
func fullScreenshot(urlstr string, quality int64, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			// get layout metrics
			_, _, contentSize, err := page.GetLayoutMetrics().Do(ctx)
			if err != nil {
				return err
			}

			width, height := int64(math.Ceil(contentSize.Width)), int64(math.Ceil(contentSize.Height))

			// force viewport emulation
			err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
				WithScreenOrientation(&emulation.ScreenOrientation{
					Type:  emulation.OrientationTypePortraitPrimary,
					Angle: 0,
				}).
				Do(ctx)
			if err != nil {
				return err
			}

			// capture screenshot
			*res, err = page.CaptureScreenshot().
				WithQuality(quality).
				WithClip(&page.Viewport{
					X:      contentSize.X,
					Y:      contentSize.Y,
					Width:  contentSize.Width,
					Height: contentSize.Height,
					Scale:  1,
				}).Do(ctx)
			if err != nil {
				return err
			}
			return nil
		}),
	}
}

//10 итераций - ждёт пока прогрузится объект
func clickShowMore() chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
	}
}

//TODO:
// Рекурсивно, ждёт пока стрница прогрузится, проверяет существует ли кнопка ещё на стрнице, кликает при существовании, заканчивает при отсутствии.
func recursiveSerchingAndClickButton(ctx context.Context) chromedp.Tasks {
	return chromedp.Tasks{
		RunWithTimeOut(&ctx, 3, chromedp.Tasks{
			chromedp.Click(`.more`, chromedp.NodeVisible, chromedp.ByQuery),
		}),
		recursiveSerchingAndClickButton(ctx),
	}
}

//TODO: принимает @data текст с аккаунтами, создаёт map c данными юзеров
func PasrsingAccounts(data *string) string {

	fmt.Println("реализовать PasrsingAccounts")
	return ""
}

//ждём появления кнопки
func RunWithTimeOut(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		return tasks.Do(timeoutContext)
	}
}
