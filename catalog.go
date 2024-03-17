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

const catalogPage = `
<!doctype html>
<head>
	<meta charset="UTF-8">
	<title>{{ .Title }}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<script src="https://cdn.tailwindcss.com/3.4.1"></script>
</head>
<body>
	{{ .Example }}
</body>
`

var catalogTmpl = template.Must(template.New("catalog").Parse(catalogPage))

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

	for name, component := range reg {
		c, ok := component.(Examples)
		if !ok {
			continue
		}

		examplesFS := c.Examples()

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

		return err
	}

	return nil
}
