package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/RitvijSrivastava/go_exercises/sitemap_builder/link_extractor"
	"golang.org/x/net/html"
)

const domain string = "https://gophercises.com"

// GET the response and return the parsed HTML
func parseLink(link string) (*html.Node, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Check if a link is valid or not
// Return a valid path to the link, if the link is valid
func isValidLink(link, domain string) (bool, string) {
	// Remove last '/' from domain
	if domain[len(domain)-1] == '/' {
		domain = domain[:len(domain)-1]
	}

	if link[0] == '/' || strings.HasPrefix(link, domain) {
		if link[0] == '/' {
			link = domain + "" + link
		}
		return true, link
	}
	return false, ""
}

// Extract "hrefs" from a HTML Page, and
// return a list of valid links.
func extractLinks(link, domain string) []string {
	doc, err := parseLink(link)
	if err != nil {
		panic(err)
	}

	// Extract all links from a page
	extracted_links := link_extractor.ExtractLinks(doc)

	var valid_links []string
	for _, link := range extracted_links {
		if ok, val := isValidLink(link.Href, domain); ok {
			valid_links = append(valid_links, val)
		}
	}
	return valid_links
}

// Get all links from a domain
func getLinks(domain string) map[string]string {
	set := make(map[string]string)

	var queue []string
	queue = append(queue, domain)

	for len(queue) > 0 {
		front := queue[0]
		queue = queue[1:]

		// fmt.Println("Processing  ", front)

		valid_links := extractLinks(front, front)
		for _, link := range valid_links {
			if _, exists := set[link]; !exists {
				set[link] = front
				queue = append(queue, link)
			}
		}
	}
	return set
}

func main() {
	links := getLinks(domain)
	for link := range links {
		fmt.Println(link)
	}

	// TODO: Improve Base URL
	// TODO: Add XML
}
