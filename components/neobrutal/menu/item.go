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

//go:embed item.html
var itemHTML string
var itemTMPL = template.Must(template.New("MenuItem").Parse(itemHTML))

type Item struct {
	log *slog.Logger
}

func NewItem(logger *slog.Logger) *Item {
	if logger == nil {
		logger = slog.Default()
	}
	return &Item{
		log: logger,
	}
}

func (*Item) Render(r *http.Request, attrs ui.Attributes) (*html.Node, error) {
	isActive := false
	if p := attrs.Get("href"); p != "" && strings.Contains(r.URL.Path, p) {
		isActive = true
	}

	classes := ui.Classes{
		{"inliine-block hover:underline hover:underline-offset-8", true},
		{"lg:font-bold lg:underline lg:underline-offset-8", isActive},
	}
	attrs.Set("_class", classes.String())

	b := &bytes.Buffer{}
	if err := itemTMPL.Execute(b, attrs); err != nil {
		return nil, fmt.Errorf("execute MenuItem template: %w", err)
	}
	return ui.ParseComponent(b)
}
