package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type QueueData struct {
	Depth int
	Path  string
}

type Link struct {
	Href string
	Text string
}

var maxDepth = 3
var queueSize = 100

func RetrieveAllLinks(f io.ReadCloser) []Link {
	z := html.NewTokenizer(f)
	var result []Link
	for {
		t := z.Next()
		if t == html.ErrorToken {
			break
		}
		tk := z.Token()

		if tk.Data == "a" {
			var url string
			for _, v := range tk.Attr {
				if v.Key == "href" {
					url = v.Val
				}
			}
			l := Link{
				Href: url,
			}
			for {
				t := z.Next()
				if t == html.ErrorToken {
					break
				}
				tk := z.Token()
				if t == html.EndTagToken && tk.Data == "a" {
					break
				}
				if t == html.TextToken {
					l.Text += strings.Trim(tk.Data, "\n ")
				}
			}

			result = append(result, l)
		}
	}

	return result
}

func retrieve(c chan QueueData, hostname string) map[string]int {
	memo := make(map[string]int)
	for {
		select {
		case q, ok := <-c:
			fmt.Printf("Q: %s\n", q.Path)
			if !ok {
				return memo
			}
			if q.Depth >= maxDepth {
				continue
			}

			nextURL := hostname + q.Path
			res, err := http.Get(nextURL)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()
			links := RetrieveAllLinks(res.Body)

			for i, _ := range links {
				if matched, err := regexp.Match(hostname, []byte(links[i].Href)); matched == true && err == nil {
					// Trim Hostname.
					path := strings.TrimPrefix(links[i].Href, hostname)
					fmt.Printf("key: %s, Link: %s\n", links[i].Href, path)

					// update depth map.
					if v, ok := memo[path]; !ok || v > q.Depth+1 {
						memo[path] = q.Depth + 1
						nextData := QueueData{
							Depth: memo[path],
							Path:  path,
						}
						c <- nextData
					}
				} else if err != nil {
					panic(err)
				}
			}
		default:
			return memo
		}
	}

	return memo
}

func main() {
	inputURL := os.Args[1]
	u, err := url.Parse(inputURL)
	if err != nil {
		panic(err)
	}

	hostname := strings.TrimRight(inputURL, u.RequestURI())
	fmt.Printf("Host: %s\n", hostname)

	c := make(chan QueueData, queueSize)
	defer close(c)
	initData := QueueData{
		Depth: 0,
		Path:  u.RequestURI(),
	}
	c <- initData
	memo := retrieve(c, hostname)

	for key, value := range memo {
		if value <= maxDepth {
			// Todo:: Make XML
			fmt.Printf("Link: %s\n", key)
		}
	}
}
