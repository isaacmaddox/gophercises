package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func main() {
	fileFlag := flag.String("file", "./ex1.html", "Which file to parse using the reader. Default is ./ex1.html")
	flag.Parse()

	file, fileReadErr := os.ReadFile(*fileFlag)

	if fileReadErr != nil {
		log.Fatalf("could not read the file %s", *fileFlag)
	}

	root, htmlParseErr := html.Parse(strings.NewReader(string(file)))

	if htmlParseErr != nil {
		log.Fatalf("an error occurred parsing the file. check the syntax")
	}

	linkList := []Link{}
	recurse(root, &linkList)

	fmt.Printf("%+v", linkList)
}

func recurse(node *html.Node, list *[]Link) {
	for n := range node.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			link := Link{}
			text := &bytes.Buffer{}
			collectText(n, text)

			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link.Href = attr.Val
				}
			}

			link.Text = strings.Trim(text.String(), " \n")
			*list = append(*list, link)
		}
	}
}

func collectText(n *html.Node, buf *bytes.Buffer) {
	if n.Type == html.TextNode {
		buf.WriteString(n.Data)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		collectText(c, buf)
	}
}
