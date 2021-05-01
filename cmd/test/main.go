package main

import (
	"fmt"
	"time"
)

func main() {

	chanUser:=make(chan string, 4)

	go func() {
		for{
			time.Sleep(time.Second * 1)
			fmt.Println(<-chanUser)

		}	}()

	countSend:=0
	for{
		countSend++
		chanUser<-"message"
		fmt.Println("send the ", countSend, "message")
	}


}

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
