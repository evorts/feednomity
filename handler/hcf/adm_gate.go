package hcf

import (
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"net/http"
)

func AdminGate(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)

	log.Log("admin_gate_handler", "request received")

	if !req.IsLoggedIn() {
		http.Redirect(w, r, "/adm/login", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/adm/dashboard", http.StatusTemporaryRedirect)
	return
}

