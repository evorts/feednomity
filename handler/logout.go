package handler

import (
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	sm := req.GetContext().Get("sm").(session.IManager)
	if !req.IsLoggedIn() {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	if err := sm.Destroy(r.Context()); err != nil {
		_ = sm.RenewToken(r.Context())
	}
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
