package dialog

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/kheepercom/ui"
	"golang.org/x/net/html"
)

//go:embed dialog.html
var dialogHTML string
var dialogTmpl = template.Must(template.New("dialog").Parse(dialogHTML))

type Dialog struct {
	log *slog.Logger
}

func New(logger *slog.Logger) *Dialog {
	if logger == nil {
		logger = slog.Default()
	}
	return &Dialog{
		log: logger,
	}
}

const width = "_width_class"

func (d *Dialog) Render(r *http.Request, attrs ui.Attributes) (*html.Node, error) {
	switch attrs.Get("width") {
	case "fit":
		attrs.Set(width, "w-fit")
	case "full":
		attrs.Set(width, "w-full")
	case "1/2":
		attrs.Set(width, "w-1/2")
	case "1/3":
		attrs.Set(width, "w-1/3")
	default:
		attrs.Set(width, "w-72")
	}

	b := &bytes.Buffer{}
	if err := dialogTmpl.Execute(b, attrs); err != nil {
		return nil, fmt.Errorf("execute dialog template: %w", err)
	}

	return ui.ParseComponent(b)
}
