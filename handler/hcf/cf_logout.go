package hcf

import (
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)

	log.Log("member_logout_handler", "request received")

	// remove cookie
	http.SetCookie(w, &http.Cookie{
		Name:       "feednomisess",
		Value:      "",
		Path:       "/",
		Expires:    time.Unix(0, 0),
	})

	sm := req.GetContext().Get("sm").(session.IManager)

	if err := sm.Destroy(r.Context()); err != nil {
		_ = sm.RenewToken(r.Context())
	}

	http.Redirect(w, r, "/mbr/login", http.StatusTemporaryRedirect)
}
