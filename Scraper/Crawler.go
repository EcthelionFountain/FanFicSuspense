package main

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"

	"golang.org/x/net/html"
)

func main() {
	foundLinks := make(chan string)
	visited := make(map[string]bool)

	var wg sync.WaitGroup
	wg.Add(1)
	// Start crawling in a goroutine
	go crawl("https://archive.transformativeworks.org/tags/Harry%20Potter%20-%20J*d*%20K*d*%20Rowling/works?page=1", foundLinks, visited, &wg)
	go func() {
		wg.Wait()
		close((foundLinks))
	}()
	// Print whatever comes back to the channel
	for link := range foundLinks {
		fmt.Println("Found link:", link)
	}
}

func scrapetext(url string) string {
	return "nothing"
}

func scrapecomment(url string) string {
	return "nothing"
}


func crawl(url string, ch chan string, visited map[string]bool, wg *sync.WaitGroup) {
	if visited[url] {
		return
	}
	visited[url] = true

	re := regexp.MustCompile(`/works/\d+$`)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	defer wg.Done()

	tokens := html.NewTokenizer(resp.Body)

	for {
		tokenType := tokens.Next()

		switch tokenType {
		case html.ErrorToken:
			// End of the document
			return

		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokens.Token()

			// Check if the tag is an anchor <a>
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := attr.Val

						if re.MatchString(link) {
							ch <- link
						}
					}
				}
			}
		}
	}
}
