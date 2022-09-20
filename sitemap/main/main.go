package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	link "github.com/sixsat/gophercises/sitemap"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	flag.Parse()

	pages := get(*urlFlag)
	for _, page := range pages {
		fmt.Println(page)
	}
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, li := range links {
		if keepFn(li) {
			ret = append(ret, li)
		}
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var ret []string
	for _, li := range links {
		switch {
		case strings.HasPrefix(li.Href, "/"):
			ret = append(ret, base+li.Href)
		case strings.HasPrefix(li.Href, "./"):
			ret = append(ret, base+li.Href[1:])
		case strings.HasPrefix(li.Href, "http"):
			ret = append(ret, li.Href)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
