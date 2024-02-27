package ui

import (
	"errors"
	"io"

	"golang.org/x/net/html"
)

// ParseComponent removes the doc, html, body wrapper added by net/html.Parse.
// If you call html.Parse("<MyComponent></MyComponent>") it will return
// "<doc><html><head></head><body><MyComponent></MyComponent></body></html>"
func ParseComponent(r io.Reader) (*html.Node, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	node := doc.FirstChild.LastChild.FirstChild
	if node == nil {
		return nil, errors.New("node not parsed")
	}
	node.Parent = nil

	return node, nil
}

var ParsePage = html.Parse
