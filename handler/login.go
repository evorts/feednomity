package handler

import (
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	view := req.GetContext().Get("view").(template.IManager)
	sm := req.GetContext().Get("sm").(session.IManager)

	log.Log("login_handler", "request received")

	if req.IsLoggedIn() {
		http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		return
	}
	renderData := map[string]interface{}{
		"PageTitle": "Login Page",
	}
	// render login page
	sm.Put(r.Context(), "token", req.GetToken())
	if err := view.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "login.html", renderData); err != nil {
		log.Log("login_handler", err.Error())
	}
}
