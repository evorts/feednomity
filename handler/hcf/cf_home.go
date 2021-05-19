package hcf

import (
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	vm := req.GetContext().Get("view").(view.ITemplateManager)

	if !req.IsMethodGet() && req.GetPath() != "/" {
		_ = vm.Render(w, http.StatusNotFound, "404.html", map[string]interface{}{
			"PageTitle": "404 Page Not Found",
		})
		return
	}
	if req.IsLoggedIn() {
		http.Redirect(w, r, "/mbr/review/list", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, "/mbr/login", http.StatusTemporaryRedirect)
	}
}
