package sgparser

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestBareToml(t *testing.T) {
	list := LoadToml("filters")
	url := "https://www.theonlinecitizen.com/2023/04/20/singaporean-activist-and-human-rights-lawyer-raise-concerns-over-impending-execution-amid-troubling-case-detail/"
	bodytoml, err := Process(url, list)
	if err != nil {
		panic(err)
	}
	if bodytoml != nil {
		fmt.Println(bodytoml.Title)
		fmt.Println(bodytoml.Body)
	} else {
		fmt.Println("empty body")
	}
}

func TestCDP(t *testing.T) {
	list := LoadCDP("filters")
	// Open the text file
	file, err := os.Open("testurls.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create an empty array to store the lines
	var lines []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Add each line to the array
		lines = append(lines, scanner.Text())
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Print the lines
	for _, line := range lines {
		output, err := Browse(line, list)
		if err != nil {
			log.Printf("errur %s\n", err.Error())
			continue
		}
		fmt.Println(output)
	}

}

func TestHTML(t *testing.T) {
	resp, err := http.Get("https://www.planetf1.com/news/charles-leclerc-grid-penalty-saudi-arabia/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	var pattern strings.Builder
	doc.Find("div.ciam-article-pf1 > p").Each(func(i int, s *goquery.Selection) {
		if s.Children().Length() > 3 {
			for _, child := range s.Children().Nodes {
				pattern.WriteString(child.Data)
			}
		}
	})
	fmt.Println(pattern.String())
}
