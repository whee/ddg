// Copyright 2012, Brian Hetro <whee@smaertness.net>
// Use of this source code is governed by the ISC license
// that can be found in the LICENSE file.

package ddg

import (
	"testing"
)

func TestZeroClick(t *testing.T) {
	res, err := ZeroClick("DuckDuckGo")
	if err != nil {
		t.Fatal(err)
	}
	verifyDDG(t, &res)
}

func TestClient(t *testing.T) {
	c := &Client{}
	res, err := c.ZeroClick("DuckDuckGo")
	if err != nil {
		t.Fatal(err)
	}
	verifyDDG(t, &res)

	c = &Client{Secure: true}
	res, err = c.ZeroClick("DuckDuckGo")
	if err != nil {
		t.Fatal(err)
	}
	verifyDDG(t, &res)
}

func verifyDDG(t *testing.T, r *Response) {
	expected := "DuckDuckGo is an Internet search engine."
	if r.Abstract != expected {
		t.Errorf("got %q, expected %q\n", r.Abstract, expected)
	}
}
