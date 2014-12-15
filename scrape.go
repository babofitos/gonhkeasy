package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

func Scrape(url, separator string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	title := doc.Find("title").Text()
	var text string

	fmt.Println(prettyTitle(title))

	newsarticle := doc.Find("#newsarticle")
	newsarticle.Children().Each(func(i int, paragraph *goquery.Selection) {
		if paragraph.Children().Length() > 0 {
			text += getTextNodes(paragraph)
			text += separator
		}
	})

	fmt.Print(text)
}

func getTextNodes(s *goquery.Selection) string {
	var output string

	nodes := s.Contents()

	if nodes.Length() < 1 {
		//base case, text node
		return s.Text()
	}
	nodes.Each(func(i int, node *goquery.Selection) {
		if !node.Is("rt") {
			output += getTextNodes(node)
		}
	})
	return output
}

func prettyTitle(title string) string {
	return strings.Split(title, "|")[1]
}

func main() {
	var (
		url       string
		separator string
	)
	fmt.Print("Input URL:")
	fmt.Scan(&url)
	fmt.Print("Input Separator:")
	fmt.Scan(&separator)
	Scrape(url, separator)
}
