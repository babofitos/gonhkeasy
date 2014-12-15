package nhkeasy

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

func Scrape(url, separator string) (string, string, error) {
	resp, err := request(url)
	if err != nil {
		return "", "", err
	}
	doc, err := getDocument(resp)
	if err != nil {
		return "", "", err
	}
	title := doc.Find("title").Text()
	title = prettyTitle(title)
	var text string
	newsarticle := doc.Find("#newsarticle")
	newsarticle.Children().Each(func(i int, paragraph *goquery.Selection) {
		text += getTextNodes(paragraph)
		if paragraph.Next().Children().Length() > 0 {
			text += separator
		}
	})
	return title, text, err
}

func getDocument(resp *http.Response) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	defer resp.Body.Close()
	if err != nil {
		err = errors.New("Could not parse response body")
		return nil, err
	}
	return doc, nil
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

func request(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		err = errors.New("HTTP Protocol Error")
		return nil, err
	}
	if resp.StatusCode == 404 {
		err = errors.New("Status code 404. Check URL")
		return nil, err
	}
	return resp, nil
}
