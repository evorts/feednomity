package hcf

import (
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/logger"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/utils"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	req := reqio.NewRequest(w, r).Prepare()
	log := req.GetContext().Get("logger").(logger.IManager)
	cfg := req.GetContext().Get("cfg").(config.IManager)

	log.Log("logout_handler", "request received")

	redirectUrl := utils.IIf(len(req.GetQueryParam("ref")) > 0, req.GetQueryParam("ref"), "/mbr/login")

	// remove cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "feednomisess",
		Domain:   cfg.GetConfig().App.CookieDomain,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})

	sm := req.GetContext().Get("sm").(session.IManager)

	if err := sm.Destroy(r.Context()); err != nil {
		_ = sm.RenewToken(r.Context())
	}

	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}
