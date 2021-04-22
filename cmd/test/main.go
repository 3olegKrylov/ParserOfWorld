package main

import (
	"fmt"
	"sync"
)

func main() {
	wg:= sync.WaitGroup{}
	wg.Add(1000)

	print:=make(chan int)


	for i:= 0; i<1000;i++{
		go func() {
			print <- i
			wg.Done()
		}()
	}

	go func(chan int) {
		total := 0

		for{
			var newPrin int
			newPrin = <- print
			total ++
			fmt.Println("total: ", total, " Value: ", newPrin)
		}
	}(print)

	wg.Wait()
	//start := time.Now()
	//
	//
	//	ctx, cancel := chromedp.NewContext(
	//		context.Background(),
	//		chromedp.WithLogf(log.Printf),
	//	)
	//	defer cancel()
	//
	//	parsing.ParsingAccountData("e.zhabrov", ctx, 1)
	//	wg.Done()
	//
	//
	//
	//	ctx, cancel := chromedp.NewContext(
	//		context.Background(),
	//		chromedp.WithLogf(log.Printf),
	//	)
	//	defer cancel()
	//
	//	parsing.ParsingAccountData("eyeyboo", ctx, 1)
	//	wg.Done()
	//
	//
	//fmt.Println(time.Since(start))
}
