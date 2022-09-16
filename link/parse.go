package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represents an anchor tag in an HTML document.
type Link struct {
	Href string
	Text string
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
			break
		}
	}
	ret.Text = getText(n)
	return ret
}

// getText extract all of the text inside the HTML node and return them.
func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	} else if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += getText(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}

// linkNode return a slice of HTML node containing the a tag.
func linkNode(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNode(c)...)
	}
	return ret
}

// Parse takes in an HTML document and returns a slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	nodes := linkNode(doc)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}
