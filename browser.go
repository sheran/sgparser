package sgparser

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/sheran/sgparser/models"
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

func (b *BrowserImpl) Run(urlToFetch string) (*models.Post, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(
		ctx,
	)
	defer cancel()

	var res []string
	var thumb string
	var title string
	newUrl := fixAmpSuffix(urlToFetch)
	post := &models.Post{
		Url: newUrl,
	}

	if b.Thumb == "-" {
		err := chromedp.Run(ctx,
			chromedp.Navigate(newUrl),
			chromedp.WaitReady("body"),
			chromedp.Text(b.Title, &title),
			chromedp.Evaluate(fmt.Sprintf(`Array.from(document.querySelectorAll("%s")).map(i => i.innerText)`, b.Body), &res),
		)
		if err != nil {
			return nil, err
		}
	} else {
		err := chromedp.Run(ctx,
			chromedp.Navigate(newUrl),
			chromedp.WaitReady("body"),
			chromedp.Text(b.Title, &title),
			chromedp.Evaluate(fmt.Sprintf(`Array.from(document.querySelectorAll("%s")).map(i => i.innerText)`, b.Body), &res),
			chromedp.Evaluate(fmt.Sprintf(`document.querySelector("%s").src`, b.Thumb), &thumb),
		)
		if err != nil {
			return nil, err
		}
	}

	post.Body = formatTextBody(res)
	post.Thumb = fixAmpSuffix(thumb)
	post.Title = title
	return post, nil
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

func formatTextBody(res []string) string {
	// We can run further filters here to remove
	// words or phrases based on our toml config
	var s strings.Builder
	for _, node := range res {
		if len(node) > 0 {
			s.WriteString(node)
			s.WriteString("\n\n")
		}
	}
	return s.String()
}

func fixAmpSuffix(urlStr string) string {
	toCheck, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	ampSuffixes := []string{"/amp", "/amp/"}
	newUrl := fmt.Sprintf("%s://%s%s", toCheck.Scheme, toCheck.Host, toCheck.Path)
	for _, suffix := range ampSuffixes {
		newUrl = strings.TrimSuffix(newUrl, suffix)
	}
	return newUrl
}
