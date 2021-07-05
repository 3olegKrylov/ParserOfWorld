package fullscreen

import (
	"github.com/chromedp/chromedp"
)

//создаёт полноразмерный скриншот экрана браузрера
func FullScreenshot( quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.FullScreenshot(res, quality),
	}
}
