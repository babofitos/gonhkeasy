package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
)

var url = "http://www3.nhk.or.jp/news/easy/k10013455581000/k10013455581000.html"

func Scrape() {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	title := doc.Find("title").First().Text()
	fmt.Printf("%v\n", title)

	newsarticle := doc.Find("#newsarticle")

	var text string

	text = getTextNodes(newsarticle.Children())

	fmt.Printf("%v", text)
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

func main() {
	Scrape()
}
