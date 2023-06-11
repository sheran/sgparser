package sgparser

import (
	"context"
	"strings"

	"github.com/chromedp/chromedp"
)

// This is the browser package. It will use ChromeDP instead of goquery like in
// filter.go. We will attempt to write it as we did, to use the filters files

type BrowserImpl struct {
	Host         string   `toml:"host"`
	Path         string   `toml:"path"`
	Title        string   `toml:"title"`
	Body         string   `toml:"body"`
	Thumb        string   `toml:"thumb"`
	SkipChildren bool     `toml:"skip_children"`
	SkipClasses  []string `toml:"skip_classes"`
	SkipText     []string `toml:"skip_text"`
	SkipElements []string `toml:"skip_elements"`
	RSS          string   `toml:"rss"`
}

func (b *BrowserImpl) Run(urlToFetch string) string {
	opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("disable-dev-shm-usage", true))

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Create a new tab
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	var pageTitle string
	err := chromedp.Run(ctx,
		chromedp.Navigate(urlToFetch),
		chromedp.WaitVisible(b.Title),
		chromedp.Text(b.Title, &pageTitle),
	)
	if err != nil {
		return err.Error()
	}
	return pageTitle
}

func (rf *BrowserImpl) Match(host string) bool {
	return strings.Contains(rf.Host, host)
}

func (rf *BrowserImpl) Snippet(path string) bool {
	if rf.Path != "" {
		if rf.Path == "-" {
			return true
		}
		return strings.HasPrefix(path, rf.Path)
	}
	return false
}

func (rf *BrowserImpl) GetHost() string {
	return rf.Host
}
