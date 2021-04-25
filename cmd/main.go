package main

import (
	"bufio"
	"fmt"
	internal "github.com/testSpace/internal/data/send"
	"github.com/testSpace/internal/parsing"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
)

func main() {
	urlStr := []string{
		"https://www.tiktok.com/search?q=настюля&lang=ru-RU",
		"https://www.tiktok.com/search?q=настя&lang=ru-RU",
		"https://www.tiktok.com/search?q=анастасия&lang=ru-RU",
		"https://www.tiktok.com/search?q=настюня&lang=ru-RU",
		"https://www.tiktok.com/search?q=нутик&lang=ru-RU",
		"https://www.tiktok.com/search?q=анастэйша&lang=ru-RU",
		"https://www.tiktok.com/search?q=маша&lang=ru-RU",
		"https://www.tiktok.com/search?q=машуня&lang=ru-RU",
		"https://www.tiktok.com/search?q=мария&lang=ru-RU",
		"https://www.tiktok.com/search?q=мила&lang=ru-RU",
		"https://www.tiktok.com/search?q=милана&lang=ru-RU",
		"https://www.tiktok.com/search?q=русская&lang=ru-RU",
		"https://www.tiktok.com/search?q=ксюша&lang=ru-RU",
		"https://www.tiktok.com/search?q=ксения&lang=ru-RU",
		"https://www.tiktok.com/search?q=инст&lang=ru-RU",
		"https://www.tiktok.com/search?q=инста&lang=ru-RU",
		"https://www.tiktok.com/search?q=катя&lang=ru-RU",
		"https://www.tiktok.com/search?q=катюша&lang=ru-RU",
		"https://www.tiktok.com/search?q=екатерина&lang=ru-RU",
		"https://www.tiktok.com/search?q=люба&lang=ru-RU",
		"https://www.tiktok.com/search?q=вероника&lang=ru-RU",
		"https://www.tiktok.com/search?q=клава&lang=ru-RU",
		"https://www.tiktok.com/search?q=вера&lang=ru-RU",
		"https://www.tiktok.com/search?q=ника&lang=ru-RU",
		"https://www.tiktok.com/search?q=аня&lang=ru-RU",
		"https://www.tiktok.com/search?q=анна&lang=ru-RU",
		"https://www.tiktok.com/search?q=подняла&lang=ru-RU",
		"https://www.tiktok.com/search?q=насосала&lang=ru-RU",
		"https://www.tiktok.com/search?q=таня&lang=ru-RU",
		"https://www.tiktok.com/search?q=танюша&lang=ru-RU",
	}

	//подключение к clickhouse
	//dbConnect := db.DBconnect()
	//db.DBinit(dbConnect)
	//defer dbConnect.Close()
	var userMap map[string]int32
	if _, err := os.Stat("users.tsv"); os.IsNotExist(err) {
		// path/to/whatever does not exist
	}else{
		userMap = UserMapInit("users.tsv")
	}

	file, err := os.OpenFile("users.tsv", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	str := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "Id", "LinkAccount", "Title", "SubTitle", "Comment", "Mail", "Telegram", "Instagram", "Links", "LanguageAccount", "Phone", "Following", "Followers", "Likes", "LastPostShowTotal", "AverageShows", "MedianShows", "TotalPosts", "LastActionTime", "ParsingTime")
	file.Write([]byte(str))

	start := time.Now()
	countOfUsers := int32(0)

	fmt.Println(userMap)
	fmt.Println(len(userMap))
	//парсинг аккаунтов
	for i := 0; i < len(urlStr); i++ {
		//инициализация карты пользователей

		countOfUsers = int32(len(userMap))
		text := parsing.ParseFindList(urlStr[i])

		lines := strings.Split(text, "\n\n")

		var newUsers []string

		for num, value := range lines {
			if strings.HasSuffix(value, "Подписчики") {
				_, ok := userMap[strings.TrimSpace(lines[num-1])]
				if !ok {
					newUsers = append(newUsers, lines[num-1])
				}
			}
		}
		if len(newUsers) > 0 {
			//internal.SendData(newUsers, dbConnect, &countOfUsers)
			internal.SendDataTsv(newUsers, file, &countOfUsers)
			countOfUsers = countOfUsers + int32(len(newUsers))
			fmt.Println("Время работы: ", time.Since(start))
		}
	}

	fmt.Println("Время: ", time.Since(start), " \nКол-во юзеров: ", countOfUsers)

}

func UserMapInit(fileName string) map[string]int32{

	newMapName := make(map[string]int32)

	tsv, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	scn := bufio.NewScanner(tsv)

	var lines []string

	for scn.Scan() {
		line := scn.Text()
		lines = append(lines, line)
	}

	if err := scn.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	for num, line := range lines {
		fmt.Println(line)
		record := strings.Split(line, "\t")
		if num>=2 {
			if len(record) > 3 {
				_, ok := newMapName[record[2]]
				if ok {
				} else {
					newMapName[record[2]] = int32(0)
				}
			}
		}

	}

	return newMapName
}

