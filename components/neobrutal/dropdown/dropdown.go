package dropdown

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/kheepercom/ui"
	"golang.org/x/net/html"
)

//go:embed dropdown.html
var dropdownHTML string
var dropdownTmpl = template.Must(template.New("dropdown").Parse(dropdownHTML))

//go:embed dropdown.js
var dropdownJS string

type Dropdown struct {
	log *slog.Logger
}

func New(logger *slog.Logger) *Dropdown {
	if logger == nil {
		logger = slog.Default()
	}
	return &Dropdown{
		log: logger,
	}
}

var colors = map[string]string{
	"violet": "bg-violet-200 hover:bg-violet-300",
	"pink":   "bg-pink-200 hover:bg-pink-300",
	"red":    "bg-red-200 hover:bg-red-300",
	"orange": "bg-orange-200 hover:bg-orange-300",
	"yellow": "bg-yellow-200 hover:bg-yellow-300",
	"lime":   "bg-lime-200 hover:bg-lime-300",
	"cyan":   "bg-cyan-200 hover:bg-cyan-300",
}

func (d *Dropdown) Render(r *http.Request, attrs ui.Attributes) (*html.Node, error) {
	classes := ui.Classes{
		{"inline-flex w-72 justify-center gap-x-1.5 px-3 py-2 border-black border-2 focus:outline-none focus:shadow-[2px_2px_0px_rgba(0,0,0,1)]", true},
		{colors[attrs.Get("color")], true},
		{attrs.Get("class"), true},
	}
	attrs.Set("_button_class", classes.String())

	attrs.Set("_menu_id", uuid.NewString())

	b := &bytes.Buffer{}
	if err := dropdownTmpl.Execute(b, attrs); err != nil {
		return nil, fmt.Errorf("execute dropdown template: %w", err)
	}

	return ui.ParseComponent(b)
}

func (d *Dropdown) JS() (string, error) {
	return dropdownJS, nil
}
