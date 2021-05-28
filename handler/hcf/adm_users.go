package hcf

import (
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func Users(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.ITemplateManager)

	log.Log("users_handler", "request received")

	if !req.IsLoggedIn() {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// render dashboard page
	if err := vm.Render(w, http.StatusOK, "admin-users.html", map[string]interface{}{
		"PageTitle": "Admin Dashboard Page",
	}); err != nil {
		log.Log("dashboard_handler", err.Error())
	}
}
