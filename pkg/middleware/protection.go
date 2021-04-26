package middleware

import (
	"fmt"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/template"
	"github.com/evorts/feednomity/pkg/utils"
	"net/http"
	"strings"
)

type ProtectionLib struct {
	Acl  acl.IManager
	Sm   session.IManager
	View template.IManager
}

type ProtectionArgs struct {
	Path           string
	Method         string
	AllowedMethods []string
	AllowedOrigins []string
	RenderType     string // json, html
}

const forbiddenTemplate = "forbidden.html"

//WithProtection when using session or jwe
func WithProtection(lib ProtectionLib, args ProtectionArgs, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			status int
			errRs  *api.ResponseError
		)
		ctx := r.Context()
		if len(args.AllowedMethods) > 0 && len(args.AllowedOrigins) > 0 {
			status, errRs = evalCors(r, args.AllowedMethods, args.AllowedOrigins)
		}
		if errRs != nil {
			render(status, errRs, w, "", lib.View)
			return
		}
		if lib.Sm != nil && lib.Acl != nil && len(args.Path) > 0 && len(args.Method) > 0 {
			var user reqio.UserSession
			err := lib.Sm.GetJson(ctx, "user", &user)
			//check access permission
			//render template when violate permission
			allowed, accessScope := lib.Acl.IsAllowed(user.Id, args.Method, args.Path)
			if err != nil || !allowed {
				render(
					http.StatusForbidden, api.NewResponseError("ACC:PERM:DND", "Permission denied!", nil, nil),
					w, utils.IIf(args.RenderType == "json", "", forbiddenTemplate), lib.View,
				)
				return
			}
			lib.Sm.Put(ctx, "access_scope", accessScope)
			lib.View.InjectData("user", user)
		}
		if len(args.Method) > 0 && strings.ToUpper(r.Method) != args.Method {
			render(
				http.StatusMethodNotAllowed,
				api.NewResponseError("FIL:MTD:NA", "Not allowed!", nil, nil),
				w, utils.IIf(args.RenderType == "json", "", forbiddenTemplate), lib.View,
			)
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = fmt.Fprintln(w, "")
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithCors(view template.IManager, allowedMethods []string, allowedOrigins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status, err := evalCors(r, allowedMethods, allowedOrigins)
		if err != nil {
			_ = view.RenderJson(w, status, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
