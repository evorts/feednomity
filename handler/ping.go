package handler

import (
	"github.com/evorts/godash/pkg/reqio"
	"github.com/evorts/godash/pkg/template"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	view := req.GetContext().Get("view").(template.IManager)
	if !req.IsMethodGet() {
		_ = view.RenderRaw(w, "NOK")
		return
	}
	_ = view.RenderRaw(w, "OK")
}

