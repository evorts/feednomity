package middleware

import (
	"encoding/json"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/template"
	"net/http"
	"net/url"
	"strings"
)

func evalCors(r *http.Request, methods, origins []string) (status int, responseError *api.ResponseError) {
	allowed := false
	for _, m := range methods {
		if strings.ToUpper(m) == strings.ToUpper(r.Method) {
			allowed = true
			break
		}
	}
	if !allowed {
		return http.StatusMethodNotAllowed, api.NewResponseError("COR:MTD:NAL", "", nil, nil)
	}
	ref := r.Referer()
	if len(ref) < 1 {
		return http.StatusNotAcceptable, api.NewResponseError("COR:REF:NAC", "", nil, nil)
	}
	uri, err := url.Parse(ref)
	if err != nil {
		return http.StatusExpectationFailed, api.NewResponseError("COR:URI:EXF", "", nil, nil)
	}
	allowed = false
	for _, o := range origins {
		if uri.Host == o {
			allowed = true
			break
		}
	}
	if !allowed {
		return http.StatusExpectationFailed, api.NewResponseError("COR:ALW:FBD", "", nil, nil)
	}
	return http.StatusAccepted, nil
}

func render(status int, rs interface{}, w http.ResponseWriter, tpl string, view template.IManager)  {
	switch true {
	case view != nil:
		if len(tpl) < 1 {
			_ = view.RenderJson(w, status, rs)
			return
		}
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
