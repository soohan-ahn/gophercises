package main

import (
	"./parselinks"
	"fmt"
	"os"
)

func main() {
	//r, err := ioutil.ReadFile(os.Args[1])
	r, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	result := parselinks.RetrieveAllLinks(r)
	for i, _ := range result {
		fmt.Printf("HREF: %s\nText: %s\n\n", result[i].Href, result[i].Text)
	}
}
