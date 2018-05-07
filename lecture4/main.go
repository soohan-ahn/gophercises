package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"os"
	"strings"
)

type Link struct {
	Href string
	Text string
}

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

func main() {
	//r, err := ioutil.ReadFile(os.Args[1])
	r, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	result := RetrieveAllLinks(r)
	for i, _ := range result {
		fmt.Printf("HREF: %s\nText: %s\n\n", result[i].Href, result[i].Text)
	}
}
