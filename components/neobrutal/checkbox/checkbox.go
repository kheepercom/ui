package checkbox

import (
	"net/http"

	"github.com/kheepercom/ui"
	"golang.org/x/net/html"
)

type Checkbox struct{}

var colors = map[string]string{
	"violet": "before:bg-violet-200 before:checked:bg-violet-300",
	"pink":   "before:bg-pink-200 before:checked:bg-pink-300",
	"red":    "before:bg-red-200 before:checked:bg-red-300",
	"orange": "before:bg-orange-200 before:checked:bg-orange-300",
	"yellow": "before:bg-yellow-200 before:checked:bg-yellow-300",
	"lime":   "before:bg-lime-200 before:checked:bg-lime-300",
	"cyan":   "before:bg-cyan-200 before:checked:bg-cyan-300",
}

var sizes = map[string]string{
	"sm": "w-5 h-5 before:w-5 before:h-5 before:border-2 before:hover:shadow-[2px_2px_0px_rgba(0,0,0,1)] before:checked:shadow-[2px_2px_0px_rgba(0,0,0,1)] after:left-1.5 after:top-0.5 after:w-2 after:h-3 after:border-r-2 after:border-b-2",
	"lg": "w-10 h-10 before:w-8 before:h-8 before:border-4 before:hover:shadow-[4px_4px_0px_rgba(0,0,0,1)] before:checked:shadow-[4px_4px_0px_rgba(0,0,0,1)] after:left-2.5 after:top-1.5 after:w-3 after:h-4 after:border-r-4 after:border-b-4",
}

func (*Checkbox) Render(_ *http.Request, attrs ui.Attributes) (*html.Node, error) {
	classes := ui.Classes{
		{"apeearance-none outline-none block relative text-center cursor-pointer m-auto before:rounded-sm before:block before:absolute before:content-[''] before:bg-[#FFC29F] before:rounded-sm before:border-black after:block after:content-[''] after:absolute after:border-black after:origin-center after:rotate-45 ", true},
		{sizes[attrs.GetOr("size", "sm")], true},
		{"[&:not(:checked)]:after:opacity-0", true},
		{"after:checked:opacity-1 before:checked:bg-[#FF965B]", true},
		{colors[attrs.GetOr("color", "orange")], true},
		// Custom classes set on the component
		{attrs.Get("class"), true},
	}

	a := []html.Attribute{
		{Key: "type", Val: "checkbox"},
		{Key: "class", Val: classes.String()},
	}
	if attrs.Has("checked") {
		a = append(a, html.Attribute{Key: "checked"})
	}

	input := &html.Node{
		Type: html.ElementNode,
		Data: "input",
		Attr: a,
	}
	label := &html.Node{
		Type:       html.ElementNode,
		Data:       "label",
		FirstChild: input,
		LastChild:  input,
	}
	input.Parent = label

	return label, nil
}
