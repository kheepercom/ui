package leftsidebar

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

//go:embed sidebar.html
var sidebarHTML string
var sidebarTMPL = template.Must(template.New("LeftSidebar").Parse(sidebarHTML))

type LeftSidebar struct {
	log *slog.Logger
}

func New(logger *slog.Logger) *LeftSidebar {
	if logger == nil {
		logger = slog.Default()
	}
	return &LeftSidebar{
		log: logger,
	}
}

func (*LeftSidebar) Render(r *http.Request, attrs ui.Attributes) (*html.Node, error) {
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
	if err := sidebarTMPL.Execute(b, attrs); err != nil {
		return nil, fmt.Errorf("execute LeftSidebar template: %w", err)
	}

	return ui.ParseComponent(b)
}
