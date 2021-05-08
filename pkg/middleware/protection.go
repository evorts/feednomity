package middleware

import (
	"context"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strings"
)

type ProtectionLib struct {
	Jwe  jwe.IManager
	Acl  acl.IManager
	Sm   session.IManager
	View view.IManager
}

type ProtectionArgs struct {
	Path           string
	Method         string
	AllowedMethods []string
	AllowedOrigins []string
	RenderType     string // json, html
}

const forbiddenTemplate = "forbidden.html"

// WithSessionProtection when using session, e.g. web
func WithSessionProtection(sm session.IManager, vm view.ITemplateManager, acc acl.IManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userSession reqio.UserData
		ctx := r.Context()
		err := sm.GetJson(ctx, "user", &userSession)
		if err != nil {
			render(
				http.StatusForbidden, api.NewResponseError("ACC:PERM:DND", "Permission denied!", nil, nil),
				w, forbiddenTemplate, vm,
			)
			return
		}
		//check access permission
		//render template when violate permission
		allowed, accessScope := acc.IsAllowed(userSession.Id, r.Method, r.URL.Path)
		if !allowed {
			render(
				http.StatusForbidden, api.NewResponseError("ACC:PERM:DND", "Permission denied!", nil, nil),
				w, forbiddenTemplate, vm,
			)
			return
		}
		sm.Put(ctx, "access_scope", accessScope)
		vm.InjectData("user", userSession)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//WithTokenProtection when using jwe, e.g. API
func WithTokenProtection(
	method string, allowedMethods, allowedOrigins []string, acc acl.IManager, jw jwe.IManager, next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) != method {
			renderJson(
				http.StatusMethodNotAllowed,
				api.NewResponseError("FIL:MTD:NA", "Not allowed!", nil, nil),
				w, nil,
			)
			return
		}
		var (
			status int
			errRs  *api.ResponseError
			err    error
		)
		ctx := r.Context()
		if len(allowedMethods) > 0 && len(allowedOrigins) > 0 {
			status, errRs = evalCors(r, allowedMethods, allowedOrigins)
		}
		if errRs != nil {
			renderJson(status, errRs, w, nil)
			return
		}
		var (
			userData reqio.UserData
			pri      jwe.PrivateClaims
		)
		jweToken := strings.Trim(r.Header.Get("X-Authorization"), " ")
		_, pri, err = jw.Decode(jweToken)
		if err != nil {
			renderJson(
				http.StatusForbidden, api.NewResponseError("ACC:PERM:DND", "Permission denied!", nil, nil),
				w, nil,
			)
			return
		}
		userData = reqio.UserData{
			Id:          pri.Id,
			Username:    pri.Username,
			DisplayName: pri.DisplayName,
			Attributes:  pri.Attributes,
			Email:       pri.Email,
			Phone:       pri.Phone,
			AccessRole:  pri.AccessRole,
			JobRole:     pri.JobRole,
			Assignment:  pri.Assignment,
			GroupId:     pri.GroupId,
			OrgId:       pri.OrgId,
			OrgGroupIds: pri.OrgGroupIds,
		}

		//check access permission
		//render template when violate permission
		allowed, accessScope := acc.IsAllowed(userData.Id, method, r.URL.Path)
		if !allowed {
			renderJson(
				http.StatusForbidden, api.NewResponseError("ACC:PERM:DND", "Permission denied!", nil, nil),
				w, nil,
			)
			return
		}
		ctx = context.WithValue(ctx, "user", &userData)
		ctx = context.WithValue(ctx, "access_scope", accessScope)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithCorsProtection(allowedMethods []string, allowedOrigins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status, err := evalCors(r, allowedMethods, allowedOrigins)
		if err != nil {
			renderJson(status, err, w, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
