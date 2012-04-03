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

	expected := "DuckDuckGo is an Internet search engine."
	if res.Abstract != expected {
		t.Errorf("got %q, expected %q\n", res.Abstract, expected)
	}
}
