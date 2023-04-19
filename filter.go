package sgparser

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sheran/sgparser/models"
)

type FilterImpl struct {
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
	Doc          *goquery.Document
}

func (rf *FilterImpl) Init(urlPage string) error {
	// we always strip out the query parameters
	toCheck, err := url.Parse(urlPage)
	if err != nil {
		return err
	}
	// we check for an amp suffix
	ampSuffixes := []string{"/amp", "/amp/"}
	newUrl := fmt.Sprintf("%s://%s%s", toCheck.Scheme, toCheck.Host, toCheck.Path)
	for _, suffix := range ampSuffixes {
		newUrl = strings.TrimSuffix(newUrl, suffix)
	}

	log.Printf("fetching %s\n", newUrl)
	resp, err := http.Get(newUrl)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	rf.Doc = doc
	return nil
}

func (rf *FilterImpl) Run() (*models.Post, error) {
	title := strings.Trim(rf.Doc.Find(rf.Title).Text(), " \n\t")
	post := &models.Post{
		Title: title,
	}
	var text string
	rf.Doc.Find(rf.Body).Each(func(i int, s *goquery.Selection) {
		for _, class := range rf.SkipClasses {
			if s.HasClass(class) {
				return
			}
		}

		for _, elem := range rf.SkipElements {
			id := s.Find(elem).Each(func(i int, s1 *goquery.Selection) {
				if s1.Is(elem) {
					return
				}
			})
			if id.Length() > 0 {
				return
			}
		}

		for _, text := range rf.SkipText {
			if strings.Contains(s.Text(), text) {
				return
			}
		}
		if rf.SkipChildren {
			if s.Children().Length() > 0 {
				return
			}
		}
		if s.Text() == "" {
			return
		}

		text += strings.Trim(s.Text(), " \n\t") + "\n\n"
	})
	post.Body = text
	return post, nil
}

func (rf *FilterImpl) Match(host string) bool {
	return strings.Contains(rf.Host, host)
}

func (rf *FilterImpl) Snippet(path string) bool {
	if rf.Path != "" {
		if rf.Path == "-" {
			return true
		}
		return strings.HasPrefix(path, rf.Path)
	}
	return false
}

func (rf *FilterImpl) GetHost() string {
	return rf.Host
}
