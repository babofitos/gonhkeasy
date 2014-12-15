package nhkeasy

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type NHKScraper struct {
	url       string
	separator string
	doc       *goquery.Document
}

func New(url, separator string) *NHKScraper {
	s := new(NHKScraper)
	s.url = url
	s.separator = separator
	doc, err := s.getDocument()
	if err != nil {
		panic(err)
	}
	s.doc = doc
	return s
}

func (s *NHKScraper) Scrape() (string, string) {
	title := s.doc.Find("title").Text()
	title = prettyTitle(title)
	var text string
	newsarticle := s.doc.Find("#newsarticle")
	newsarticle.Children().Each(func(i int, paragraph *goquery.Selection) {
		text += getTextNodes(paragraph)
		if paragraph.Next().Children().Length() > 0 {
			text += s.separator
		}
	})
	return title, text
}

func (s *NHKScraper) request() (*http.Response, error) {
	resp, err := http.Get(s.url)
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

func (s *NHKScraper) getDocument() (*goquery.Document, error) {
	resp, err := s.request()
	if err != nil {
		return nil, err
	}
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
