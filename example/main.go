package main

import (
	"embed"
	"math/rand"
	"net/http"

	"github.com/kheepercom/ui"
	"github.com/kheepercom/ui/example/components/loginout"
)

//go:embed pages public
var appfs embed.FS

func main() {
	reg := ui.Registry{}

	reg.Add("loginout", loginout.Must())

	mux := ui.Must(reg, appfs, ui.Options{})

	http.ListenAndServe(":8888", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Intn(2) == 1 {
			ctx := loginout.WithUser(r.Context(), "j_smith")
			r = r.WithContext(ctx)
		}
		mux.ServeHTTP(w, r)
	}))
}
