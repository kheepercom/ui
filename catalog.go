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
	// Prefix is the path where the UI server is mounted if other than the root
	// path.
	Prefix string
	// Path is the path where the catalog is served. Do not include the Prefix.
	// Default is "/catalog".
	Path   string
	Logger *slog.Logger
}

// catalog page is the template used to render a single component example
type catalogPageData struct {
	Title   string
	Example string
	Prefix  string
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
	<script src="{{ .Prefix }}/js/components.js"></script>
</body>`

// menu page is the template used to render links to multiple components or examples
type link struct {
	Href string
	Text string
}

type menuPageData struct {
	Title string
	links []link
}

const menuPage = `
<!doctype html>
<head>
	<meta charset="UTF-8">
	<title>{{ .Title }}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<script>fetch("/reload").catch(() => setTimeout(() => { location.reload(); }, 1000))</script>
</head>
<body>
	{{ range .links }}
	<p><a href={{ .Href }}>{{ .Text }}</a></p>
	{{ end }}
</body>`

var (
	catalogTmpl = template.Must(template.New("catalog").Parse(catalogPage))
	menuTmpl    = template.Must(template.New("menu").Parse(menuPage))
)

// TODO return a mux
func catalog(reg Registry, mux *http.ServeMux, opts CatalogOptions) error {
	log := opts.Logger
	if log == nil {
		log = slog.Default()
	}

	if opts.Path == "" {
		opts.Path = "/catalog"
	}
	base := path.Join(opts.Prefix, opts.Path)

	var componentlinks []link

	// Iterate over all components in the registry to see if they have any examples.
	// Components provide examples by implementing the Examples method returning an
	// embed.FS.
	for name, component := range reg {
		c, ok := component.(Examples)
		if !ok {
			continue
		}

		examplesFS := c.Examples()

		// Links to every component example page. Used to construct a menu page.
		var exampleLinks []link

		// Add a handler for every example file the component provides.
		err := fs.WalkDir(examplesFS, "examples", func(filename string, f fs.DirEntry, err error) error {
			if f == nil {
				// If you do `var examples embed.FS` with no embed directive you'll get this error.
				return fmt.Errorf("got nil file walking examples FS for component %s, possibily caused by an empty embed.FS")
			}
			if f.IsDir() {
				return nil
			}
			if filepath.Ext(filename) != ".html" {
				return nil
			}

			// used as the link text on the menu page
			exampleName := strings.TrimSuffix(filepath.Base(filename), ".html")

			pattern := path.Join(base, name, exampleName)

			exampleHTML, err := examplesFS.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("catalog failed to read %s: %w", filename, err)
			}

			data := &catalogPageData{
				Title:   fmt.Sprintf("%s: %s", name, exampleName),
				Example: string(exampleHTML),
			}
			fullExamplePage := &bytes.Buffer{}
			if err := catalogTmpl.Execute(fullExamplePage, data); err != nil {
				return fmt.Errorf("catalog failed to render example page for %s %s: %w", name, exampleName, err)
			}
			parsedExamplePage, err := ParsePage(fullExamplePage)
			if err != nil {
				return fmt.Errorf("catalog failed to parse example page for %s %s: %w", name, exampleName, err)
			}

			exampleLinks = append(exampleLinks, link{
				Href: pattern,
				Text: exampleName,
			})

			mux.HandleFunc("GET "+pattern, func(w http.ResponseWriter, r *http.Request) {
				c, err := reg.Render(r, Clone(parsedExamplePage))
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

		// component menu page with links to all component examples
		exampleMenuPath := path.Join(base, name, "/")
		mux.HandleFunc("GET "+exampleMenuPath, func(w http.ResponseWriter, r *http.Request) {
			data := &menuPageData{
				Title: name + " Examples",
				links: exampleLinks,
			}
			err := menuTmpl.Execute(w, data)
			if err != nil {
				log.Error(err.Error())
			}
		})

		componentlinks = append(componentlinks, link{
			Href: path.Join(base, name),
			Text: name,
		})
	}

	// catalog menu page with links to all components
	mux.HandleFunc("GET "+base+"/", func(w http.ResponseWriter, r *http.Request) {
		data := menuPageData{
			Title: "Component Catalog",
			links: componentlinks,
		}
		err := menuTmpl.Execute(w, data)
		if err != nil {
			log.Error(err.Error())
		}
	})

	return nil
}
