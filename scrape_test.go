package nhkeasy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScrape(t *testing.T) {
	dat, err := ioutil.ReadFile("mock.html")
	if err != nil {
		panic(err)
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, string(dat))
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()
	separator := "\n"
	n := New(ts.URL, separator)
	title, text := n.Scrape()
	if title != "長野県で震度６弱の地震　これからも十分に注意して" {
		t.Error("Bad title")
	}
	if text != fmt.Sprintf("２２日午後１０時ころ、長野県でマグニチュード%sこの地震で４５人がけがをしたと、２５", separator) {
		t.Error("Bad text")
	}
}
