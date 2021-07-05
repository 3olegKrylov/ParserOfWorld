   package parsing

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"
)

//Даёт определённое время на опрерации chromedp
func RunWithTimeOut(timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutContext, cancel := context.WithTimeout(ctx, timeout*time.Second)

		defer cancel()

		return tasks.Do(timeoutContext)
	}
}
