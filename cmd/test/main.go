package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
)

func main() {

	// create chrome instance
	ctx, _ :=  chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://golang.org/pkg/time/`))
	if err != nil {
		log.Fatal(err)
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
