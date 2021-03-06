package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/RitvijSrivastava/go_exercises/sitemap_builder/link_extractor"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Url string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {

	url := flag.String("url", "https://gophercises.com", "url for sitemap builder")
	maxDepth := flag.Int("depth", 3, "how deep should the builder go")

	flag.Parse()

	links := bfs(*url, *maxDepth)

	xmlData := urlset{Xmlns: xmlns}
	for _, link := range links {
		xmlData.Urls = append(xmlData.Urls, loc{link})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", " ")
	if err := enc.Encode(xmlData); err != nil {
		panic(err)
	}
	fmt.Println()
}

func bfs(baseUrl string, maxDepth int) []string {
	visited := make(map[string]struct{})
	q := make(map[string]struct{})
	q[baseUrl] = struct{}{}

	for depth := 0; depth <= maxDepth; depth++ {

		if len(q) == 0 {
			break
		}

		tmpq := make(map[string]struct{})
		for link, _ := range q {
			if _, ok := visited[link]; ok {
				continue
			}
			visited[link] = struct{}{}
			links := get(link)
			for _, url := range links {
				tmpq[url] = struct{}{}
			}
		}
		// Move to the new layer
		q = tmpq
	}

	links := make([]string, 0)
	for url, _ := range visited {
		links = append(links, url)
	}
	return links
}

func get(link string) []string {
	resp, err := http.Get(link)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	return filter(hrefs(resp.Body, baseURL.String()), withPrefix(baseURL.String()))
}

func hrefs(r io.Reader, baseURL string) []string {
	links := link_extractor.Parse(r)
	var urls []string

	for _, link := range links {
		if strings.HasPrefix(link.Href, "/") {
			urls = append(urls, baseURL+link.Href)
		} else if strings.HasPrefix(link.Href, "http") {
			urls = append(urls, link.Href)
		}
	}
	return urls
}

func filter(links []string, filterFn func(string) bool) []string {
	var urls []string
	for _, link := range links {
		if filterFn(link) {
			urls = append(urls, link)
		}
	}
	return urls
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
