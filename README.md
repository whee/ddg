# ddg - A Go package for DuckDuckGo Zero-click API access

## Install

	go get github.com/whee/ddg

## Use

	package main
	
	import (
		"ddg"
		"fmt"
		"os"
	)

	func main() {
		for _, s := range os.Args[1:] {
			if r, err := ddg.ZeroClick(s); err == nil {
				fmt.Printf("%s: %s\n", s, r.Abstract)
			} else {
				fmt.Printf("Error looking up %s: %v\n", s, err)
			}
		}
	}

See [godoc documentation](http://www.gopkgdoc.com/pkg/github.com/whee/ddg) for more information.
