package ui

import (
	"net/http"

	"golang.org/x/net/html"
)

type Component interface {
	Render(r *http.Request, attributes Attributes) (*html.Node, error)
}

type Styled interface {
	CSS() (string, error)
}

type Scripted interface {
	JS() (string, error)
}
