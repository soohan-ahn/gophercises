package main

import (
	"fmt"
	"github.com/gophercises/lecture4/parselinks"
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

var maxDepth = 3
var queueSize = 100

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
			links := parselinks.RetrieveAllLinks(res.Body)

			for i, _ := range links {
				if matched, err := regexp.Match(hostname, []byte(links[i].Href)); err != nil {
					panic(err)
				} else if matched {
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
