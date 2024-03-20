package button

import (
	"net/http"

	"github.com/kheepercom/ui"
	"golang.org/x/net/html"
)

type Button struct{}

func (b *Button) Render(_ *http.Request, attrs ui.Attributes) (*html.Node, error) {
	classes := ui.Classes{
		{"border-black border-2", true},
		{
			"bg-violet-200 hover:bg-violet-300 active:bg-violet-400",
			attrs.Get("color") == "violet" && !attrs.Has("disabled"),
		},
		{
			"bg-pink-200 hover:bg-pink-300 active:bg-pink-400",
			attrs.Get("color") == "pink" && !attrs.Has("disabled"),
		},
		{
			"bg-red-200 hover:bg-red-300 active:bg-red-400",
			attrs.Get("color") == "red" && !attrs.Has("disabled"),
		},
		{
			"bg-orange-200 hover:bg-orange-300 active:bg-orange-400",
			attrs.Get("color") == "orange" && !attrs.Has("disabled"),
		},
		{
			"bg-yellow-200 hover:bg-yellow-300 active:bg-yellow-400",
			attrs.Get("color") == "yellow" && !attrs.Has("disabled"),
		},
		{
			"bg-lime-200 hover:bg-lime-300 active:bg-lime-400",
			attrs.Get("color") == "lime" && !attrs.Has("disabled"),
		},
		{
			"bg-cyan-200 hover:bg-cyan-300 active:bg-cyan-400", (attrs.Get("color") == "cyan" || !attrs.Has("color")) && !attrs.Has("disabled"),
		},
		{"rounded-none", attrs.Get("rounded") == "none" || !attrs.Has("rounded")},
		{"rounded-md", attrs.Get("rounded") == "md"},
		{"rounded-full", attrs.Get("rounded") == "full"},
		{"h-10 px-4 hover:shadow-[2px_2px_0px_rgba(0,0,0,1)]", attrs.Get("size") == "sm"},
		{"h-12 px-5 hover:shadow-[2px_2px_0px_rgba(0,0,0,1)]", attrs.Get("size") == "md" || !attrs.Has("size")},
		{"h-14 px-5 hover:shadow-[4px_4px_0px_rgba(0,0,0,1)]", attrs.Get("size") == "lg"},
		{
			"border-[#727272] bg-[#D4D4D4] text-[#676767] hover:bg-[#D4D4D4] hover:shadow-none active:bg-[#D4D4D4]",
			attrs.Has("disabled"),
		},

		// Custom classes set on the component
		{attrs.Get("class"), true},
	}

	a := []html.Attribute{
		{Key: "class", Val: classes.String()},
	}
	if attrs.Has("disabled") {
		a = append(a, html.Attribute{Key: "disabled"})
	}

	children := &html.Node{
		Type: html.ElementNode,
		Data: "children",
	}
	btn := &html.Node{
		Type:       html.ElementNode,
		Data:       "button",
		Attr:       a,
		FirstChild: children,
		LastChild:  children,
	}
	children.Parent = btn
	return btn, nil
}
