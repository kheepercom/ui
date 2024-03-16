package ui

import (
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
