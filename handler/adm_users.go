package handler

import (
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/template"
	"net/http"
)

func Users(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)

	log.Log("users_handler", "request received")

	if !req.IsLoggedIn() {
		//http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//return
	}

	// render dashboard page
	if err := view.Render(w, http.StatusOK, "admin-users.html", map[string]interface{}{
		"PageTitle": "Admin Dashboard Page",
	}); err != nil {
		log.Log("dashboard_handler", err.Error())
	}
}
