package main

import (
	_ "embed"
	"math/rand"
	"net/http"
	"strings"

	"github.com/kheepercom/ui"
	"golang.org/x/net/html"
)

//go:embed home.html
var homeHTML string

func main() {
	reg := ui.Registry{}
	l, err := NewLoginout()
	if err != nil {
		panic(err)
	}
	reg.Add("loginout", l)

	home, err := ui.ParsePage(strings.NewReader(homeHTML))
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", 404)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if rand.Intn(2) == 1 {
			ctx := WithUser(r.Context(), "j_smith")
			r = r.WithContext(ctx)
		}
		h, err := reg.Render(r, ui.Clone(home))
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}
		html.Render(w, h)
	})

	http.ListenAndServe(":8888", nil)
}
