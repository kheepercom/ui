package main

import (
	"embed"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/kheepercom/ui"
	"github.com/kheepercom/ui/components/neobrutal/button"
	"github.com/kheepercom/ui/components/neobrutal/card"
	"github.com/kheepercom/ui/components/neobrutal/checkbox"
	"github.com/kheepercom/ui/components/neobrutal/dialog"
	"github.com/kheepercom/ui/components/neobrutal/dropdown"
	"github.com/kheepercom/ui/components/neobrutal/leftsidebar"
	"github.com/kheepercom/ui/components/neobrutal/menu"
	"github.com/kheepercom/ui/example/components/loginout"
)

//go:embed pages public
var appfs embed.FS

func main() {
	reg := ui.Registry{}
	logger := slog.Default()

	reg.Add("loginout", loginout.Must())
	reg.Add("NeoBrutalButton", &button.Button{})
	reg.Add("NeoBrutalCard", card.New(logger))
	reg.Add("NeoBrutalCheckbox", &checkbox.Checkbox{})
	reg.Add("NeoBrutalDialog", dialog.New(logger))
	reg.Add("NeoBrutalDropdown", dropdown.New(logger))
	reg.Add("NeoBrutalDropdownItem", &dropdown.Item{})
	reg.Add("NeoBrutalMenuGroup", menu.NewGroup(logger))
	reg.Add("NeoBrutalMenuItem", menu.NewItem(logger))
	reg.Add("NeoBrutalLeftSidebar", leftsidebar.New(logger))

	mux := ui.Must(reg, appfs, ui.Options{
		CatalogPath: "/catalog",
	})

	http.ListenAndServe(":8888", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Intn(2) == 1 {
			ctx := loginout.WithUser(r.Context(), "j_smith")
			r = r.WithContext(ctx)
		}
		mux.ServeHTTP(w, r)
	}))
}
