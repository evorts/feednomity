package handler

import (
	"github.com/evorts/godash/pkg/logger"
	"github.com/evorts/godash/pkg/reqio"
	"github.com/evorts/godash/pkg/template"
	"net/http"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	log.Log("dashboard_handler", "request received")
	//req := Request{w: w, r: r}
	// check if the request already authenticated
	/*rp := r.URL.Path
	if rp != "/" && rp != "/login" && rp != "/logout" && rp != "/reload" &&
		!strings.HasPrefix(rp, "/assets/") {
		http.Redirect(w, r, "/not-found", http.StatusPermanentRedirect)
		return
	}
	if !req.isLoggedIn() {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}*/
	// render dashboard page
	if err := view.Render(w, "dashboard.html", map[string]interface{}{
		"PageAttributes": map[string]interface{}{
			"Title": "Dashboard Page",
		},
	}); err != nil {
		log.Log("dashboard_handler", err.Error())
	}
}
