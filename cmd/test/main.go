package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/testSpace/internal/parsing"
	"log"
)

func main() {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	parsing.ParsingAccountData("e.zhabrov", ctx, 1)
}
