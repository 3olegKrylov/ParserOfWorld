package main

import (
	"fmt"
	"sync"
)

func main() {

	mesChan := make(chan int)

	wg:= sync.WaitGroup{}
	wg.Add(100)
	go func(chan int) {
		for {
			res := <-mesChan
			fmt.Println("First gorutine handle: ", res)

			wg.Done()
		}
	}(mesChan)

	go func(chan int) {
		for {
			res := <-mesChan
			fmt.Println("Second gorutine handle: ", res)

			wg.Done()
		}
	}(mesChan)

	go func(chan int) {
		for {
			res := <-mesChan
			fmt.Println("Third gorutine handle: ", res)

			wg.Done()
		}
	}(mesChan)


	for i:=0 ; i<100; i++{
		mesChan <- i
		fmt.Println("записал ", i)
	}

	wg.Wait()

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
