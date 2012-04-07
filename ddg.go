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
//		"fmt"
//		"github.com/whee/ddg"
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
	Height int        `json:"-"` // height of icon (px)
	Width  int        `json:"-"` // width of icon (px)

	// The height and width can be "" (string; we treat as 0) or an int. Unmarshal here, then populate the above two.
	RawHeight interface{} `json:"Height"`
	RawWidth interface{} `json:"Width"`
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

// A Client is a DDG Zero-click client.
type Client struct {
	// Secure specifies whether HTTPS is used.
	Secure bool
}

// ZeroClick queries DuckDuckGo's zero-click API for the specified query
// and returns the Response.
func (c *Client) ZeroClick(query string) (res Response, err error) {
	// TODO: Support some of the available configuration (e.g., no html)
	v := url.Values{}
	v.Set("q", query)
	v.Set("format", "json")
	// TODO: support the disambiguation category type
	v.Set("skip_disambig", "1")

	var scheme string
	if c.Secure {
		scheme = "https"
	} else {
		scheme = "http"
	}

	req, err := http.NewRequest("GET", scheme+"://api.duckduckgo.com/?"+v.Encode(), nil)
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

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err == io.EOF {
		err = nil
	}
	handleInterfaces(&res)

	return
}

// handleInterfaces cleans up incoming data that can be of multiple types; for example, the
// icon width and height are either a float64 or string, but we want to treat them as int.
func handleInterfaces(response *Response) {
	for _, result := range response.Results {
		if height, ok := result.Icon.RawHeight.(float64); ok {
			result.Icon.Height = int(height)
		}
		if width, ok := result.Icon.RawWidth.(float64); ok {
			result.Icon.Width = int(width)
		}
	}
}

// ZeroClick queries DuckDuckGo's zero-click API for the specified query
// and returns the Response. This helper function uses a zero-value Client.
func ZeroClick(query string) (res Response, err error) {
	c := &Client{}
	return c.ZeroClick(query)
}
