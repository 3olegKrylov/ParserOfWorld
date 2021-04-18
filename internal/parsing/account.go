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
	"sort"
	"strconv"
	"strings"
	"time"
)

func ParsingAccountData(nick string, ctx context.Context, id int32) model.UserData {
	url := "https://www.tiktok.com/@" + nick

	user := model.UserData{
		Id:                id,
		LinkAccount:       url,
		Title:             "",
		SubTitle:          "",
		Comment:           "",
		Mail:              "",
		Telegram:          "",
		Instagram:         "",
		Linkes:            "",
		LanguageAccount:   "",
		Following:         0,
		Followers:         0,
		Likes:             0,
		LastPostShowTotal: 0,
		AverageShows:      0,
		MedianShows:       0,
		TotalPosts:        0,
		LastActionTime:    time.Time{},
		ParsingTime:       time.Now(),
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


	userIsExist:=""
	err = chromedp.Run(ctx, RunWithTimeOut(
		1,
		chromedp.Tasks{chromedp.Text(`.title`, &userIsExist, chromedp.NodeVisible, chromedp.ByQuery),
		},
	))

	if userIsExist == "Couldn't find this account" {
		return model.UserData{}
	}

	err = chromedp.Run(ctx,
		chromedp.Text(`.share-title`, &user.Title, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.share-sub-title`, &user.SubTitle, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.share-desc`, &user.Comment, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.count-infos`, &numericData, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

	checkClearAccount := ""

	err = chromedp.Run(ctx, RunWithTimeOut(
		1,
		chromedp.Tasks{chromedp.Text(`.error-page`, &checkClearAccount, chromedp.NodeVisible, chromedp.ByQuery),
		},
	))
	fmt.Println(checkClearAccount)
	user.Followers, user.Following, user.Likes = parserNumericData(numericData)

	//Todo: распарсить комментарий
	//У пользователя нет контента
	if checkClearAccount != "" {
		return user
	}

	err = chromedp.Run(ctx,
		chromedp.Text(`.tt-feed`, &likesCard, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}

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

	//при всплывающем модалке начинаем всё заново (пока её не станет)
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
			return ParsingAccountData(nick, ctx, id)

		}
	}

	err = chromedp.Run(ctx,
		chromedp.Click(`._ratio_wrapper`, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Sleep(time.Millisecond*100),
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

	fmt.Println("likesCard: ", likesCard)
	user.LastPostShowTotal, _, _, _ = parserCardShows(likesCard)
	fmt.Println("parserCardShows end")
	user.LastActionTime = time.Time(parseLastActionTime(ActionTime))
	fmt.Println(user)

	return user
}

//Todo: не правильно работает, т.к не успевает прогрузить все карточки, соответсвенно надо ждать и потом всё считать
//возвращает количество просмотров последнего поста, среднее кол-во просмотров, кол-во постов
func parserCardShows(data string) (int32, int32, int32, int32) {
	cardShowArr := strings.Split(data, "\n")
	var cardsLikes []int32
	total := int32(0)

	for num, _ := range cardShowArr {
		likeOfCard := parserNum(strings.TrimSpace(cardShowArr[num]))
		if likeOfCard != -1 {
			cardsLikes = append(cardsLikes, likeOfCard)
			total = total + likeOfCard
		}
	}
	var likeFirstCard int32

	if len(cardsLikes) > 0 {
		likeFirstCard = cardsLikes[0]
	} else {
		likeFirstCard = 0
	}

	//
	//median := int32(0)
	//
	//if len(cardsLikes) > 1 {
	//	median = nlogn_median(cardsLikes)
	//} else {
	//	median = cardsLikes[0]
	//}

	//return likeFirstCard, total / int32(len(cardShowArr)), int32(len(cardsLikes)), median
	return likeFirstCard, 0, 0, 0
}

func nlogn_median(l []int32) int32 {
	slice := l
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] > slice[j]
	})

	if len(l)%2 == 1 {
		return slice[len(slice)/2]
	} else {

		return slice[len(slice)/2-1]
	}
}

//парсит строку строку аккаунта, где содержится информация об аккаунте
func parserNumericData(data string) (countFollowing int32, countFollowers int32, Likes int32) {
	dataArr := strings.Split(data, "\n")
	return parserNum(dataArr[0]), parserNum(dataArr[2]), parserNum(dataArr[4])
}

//перевдит строки вида 21.2M в 21200000 и 21.1K в 21100 вернёт -1 если не смогу распарсить
func parserNum(num string) int32 {
	num = strings.TrimSpace(num)
	switch num[(len(num) - 1):] {
	case "M":
		num = strings.Replace(num, "M", "", -1)
		val, err := strconv.ParseFloat(num, 32)
		if err != nil {
			log.Println("Не смогу распарсить данные в float32: ", num, "\n", err)
			return -1
		}
		return int32(val * 1000000)
	case "K":
		num = strings.Replace(num, "K", "", -1)
		val, err := strconv.ParseFloat(num, 32)
		if err != nil {
			log.Println("Не смогу распарсить данные в float32: ", num, "\n", err)
			return -1
		}
		return int32(val * 1000)

	default:
		val, err := strconv.ParseFloat(num, 32)
		if err != nil {
			log.Println("Не смогу распарсить данные в float32: ", num, "\n", err)
			return -1
		}
		return int32(val)
	}

}

func parseLastActionTime(data string) clickhouse.Date {
	dataArr := strings.Split(data, "·")
	resStr := strings.TrimSpace(dataArr[1])

	fmt.Println(resStr)
	dataType := strings.Split(resStr, "-")

	if strings.Contains(dataType[0]," "){
		return clickhouse.Date(time.Now())
	}

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
		fmt.Println(dataType)
		fmt.Println(len(dataType))
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
