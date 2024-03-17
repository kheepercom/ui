package ui

import (
	"embed"
	"net/http"

	"golang.org/x/net/html"
)

type Component interface {
	Render(r *http.Request, attributes Attributes) (*html.Node, error)
}

// If a Component implements the Styled interface then its CSS will
// be added to the app's CSS bundle served with every page.
type Styled interface {
	CSS() (string, error)
}

// If a Component implements the Scripted interface then its JS will
// be added to the app's JS bundle served with every page.
type Scripted interface {
	JS() (string, error)
}

// Examples is an optional component interface used to add variants of the
// component to the catalog. It is also used to ensure that all classes the
// component makes use of are included in the Tailwind release build.
type Examples interface {
	Examples() embed.FS
}
