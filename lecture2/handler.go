package main

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if redirectURL, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, redirectURL, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
	return fn
}

type YAMLMap struct {
	Path string `yaml:"path,omitempty"`
	Url  string `yaml:"url,omitempty"`
}

func buildMap(y []YAMLMap) map[string]string {
	pathToURLs := map[string]string{}
	for i, _ := range y {
		pathToURLs[y[i].Path] = y[i].Url
	}

	return pathToURLs
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	yamlMap := new([]YAMLMap)
	err := yaml.Unmarshal(yml, &yamlMap)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(*yamlMap)

	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if redirectURL, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, redirectURL, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
	return fn, nil
}
