package parsing

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	chromedp "github.com/chromedp/chromedp"
	"github.com/mcnijman/go-emailaddress"
	"github.com/testSpace/model"
	"log"
	"mvdan.cc/xurls/v2"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ParsingAccountData(nick string, ctx context.Context) model.UserData {
	url := "https://www.tiktok.com/@" + nick

	fmt.Println("1 ", url)

	user := model.UserData{
		Id:                0,
		LinkAccount:       url,
		Title:             "",
		SubTitle:          "",
		Comment:           "",
		Mail:              "",
		Telegram:          "",
		Instagram:         "",
		Links:             "",
		LanguageAccount:   "",
		Phone:             "",
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

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("2 ", url)
	//смотрим существует ли страница - если нет title значит страницы нет
	userIsExist := ""
	err = chromedp.Run(ctx, RunWithTimeOut(
		1,
		chromedp.Tasks{chromedp.Text(`.share-title`, &userIsExist, chromedp.NodeVisible, chromedp.ByQuery)},
	))


	if userIsExist == ""{
		fmt.Println("3* ", url, "не найден")
		err = chromedp.Run(ctx, RunWithTimeOut(
			1,
			chromedp.Tasks{chromedp.Sleep(time.Millisecond*200),
				chromedp.Text(`.share-title`, &userIsExist, chromedp.NodeVisible, chromedp.ByQuery)},
		))
		if userIsExist == "" {
			fmt.Println("3* ", url, "не найден - не существует")
			return user
		}
	}

	titleUser := ""

	err = chromedp.Run(ctx,
		chromedp.Text(`.share-title`, &titleUser, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.share-sub-title`, &user.SubTitle, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.share-desc`, &user.Comment, chromedp.NodeVisible, chromedp.ByQuery),
		chromedp.Text(`.count-infos`, &numericData, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}
	user.Title = strings.TrimSpace(titleUser)

	user.Following, user.Followers, user.Likes = numericDataParser(numericData)

	//проверяем существет ли ссылка у пользователя
	linkOnTitile := ""
	err = chromedp.Run(ctx, RunWithTimeOut(
		1,
		chromedp.Tasks{chromedp.Text(`.share-links`, &linkOnTitile, chromedp.NodeVisible, chromedp.ByQuery)},
	))

	if linkOnTitile != "" {
		user.Links = linkOnTitile
	}

	moreLinks := ""

	//если комментарий не пустой то распарсим всё что в нём есть
	if user.Comment != "No bio yet." {
		moreLinks, user.Phone, user.Instagram, user.Telegram, user.Mail = commentParser(user.Comment)
	}

	user.Links = user.Links + moreLinks

	//проверяем есть ли посты у пользователя если там error-page то постов нет
	checkClearAccount := ""
	err = chromedp.Run(ctx, RunWithTimeOut(
		1,
		chromedp.Tasks{chromedp.Text(`.error-page`, &checkClearAccount, chromedp.NodeVisible, chromedp.ByQuery)},
	))
	//У пользователя нет контента
	if checkClearAccount != "" {
		return user
	}
	//Todo: записываем все лайки на карточках (чтобы записать все надо дождаться пока прогрузиться вся страница)
	err = chromedp.Run(ctx,
		chromedp.Text(`.tt-feed`, &likesCard, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		log.Fatal(err)
	}
	//TODO:логически если это будет ближе к началу
	//при всплывающем модалке начинаем всё заново (пока её не станет)
	for {
		what := ""
		err = chromedp.Run(ctx,
			RunWithTimeOut(
				2,
				chromedp.Tasks{chromedp.Text(`.iframe-container`, &what, chromedp.NodeVisible, chromedp.ByQuery),
				},
			),
		)
		if err != nil {
			break
		} else {
			return ParsingAccountData(nick, ctx)
		}
	}
	//переходим по ссылке первого видео кликая на неё
	for {
		err = chromedp.Run(ctx,
			RunWithTimeOut(1,
				chromedp.Tasks{chromedp.Click(`._ratio_wrapper`, chromedp.NodeVisible, chromedp.ByQuery),
				},
			))

		if err == nil {
			break
		}
	}

	for {
		err = chromedp.Run(ctx,
			RunWithTimeOut(1,
				chromedp.Tasks{chromedp.Text(`.author-nickname`, &ActionTime, chromedp.NodeVisible, chromedp.ByQuery),
				},
			))

		if err == nil {
			break
		}

		checkClearVideo := ""
		err = chromedp.Run(ctx,
			chromedp.Reload(),
			RunWithTimeOut(
				1,
				chromedp.Tasks{chromedp.Text(`.error-page`, &checkClearVideo, chromedp.NodeVisible, chromedp.ByQuery)},
			),
		)

		if checkClearVideo != "" {
			break
		}

	}

	//парсим данные карточек
	user.LastPostShowTotal, _, _, _ = parserCardShows(likesCard)
	if ActionTime != "" {
		user.LastActionTime = time.Time(lastActionTimeParser(ActionTime))
	}
	fmt.Println("12 ", url)
	return user
}

//Todo: не правильно работает, т.к не успевает прогрузить все карточки, соответсвенно надо ждать и потом всё считать
//возвращает количество просмотров последнего поста, среднее кол-во просмотров, кол-во постов
func parserCardShows(data string) (int32, int32, int32, int32) {
	cardShowArr := strings.Split(data, "\n")
	var cardsLikes []int32
	total := int32(0)

	for num, _ := range cardShowArr {
		likeOfCard := numParser(strings.TrimSpace(cardShowArr[num]))
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
func numericDataParser(data string) (countFollowing int32, countFollowers int32, Likes int32) {
	dataArr := strings.Split(data, "\n")
	return numParser(dataArr[0]), numParser(dataArr[2]), numParser(dataArr[4])
}

//перевдит строки вида 21.2M в 21200000 и 21.1K в 21100 вернёт -1 если не смогу распарсить
func numParser(num string) int32 {
	num = strings.TrimSpace(num)
	switch num[(len(num) - 1):] {
	case "M":
		num = strings.Replace(num, "M", "", -1)
		val, err := strconv.ParseFloat(num, 32)
		if err != nil {
			//log.Println("Не смогу распарсить данные в float32: ", num, "\n", err)
			return -1
		}
		return int32(val * 1000000)
	case "K":
		num = strings.Replace(num, "K", "", -1)
		val, err := strconv.ParseFloat(num, 32)
		if err != nil {
			//log.Println("Не смогу распарсить данные в float32: ", num, "\n", err)
			return -1
		}
		return int32(val * 1000)

	default:
		val, err := strconv.ParseFloat(num, 32)
		if err != nil {
			//log.Println("Не смогу распарсить данные в float32: ", num, "\n", err)
			return -1
		}
		return int32(val)
	}

}

//TODO: парсить дату для "3 дня назад"
//прасит текст со временем на сранице видео
func lastActionTimeParser(data string) clickhouse.Date {

	dataArr := strings.Split(data, "·")
	resStr := strings.TrimSpace(dataArr[1])

	fmt.Println(resStr)
	dataType := strings.Split(resStr, "-")

	if strings.Contains(dataType[0], " ") {
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

//парисит из описания/комментария при наличии ссылки номера телефонов, инстаграм, тг, почты
func commentParser(comment string) (links string, phoneNum string, instagram string, telegram string, mail string) {
	rxRelaxed := xurls.Relaxed()
	relUrl := rxRelaxed.FindString(comment)

	rxStrict := xurls.Strict()
	linksArr := rxStrict.FindAllString(comment, -1)

	for _, l := range linksArr {
		links = links + " " + l
	}
	if relUrl != "" {
		links = links + " " + relUrl
	}

	comment = strings.ToLower(comment)
	comment = strings.Replace(comment, "\n", " ", -1)
	comment = strings.Replace(comment, ",", " ", -1)
	words := strings.Split(comment, " ")

	var res []string
	for _, val := range words {
		newWords := strings.Split(val, "\n")
		res = append(res, newWords...)
	}

	words = res
	//парсер для номеров
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	//fmt.Println("Words", words)
	for num, value := range words {

		//instagram
		if value == "insta:" || value == "insta" || value == "inst:" || value == "inst" || value == "instagram:" || value == "instagram" || value == "инст:" || value == "инст" || value == "инстаграм:" || value == "инстаграм" || value == "инста" || value == "инста:" || value == "инстаграмм:" || value == "инстаграмм" || value == "ig" || value == "ig:" {
			if len(words) >= num+2 {
				if words[num+1] == "-" || words[num+1] == ":" {
					if len(words) >= num+3 {
						instagram += words[num+2] + " "
					}
				} else {
					instagram += words[num+1] + " "
				}
			}
		} else if strings.HasSuffix(value, "insta:") || strings.HasSuffix(value, "inst:") || strings.HasSuffix(value, "instagram:") || strings.HasSuffix(value, "инст:") || strings.HasSuffix(value, "инстаграм:") || strings.HasSuffix(value, "инста:") || strings.HasSuffix(value, "инстаграмм:") || strings.HasSuffix(value, "ig:") {
			instagram += words[num] + " "
		}

		//phone
		arrNumber := re.MatchString(value)
		if arrNumber == true && value != "-" {
			phoneNum = phoneNum + " " + value
		}
		if value == "phone" || value == "phone:" {
			if len(words) >= num+2 {
				if words[num+1] == "-" || words[num+1] == ":" {
					if len(words) >= num+3 {
						phoneNum += words[num+2] + " "
					}
				} else {
					phoneNum += words[num+1] + " "
				}
			}
		}

		//telegram
		if value == "telegram" || value == "tg" || value == "telegram:" || value == "tg:" || value == "тг" || value == "телеграм:" || value == "телеграм" || value == "тг:" || value == "телеграмм:" || value == "телеграмм" {
			if len(words) >= num+2 {
				if words[num+1] == "-" || words[num+1] == ":" {
					if len(words) >= num+3 {
						telegram += words[num+2] + " "
					}
				} else {
					telegram += words[num+1] + " "
				}
			}
		} else if strings.HasSuffix(value, "telegram:") || strings.HasSuffix(value, "tg:") || strings.HasSuffix(value, "телеграм:") || strings.HasSuffix(value, "тг:") || strings.HasSuffix(value, "телеграмм:") {
			instagram += words[num] + " "
		}

	}

	text := []byte(comment)
	validateHost := false

	emails := emailaddress.Find(text, validateHost)

	for _, e := range emails {
		mail = mail + e.String() + ""
	}

	return links, phoneNum, instagram, telegram, mail
}
