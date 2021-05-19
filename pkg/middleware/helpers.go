package middleware

import (
	"encoding/json"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/view"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func parseToken(jw jwe.IManager, value string) (status int, code, message string, pub jwt.Claims, pri jwe.PrivateClaims){
	pub, pri, err := jw.Decode(value)
	if err != nil {
		return http.StatusForbidden, "ACC:PERM:DND", "Permission denied!", pub, pri
	}
	if pub.Expiry.Time().Before(time.Now()) {
		return http.StatusBadRequest, "ACC:PERM:EXP", "token has expired!", pub, pri
	}
	return http.StatusContinue, "", "", pub, pri
}

func getReferrer(r *http.Request) string {
	origin := r.Header.Get("X-HTTP-FORWARD-FOR")
	if len(origin) < 1 {
		origin = r.Header.Get("Origin")
	}
	if len(origin) < 1 {
		origin = r.Referer()
	}
	if len(origin) < 1 {
		origin = r.RemoteAddr
	}
	return origin
}

func evalCors(r *http.Request, methods, origins []string) (status int, responseError *api.ResponseError, u *url.URL) {
	allowed := false
	for _, m := range methods {
		if strings.ToUpper(m) == strings.ToUpper(r.Method) {
			allowed = true
			break
		}
	}
	if !allowed {
		return http.StatusMethodNotAllowed, api.NewResponseError("COR:MTD:NAL", "", nil, nil), nil
	}
	ref := getReferrer(r)
	if len(ref) < 1 {
		return http.StatusNotAcceptable, api.NewResponseError("COR:REF:NAC", "", nil, nil), nil
	}
	uri, err := url.Parse(ref)
	if err != nil {
		return http.StatusExpectationFailed, api.NewResponseError("COR:URI:EXF", "", nil, nil), uri
	}
	allowed = false
	for _, o := range origins {
		if uri.Host == o {
			allowed = true
			break
		}
	}
	if !allowed {
		return http.StatusExpectationFailed, api.NewResponseError("COR:ALW:FBD", "", nil, nil), uri
	}
	return http.StatusAccepted, nil, uri
}

func render(status int, rs interface{}, w http.ResponseWriter, tpl string, view view.ITemplateManager)  {
	switch true {
	case view != nil:
		_ = view.RenderFlex(w, status, tpl, rs)
	default:
		w.WriteHeader(status)
		j, err := json.Marshal(rs)
		if err != nil {
			_, _ = w.Write([]byte("{}"))
			return
		}
		_, _ = w.Write(j)
	}
}

func renderJson(status int, rs interface{}, w http.ResponseWriter, view view.IManager)  {
	switch true {
	case view != nil:
		_ = view.RenderJson(w, status, rs)
	default:
		w.WriteHeader(status)
		j, err := json.Marshal(rs)
		if err != nil {
			_, _ = w.Write([]byte("{}"))
			return
		}
		_, _ = w.Write(j)
	}
}
