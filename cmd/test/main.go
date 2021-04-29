package main

import (
	"fmt"
	"github.com/testSpace/internal/db"
)

func main() {

	dbConnect := db.DBconnect()
	db.DBinit(dbConnect)
	defer dbConnect.Close()

	fmt.Println(db.FindUserDB(dbConnect, "nargissa28048330"))

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
