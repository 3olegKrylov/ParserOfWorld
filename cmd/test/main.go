package main

import (
	"context"
	"flag"
	"github.com/chromedp/chromedp"
	"log"
)


func main() {

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var devToolWsUrl string
	flag.StringVar(&devToolWsUrl, "devtools-ws-url", "ws://127.0.0.1:42273/devtools/browser/81ac8cbd-b85f-40a2-bb79-bec68d6990d7", "DevTools Websocket URL")
	flag.Parse()

	actxt, cancelActxt := chromedp.NewRemoteAllocator(context.Background(), devToolWsUrl)
	defer cancelActxt()

	ctx, cancelCtxt := chromedp.NewContext(actxt) // create new tab
	defer cancelCtxt()                             // close tab afterwards


	example:=""
	if err := chromedp.Run(ctx,
		chromedp.Navigate(`https://golang.org/pkg/time/`),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`body > footer`),
		// find and click "Expand All" link
		chromedp.Click(`#pkg-examples > div`, chromedp.NodeVisible),
		// retrieve the value of the textarea
		chromedp.Value(`#example_After .play .input textarea`, &example),
	); err != nil {
		log.Fatalf("Failed: %v", err)
	}


	log.Printf("Go's time.After example:\n%s", example)
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

