// Copyright 2012, Brian Hetro <whee@smaertness.net>
// Use of this source code is governed by the ISC license
// that can be found in the LICENSE file.

// DuckDuckGo Zero-click API client. See https://duckduckgo.com/api.html
// for information about the API.
//
// Example command-line program to show abstracts:
//
//	package main
//	
//	import (
//		"ddg"
//		"fmt"
//		"os"
//	)
//
//	func main() {
//		for _, s := range os.Args[1:] {
//			if r, err := ddg.ZeroClick(s); err == nil {
//				fmt.Printf("%s: %s\n", s, r.Abstract)
//			} else {
//				fmt.Printf("Error looking up %s: %v\n", s, err)
//			}
//		}
//	}
package ddg

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Response represents the response from a zero-click API request via
// the ZeroClick function.
type Response struct {
	Abstract       string // topic summary (can contain HTML, e.g. italics)
	AbstractText   string // topic summary (with no HTML)
	AbstractSource string // name of Abstract source
	AbstractURL    string // deep link to expanded topic page in AbstractSource
	Image          string // link to image that goes with Abstract
	Heading        string // name of topic that goes with Abstract

	Answer     string // instant answer
	AnswerType string // type of Answer, e.g. calc, color, digest, info, ip, iploc, phone, pw, rand, regexp, unicode, upc, or zip (see goodies & tech pages for examples).

	Definition       string // dictionary definition (may differ from Abstract)
	DefinitionSource string // name of Definition source
	DefinitionURL    string // deep link to expanded definition page in DefinitionSource

	RelatedTopics []Result // array of internal links to related topics associated with Abstract

	Results []Result // array of external links associated with Abstract

	Type CategoryType // response category, i.e. A (article), D (disambiguation), C (category), N (name), E (exclusive), or nothing.

	Redirect string // !bang redirect URL
}

type Result struct {
	Result   string // HTML link(s) to external site(s)
	FirstURL string // first URL in Result
	Icon     Icon   // icon associated with FirstURL
	Text     string // text from FirstURL
}

type Icon struct {
	URL    string      // URL of icon
	Height interface{} // height of icon (px)
	Width  interface{} // width of icon (px)
}

type CategoryType string

const (
	Article        = "A"
	Disambiguation = "D"
	Category       = "C"
	Name           = "N"
	Exclusive      = "E"
	None           = ""
)

// ZeroClick queries DuckDuckGo's zero-click API for the specified query
// and returns the Response.
func ZeroClick(query string) (res Response, err error) {
	// TODO: Support some of the available configuration (e.g., no html)
	v := url.Values{}
	v.Set("q", query)
	v.Set("format", "json")
	// TODO: support the disambiguation category type
	v.Set("skip_disambig", "1")

	req, err := http.NewRequest("GET", "https://api.duckduckgo.com/?" + v.Encode(), nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "ddg.go/0.1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(&res)
	if err == io.EOF {
		err = nil
	}

	return
}
