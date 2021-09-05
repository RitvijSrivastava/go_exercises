package link_extractor

import (
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func ExtractLinks(n *html.Node) []Link {
	var links []Link
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				// texts := extractText(n)
				links = append(links, Link{a.Val, ""})
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractedLinks := ExtractLinks(c)
		links = append(links, extractedLinks...)
	}
	return links
}

func extractText(n *html.Node) string {
	var text string
	if n.Type != html.ElementNode && n.Data != "a" && n.Type != html.CommentNode {
		text = n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += extractText(c)
	}
	return strings.Trim(text, "\n")
}
