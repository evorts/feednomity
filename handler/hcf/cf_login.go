package hcf

import (
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
)

func Login(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.ITemplateManager)
	sm := req.GetContext().Get("sm").(session.IManager)

	log.Log("member_login_handler", "request received")

	if req.IsLoggedIn() {
		http.Redirect(w, r, "/mbr/list", http.StatusTemporaryRedirect)
		return
	}
	renderData := map[string]interface{}{
		"PageTitle": "Login Page",
	}
	// render login page
	sm.Put(r.Context(), "csrf", req.GetCsrfToken())
	if err := vm.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "member-login.html", renderData); err != nil {
		log.Log("member_login_handler", err.Error())
	}
}
