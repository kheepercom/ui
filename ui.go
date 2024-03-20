package ui

import (
	"bytes"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Options struct {
	// Render and serve all files under this directory
	PagesRoot string
	// Serve all files under this directory unmodified
	PublicRoot string

	// Logger to use for messages and errors
	Logger *slog.Logger

	// Serve a component catalog at this path
	CatalogPath string

	// Reload browser when air restarts the app.
	LiveReload bool
}

func Must(reg Registry, uifs fs.FS, opts Options) *http.ServeMux {
	mux, err := New(reg, uifs, opts)
	if err != nil {
		panic(err)
	}
	return mux
}

// New creates a ServeMux with a handler for each page and each public file.
// The fs argument must contain directories named "public" and "pages".
func New(reg Registry, uifs fs.FS, opts Options) (*http.ServeMux, error) {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}
	pubRoot := opts.PublicRoot
	if pubRoot == "" {
		pubRoot = "public"
	}
	pageRoot := opts.PagesRoot
	if pageRoot == "" {
		pageRoot = "pages"
	}

	mux := http.NewServeMux()

	err := fs.WalkDir(uifs, pubRoot, func(path string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		pattern := "GET " + strings.TrimPrefix(path, pubRoot)
		logger.Debug(pattern)
		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, uifs, path)
		})
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = fs.WalkDir(uifs, pageRoot, func(filepathname string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		pageHTML, err := fs.ReadFile(uifs, filepathname)
		if err != nil {
			return fmt.Errorf("failed to read %q", filepathname)
		}
		page, err := ParsePage(bytes.NewReader(pageHTML))
		if err != nil {
			return fmt.Errorf("failed to parse %q", filepathname)
		}

		pattern := strings.TrimPrefix(filepathname, pageRoot)
		if path.Base(pattern) == "index.html" {
			withIndexHTML := "GET " + pattern
			logger.Debug(withIndexHTML)
			mux.HandleFunc(withIndexHTML, func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, pattern, http.StatusMovedPermanently)
			})
			pattern = path.Dir(pattern)
		}
		pattern = "GET " + pattern
		logger.Info(pattern)

		mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			node, err := reg.Render(r, Clone(page))
			if err != nil {
				http.Error(w, "Internal Server Error", 500)
				logger.Error("Failed to render page with Kheeper UI", "page", filepathname, "error", err)
				return
			}
			err = html.Render(w, node)
			if err != nil {
				logger.Error("Failed to send HTML response with Kheeper UI", "page", filepathname, "error", err)
			}
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	if opts.CatalogPath != "" {
		err := catalog(reg, mux, CatalogOptions{
			Prefix: opts.CatalogPath,
			Logger: logger,
		})
		if err != nil {
			return nil, err
		}
	}

	if opts.LiveReload || opts.CatalogPath != "" {
		mux.HandleFunc("GET /reload", func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Hour * 24)
		})
	}

	return mux, nil
}
