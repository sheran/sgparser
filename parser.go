package sgparser

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/sheran/sgparser/models"
)

type Filter interface {
	Init(string) error
	Run() (*models.Post, error)
	Match(string) bool
	Snippet(string) bool
	GetHost() string
}

type Browser interface {
	Run(string) (*models.Post, error)
	Match(string) bool
	Snippet(string) bool
	GetHost() string
}

func LoadToml(dir string) []Filter {
	log.Println("loading filters...")
	list := make([]Filter, 0)

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range entries {
		if filepath.Ext(file.Name()) == ".toml" {
			var fl FilterImpl
			log.Printf("[+] %s/%s\n", dir, file.Name())
			_, err := toml.DecodeFile(fmt.Sprintf("%s/%s", dir, file.Name()), &fl)
			if err != nil {
				panic(err)
			}
			list = append(list, &fl)
		}
	}
	return list
}

func LoadCDP(dir string) []Browser {
	log.Println("loading filters...")
	list := make([]Browser, 0)

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range entries {
		if filepath.Ext(file.Name()) == ".toml" {
			var bl BrowserImpl
			log.Printf("[+] %s/%s\n", dir, file.Name())
			_, err := toml.DecodeFile(fmt.Sprintf("%s/%s", dir, file.Name()), &bl)
			if err != nil {
				panic(err)
			}
			list = append(list, &bl)
		}
	}
	return list
}

func Browse(page string, list []Browser) (*models.Post, error) {
	link, err := url.ParseRequestURI(page)
	if err != nil {
		return nil, err
	}
	var output *models.Post
	for _, filter := range list {
		if link.Host != "" {
			if filter.Match(link.Host) {
				log.Printf("got host: %s\n", link.Host)
				if filter.Snippet(link.Path) {
					output, err = filter.Run(page)
					if err != nil {
						return nil, err
					}
				}
			}
		}

	}
	return output, nil
}

func Process(page string, list []Filter) (*models.Post, error) {
	link, err := url.ParseRequestURI(page)
	if err != nil {
		return nil, err
	}
	var body *models.Post
	for _, filter := range list {
		if link.Host != "" {
			if filter.Match(link.Host) {
				log.Printf("got host: %s\n", link.Host)
				if filter.Snippet(link.Path) {
					filter.Init(page)
					log.Printf("[+] (%s) running filter\n", filter.GetHost())
					body, err = filter.Run()
					if err != nil {
						return nil, err
					}
					if len(body.Body) < 350 {
						return nil, fmt.Errorf("[!!](%s) parsed text is too short to post", filter.GetHost())
					}
					log.Printf("[+] (%s) text length: %d", filter.GetHost(), len(body.Body))
				}
			}
		}

	}
	return body, nil
}
