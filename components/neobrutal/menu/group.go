package menu

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kheepercom/ui"
	"golang.org/x/net/html"
)

//go:embed group.html
var groupHTML string
var groupTMPL = template.Must(template.New("MenuGroup").Parse(groupHTML))

type Group struct {
	log *slog.Logger
}

func NewGroup(logger *slog.Logger) *Group {
	if logger == nil {
		logger = slog.Default()
	}
	return &Group{
		log: logger,
	}
}

func (*Group) Render(r *http.Request, attrs ui.Attributes) (*html.Node, error) {
	isActive := false
	if p := attrs.Get("path"); p != "" && strings.Contains(r.URL.Path, p) {
		isActive = true
	}

	classes := ui.Classes{
		{"inliine-block hover:underline hover:underline-offset-8", true},
		{"lg:font-bold lg:underline lg:underline-offset-8", isActive},
	}
	attrs.Set("_class", classes.String())

	b := &bytes.Buffer{}
	if err := groupTMPL.Execute(b, attrs); err != nil {
		return nil, fmt.Errorf("execute MenuGroup template: %w", err)
	}

	return ui.ParseComponent(b)
}
