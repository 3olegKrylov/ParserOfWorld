package parsing

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	chromedp "github.com/chromedp/chromedp"
	"github.com/testSpace/internal/fullscreen"
	"github.com/testSpace/model"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

func ParsingAccountData(nick string, ctx context.Context, id int32) model.UserData {
	url := "https://www.tiktok.com/@" + nick

	user := model.UserData{
		Id:                   id,
		LinkAccount:          url,
		Title:                "",
		SubTitle:             "",
		Comment:              "",
		Mail:                 "",
		Telegram:             "",
		Instagram:            "",
		Linkes:               "",
		Following:            0,
		Followers:            0,
		Likes:                0,
		ShowTotal:            0,
		AverageNumberOfShows: 0,
		LastActionTime:       time.Time{},
		LanguageAccount:      "",
		ParsingTime:          time.Now(),
	}

	numericData := ""
	likesCard := ""
	ActionTime := ""
	var buf []byte

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Navigate to ", url)

	err = chromedp.Run(ctx,
		chromedp.Text(`.share-title`, &user.Title, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.share-sub-title`, &user.SubTitle, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.share-desc`, &user.Comment, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.count-infos`, &numericData, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.tt-feed`, &likesCard, chromedp.NodeVisible, chromedp.ByQuery),
	)

	log.Println("Получили данные со страницы")
	err = chromedp.Run(ctx,
		fullscreen.FullScreenshot(1, &buf),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("DataBeforClick.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	for {
		what := ""
		err = chromedp.Run(ctx,
			RunWithTimeOut(
				3,
				chromedp.Tasks{chromedp.Text(`.iframe-container`, &what, chromedp.NodeVisible, chromedp.ByQuery),
					fullscreen.FullScreenshot(1, &buf),
				},
			),
		)
		if err != nil {
			log.Println("Выходим из цикла текст: ", what, "\n", err)
			break
		} else {
			log.Println("Перезагружаем: ", what)
			return ParsingAccountData(nick,ctx,id)

		}
	}

	fmt.Println("err:", err)

	err = chromedp.Run(ctx,
		chromedp.Click(`._ratio_wrapper`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Sleep(time.Millisecond*200),
		chromedp.Reload(),
		fullscreen.FullScreenshot(1, &buf),
		chromedp.Text(`.author-nickname`, &ActionTime, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile("AccountPhoto.png", buf, 0o644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("numericData: ", numericData)
	user.Followers, user.Following, user.Likes = parserNumericData(numericData)

	user.LastActionTime = time.Time(parseLastActionTime(ActionTime))

	fmt.Println(user)
	fmt.Println("LastAction: ", time.Time(user.LastActionTime))

	return user
}

//парсит строку строку аккаунта, где содержится информация об аккаунте
func parserNumericData(data string) (countFollowing int32, countFollowers int32, Likes int32) {
	dataArr := strings.Split(data, "\n")
	return parserNum(dataArr[0]), parserNum(dataArr[2]), parserNum(dataArr[4])
}

//перевдит строки вида 21.2M в 21200000 и 21.1K в 21100
func parserNum(num string) int32 {
	num = strings.TrimSpace(num)
	switch num[(len(num) - 1):] {
	case "M":
		num = strings.Replace(num, "M", "", -1)
		val, err := strconv.ParseFloat(num, 32)
		if err != nil {
			log.Fatal("Не смогу распарсить данные в float32: ", num, "\n", err)
		}
		return int32(val * 1000000)
	case "K":
		num = strings.Replace(num, "K", "", -1)
		val, err := strconv.ParseFloat(num, 32)
		if err != nil {
			log.Fatal("Не смогу распарсить данные в float32: ", num, "\n", err)
		}
		return int32(val * 1000)

	default:
		val, err := strconv.ParseFloat(num, 32)
		if err != nil {
			log.Fatal("Не смогу распарсить данные в float32: ", num, "\n", err)
		}
		return int32(val)
	}

}

func parseLastActionTime(data string) clickhouse.Date {
	dataArr := strings.Split(data, "·")
	resStr := strings.TrimSpace(dataArr[1])

	fmt.Println(resStr)
	dataType := strings.Split(resStr, "-")

	if len(dataType) == 2 {
		year, err := strconv.Atoi(time.Now().Format("2006"))
		countMonth, err := strconv.Atoi(dataType[0])
		month := time.Month(countMonth)
		day, err := strconv.Atoi(dataType[1])
		if err != nil {
			log.Println(err)
		}
		return clickhouse.Date(time.Date(year, month, day, 0, 0, 0, 0, time.Local))
	} else {
		countMonth, err := strconv.Atoi(dataType[1])
		month := time.Month(countMonth)
		day, err := strconv.Atoi(dataType[2])
		if err != nil {
			log.Println(err)
		}
		year, err := strconv.Atoi(dataType[0])
		if err != nil {
			log.Println(err)
		}
		return clickhouse.Date(time.Date(year, month, day, 0, 0, 0, 0, time.Local))
	}
}
