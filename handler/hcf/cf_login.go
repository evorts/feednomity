package hcf

import (
	"fmt"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/utils"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	vm := req.GetContext().Get("view").(view.ITemplateManager)
	sm := req.GetContext().Get("sm").(session.IManager)

	isAdmin := strings.Contains(req.GetPath(), "/adm/")
	redirectUrl := utils.IIf(isAdmin, "/adm/dashboard", "/mbr/review/list")
	prefix := utils.IIf(isAdmin, "admin", "member")

	log.Log(fmt.Sprintf("%s_login_handler", prefix), "request received")

	if req.IsLoggedIn() {
		http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
		return
	}

	ref := req.GetQueryParam("ref")
	renderData := map[string]interface{}{
		"PageTitle":   fmt.Sprintf("%s Login Page", strings.Title(prefix)),
		"RedirectUrl": utils.IIf(len(ref) > 0, ref, redirectUrl),
		"UserID": req.GetQueryParam("user"),
	}

	// render login page
	sm.Put(r.Context(), "csrf", req.GetCsrfToken())
	if err := vm.InjectData("Csrf", req.GetToken()).Render(w, http.StatusOK, "member-login.html", renderData); err != nil {
		log.Log(fmt.Sprintf("%s_login_handler", prefix), err.Error())
	}
}
