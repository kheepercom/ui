package ui

import (
	"bytes"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/net/html"
)

type CatalogOptions struct {
	// default is "/catalog"
	Prefix string
	Logger *slog.Logger
}

type catalogPageData struct {
	Title   string
	Example string
}

type Link struct {
	Href string
	Text string
}

type menuPageData struct {
	Title string
	Links []Link
}

const catalogPage = `
<!doctype html>
<head>
	<meta charset="UTF-8">
	<title>{{ .Title }}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body>
	{{ .Example }}
	<script src="https://cdn.tailwindcss.com/3.4.1"></script>
	<script>fetch("/reload").catch(() => setTimeout(() => { location.reload(); }, 1000))</script>
	<script src="/js/components.js"></script>
</body>`

const menuPage = `
<!doctype html>
<head>
	<meta charset="UTF-8">
	<title>{{ .Title }}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<script>fetch("/reload").catch(() => setTimeout(() => { location.reload(); }, 1000))</script>
</head>
<body>
	{{ range .Links }}
	<p><a href={{ .Href }}>{{ .Text }}</a></p>
	{{ end }}
</body>`

var (
	catalogTmpl = template.Must(template.New("catalog").Parse(catalogPage))
	menuTmpl    = template.Must(template.New("menu").Parse(menuPage))
)

// TODO return a mux
func catalog(reg Registry, mux *http.ServeMux, opts CatalogOptions) error {
	const htmlExt = ".html"

	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}
	prefix := opts.Prefix
	if prefix == "" {
		prefix = "/catalog"
	}

	var componentLinks []Link

	for name, component := range reg {
		c, ok := component.(Examples)
		if !ok {
			continue
		}

		examplesFS := c.Examples()

		var variantLinks []Link

		err := fs.WalkDir(examplesFS, "examples", func(filename string, f fs.DirEntry, err error) error {
			if f == nil {
				return fmt.Errorf("got nil file walking examples FS for component %s, possibily caused by an empty embed.FS")
			}
			if f.IsDir() {
				return nil
			}
			if filepath.Ext(filename) != ".html" {
				return nil
			}

			variant := strings.TrimSuffix(filepath.Base(filename), ".html")

			pattern := path.Join(prefix, name, variant)

			exampleHTML, err := examplesFS.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("catalog failed to read %s: %w", filename, err)
			}

			data := &catalogPageData{
				Title:   fmt.Sprintf("%s: %s", name, variant),
				Example: string(exampleHTML),
			}
			fullCatalogPage := &bytes.Buffer{}
			if err := catalogTmpl.Execute(fullCatalogPage, data); err != nil {
				return err
			}
			parsedCatalogPage, err := ParsePage(fullCatalogPage)
			if err != nil {
				return err
			}

			variantLinks = append(variantLinks, Link{
				Href: pattern,
				Text: variant,
			})

			mux.HandleFunc("GET "+pattern, func(w http.ResponseWriter, r *http.Request) {
				c, err := reg.Render(r, Clone(parsedCatalogPage))
				if err != nil {
					http.Error(w, "Internal Server Error", 500)
					log.Error(err.Error())
					return
				}
				w.Header().Set("Content-Type", "text/html")
				if err := html.Render(w, c); err != nil {
					log.Error(err.Error())
				}
			})
			return nil
		})
		if err != nil {
			return err
		}

		data := &menuPageData{
			Title: name + " Examples",
			Links: variantLinks,
		}

		mux.HandleFunc("GET "+path.Join(prefix, name)+"/", func(w http.ResponseWriter, r *http.Request) {
			err := menuTmpl.Execute(w, data)
			if err != nil {
				log.Error(err.Error())
			}
		})

		componentLinks = append(componentLinks, Link{
			Href: path.Join(prefix, name),
			Text: name,
		})
	}

	data := menuPageData{
		Title: "Component Catalog",
		Links: componentLinks,
	}
	mux.HandleFunc("GET "+prefix+"/", func(w http.ResponseWriter, r *http.Request) {
		err := menuTmpl.Execute(w, data)
		if err != nil {
			log.Error(err.Error())
		}
	})

	return nil
}
