package main

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/kheepercom/ui"
	"github.com/kheepercom/ui/components/neobrutal/button"
	"github.com/kheepercom/ui/components/neobrutal/card"
	"github.com/kheepercom/ui/components/neobrutal/checkbox"
	"github.com/kheepercom/ui/components/neobrutal/dialog"
	"github.com/kheepercom/ui/components/neobrutal/dropdown"
	"github.com/kheepercom/ui/components/neobrutal/leftsidebar"
	"github.com/kheepercom/ui/components/neobrutal/menu"
)

//go:embed pages public
var appfs embed.FS

func main() {
	reg := ui.Registry{}
	logger := slog.Default()

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
		LiveReload:  true,
	})
	mux.HandleFunc("GET /um", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("um"))
	})

	root := http.NewServeMux()
	root.Handle("/neobrutal/", http.StripPrefix("/neobrutal", mux))
	http.ListenAndServe(":8888", root)
}
