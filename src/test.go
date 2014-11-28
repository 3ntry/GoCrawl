package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"strings"
	"net/url"
)

func main() {
	// Print welcome message.
	fmt.Println("Welcome to Roger's web crawler!")
	// Parse command line arguments.
	flag.Parse()
	arguments := flag.Args()
	if len(arguments) != 1 {
		fmt.Println("Usage: test URL.")
		os.Exit(1)
	}
	// The URL we want is the first argument.
	// Add 'http://' if not already included.
	url := arguments[0]
	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}
	// Echo arguments
	fmt.Println("Starting url: ", url)
	// Create queue channel.
	queue := make(chan string)
	// Start go routine.
	go func() {
		queue <- url
	}()
	for uri := range queue {
		enqueue(uri, queue)
	}
}

func enqueue(url string, queue chan string) {
	// Make http request
	response, error := http.Get(url)
	if error != nil {
		fmt.Println("Http error: ", error)
	}
	defer response.Body.Close()
	// Get links.
	links := getLinks(response.Body)
	// Enqueue links and swpan new go routine.
	for _, link := range links {
		absUrl := fixUrl(link, url)
		if absUrl != "" {
			fmt.Println(absUrl)
			go func() {
				queue <- absUrl
			}()
		}
	}
}

func fixUrl(href, base string) (string) {
  uri, err := url.Parse(href)
  if err != nil {
    return ""
  }
  baseUrl, err := url.Parse(base)
  if err != nil {
    return ""
  }
  uri = baseUrl.ResolveReference(uri)
	if strings.HasPrefix(uri.String(), "javascript:") {
		return ""
	}
  return uri.String()
}

func getLinks(reader io.Reader) []string {
	// Create new slice of strings.
	links := make([]string, 0)
	// Tokenize the html digest.
	page := html.NewTokenizer(reader)
	for {
		// Get next tag.
		tokenType := page.Next()
		// If we are finished, return.
		if tokenType == html.ErrorToken {
			return links
		}
		// Get next token.
		token := page.Token()
		// Look for 'a' tags.
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			// Get the 'http' attribute value.
			for _, attr := range token.Attr {
				// If we got one then append to the strings slice.
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
	}
}

