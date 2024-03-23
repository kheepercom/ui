package dropdown

import (
	"net/http"

	"github.com/kheepercom/ui"
	"golang.org/x/net/html"
)

type Item struct{}

var hovers = map[string]string{
	"violet": "hover:bg-violet-200",
	"pink":   "hover:bg-pink-200",
	"red":    "hover:bg-red-200",
	"orange": "hover:bg-orange-200",
	"yellow": "hover:bg-yellow-200",
	"lime":   "hover:bg-lime-200",
	"cyan":   "hover:bg-cyan-200",
}

func (*Item) Render(r *http.Request, attrs ui.Attributes) (*html.Node, error) {
	classes := ui.Classes{
		{"block px-4 py-2 text-sm border-black border-b-2 hover:font-medium", true},
		{hovers[attrs.Get("color")], true},
		{attrs.Get("class"), true},
	}

	a := &html.Node{
		Type: html.ElementNode,
		Data: "a",
		Attr: []html.Attribute{
			{Key: "class", Val: classes.String()},
			{Key: "id", Val: attrs.Get("id")},
			{Key: "href", Val: attrs.Get("href")},
			{Key: "role", Val: "menuitem"},
		},
	}
	label := &html.Node{
		Type:   html.TextNode,
		Data:   attrs.Get("label"),
		Parent: a,
	}
	a.FirstChild, a.LastChild = label, label

	return a, nil
}
