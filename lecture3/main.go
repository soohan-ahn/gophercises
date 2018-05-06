package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type StoryHtml struct {
	StoryArc
	Title string
}

type StoryArc struct {
	Title   string        `json:"title,omitempty"`
	Story   []string      `json:"story,omitempty"`
	Options []StoryOption `json:"options,string,omitempty"`
}

type StoryOption struct {
	Text string `json:"text,omitempty"`
	Arc  string `json:"arc,omitempty"`
}

func parseJSON() (stories map[string]StoryArc) {
	jsonContent, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonContent, &stories)
	if err != nil {
		panic(err)
	}

	return stories
}

func router(stories map[string]StoryArc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Base(r.URL.Path)
		fmt.Printf("Path: %s\n", path)
		if _, ok := stories[path]; !ok {
			path = "intro"
		}

		t, err := template.ParseFiles("template.html")
		if err != nil {
			panic(err)
		}
		s := StoryHtml{
			StoryArc: stories[path],
			Title:    path,
		}
		t.Execute(w, s)
		//fmt.Fprint(w, stories[path])
	}

	return fn
}

func main() {
	mainHandler := router(parseJSON())
	fmt.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mainHandler))
}
