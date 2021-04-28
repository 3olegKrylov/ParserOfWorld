package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {

	// Зачитываем содержимое файла
	data, err := ioutil.ReadFile("cmd/name.txt")
	// Если во время считывания файла произошла ошибка
	// выводим ее
	if err != nil {
		fmt.Println(err)
	}


	// Если чтение данных прошло успено
	// выводим их в консоль
	fmt.Print(string(data))

	nameArr := strings.Split(string(data),"\n")
	fmt.Println(nameArr)
	fmt.Println(len(nameArr))
	//wg:= sync.WaitGroup{}
	//wg.Add(1000)
	//
	//print:=make(chan int)
	//
	//
	//for i:= 0; i<1000;i++{
	//	go func() {
	//		print <- i
	//		wg.Done()
	//	}()
	//}
	//
	//go func(chan int) {
	//	total := 0
	//
	//	for{
	//		var newPrin int
	//		newPrin = <- print
	//		total ++
	//		fmt.Println("total: ", total, " Value: ", newPrin)
	//	}
	//}(print)
	//
	//wg.Wait()

	//парсер аккаунтов
	//start := time.Now()
	//
	//
	//	ctx, cancel := chromedp.NewContext(
	//		context.Background(),
	//		chromedp.WithLogf(log.Printf),
	//	)
	//	defer cancel()
	//
	//	parsing.ParsingAccountData("mariakurashina5", ctx )
	//
	//
	//
	//	fmt.Println(time.Since(start))
}
