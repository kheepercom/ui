package card

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

//go:embed card.html
var cardHTML string
var cardTmpl = template.Must(template.New("card").Parse(cardHTML))

type Card struct {
	log *slog.Logger
}

func New(logger *slog.Logger) *Card {
	if logger == nil {
		logger = slog.Default()
	}
	return &Card{
		log: logger,
	}
}

func (c *Card) Render(r *http.Request, attrs ui.Attributes) (*html.Node, error) {
	b := &bytes.Buffer{}
	if err := cardTmpl.Execute(b, attrs); err != nil {
		return nil, fmt.Errorf("execute card template: %w", err)
	}
	return ui.ParseComponent(b)
}
