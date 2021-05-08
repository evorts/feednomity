package hadm

import (
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	if !req.IsLoggedIn() {
		http.Redirect(w, r, "/adm/login", http.StatusTemporaryRedirect)
		return
	}
	sm := req.GetContext().Get("sm").(session.IManager)
	if err := sm.Destroy(r.Context()); err != nil {
		_ = sm.RenewToken(r.Context())
	}
	http.Redirect(w, r, "/adm/login", http.StatusTemporaryRedirect)
}
