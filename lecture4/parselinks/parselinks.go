package parselinks

import (
	"golang.org/x/net/html"
	"io"
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
