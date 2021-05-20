package middleware

import (
	"context"
	"fmt"
	"github.com/evorts/feednomity/pkg/acl"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/config"
	"github.com/evorts/feednomity/pkg/jwe"
	"github.com/evorts/feednomity/pkg/reqio"
	"github.com/evorts/feednomity/pkg/session"
	"github.com/evorts/feednomity/pkg/view"
	"net/http"
	"strings"
	"time"
)

const (
	forbiddenTemplate = "forbidden.html"
	cookieToken       = "feednomisess"
)

// WithSessionProtection when using session, e.g. web
func WithSessionProtection(sm session.IManager, vm view.ITemplateManager, acc acl.IManager, jw jwe.IManager, cfg config.IManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userData reqio.UserData
			jweToken string
			err error
		)
		ctx := r.Context()
		err = sm.GetJson(ctx, "user", &userData)
		if err != nil {
			fmt.Println(err)
		}
		if userData.Id < 1 {
			//try to get from cookie
			cookie, err2 := r.Cookie(cookieToken)
			if err2 == nil {
				jweToken = cookie.Value
			}
			if len(jweToken) > 0 {
				status, code, message, _, pri := parseToken(jw, jweToken)
				if len(message) > 0 {
					// remove cookie
					http.SetCookie(w, &http.Cookie{
						Domain:  cfg.GetConfig().App.CookieDomain,
						Name:    cookieToken,
						Value:   "",
						Path:    "/",
						Expires: time.Unix(0, 0),
					})
					render(
						status, api.NewResponseError(code, message, nil, nil),
						w, forbiddenTemplate, vm,
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
				err = sm.PutJson(ctx, "user", userData)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		if userData.Id < 1 {
			render(
				http.StatusForbidden, api.NewResponseError("ACC:PERM:DND", "Permission denied!", nil, nil),
				w, forbiddenTemplate, vm,
			)
			return
		}
		//check access permission
		//render template when violate permission
		allowed, accessScope := acc.IsAllowed(userData.Id, r.Method, r.URL.Path)
		if !allowed {
			render(
				http.StatusForbidden, api.NewResponseError("ACC:PERM:DND", "Permission denied!", nil, nil),
				w, forbiddenTemplate, vm,
			)
			return
		}
		sm.Put(ctx, "access_scope", accessScope)
		vm.InjectData("user", userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//WithTokenProtection when using jwe, e.g. API
func WithTokenProtection(
	method string, allowedMethods, allowedOrigins []string, acc acl.IManager, jw jwe.IManager, vm view.IManager, next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status, err, origin := evalFilters(r, method, allowedMethods, allowedOrigins)
		vm.ResetHeaders()
		vm.InjectHeader("Access-Control-Allow-Methods", r.Method)
		vm.InjectHeader("Access-Control-Allow-Origin", origin)
		if err != nil {
			_ = vm.RenderJson(w, status, map[string]interface{}{"error": err.Error()})
			return
		}
		if status != http.StatusContinue {
			_ = vm.RenderJson(w, status, make([]string, 0))
			return
		}
		ctx := r.Context()
		userData := reqio.UserData{}
		jweToken := strings.Trim(r.Header.Get("X-Authorization"), " ")
		status, code, message, _, pri := parseToken(jw, jweToken)
		if len(message) > 0 {
			_ = vm.RenderJson(w, status, api.NewResponseError(code, message, nil, nil))
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
			_ = vm.RenderJson(w, http.StatusForbidden, api.NewResponseError("ACC:PERM:DND", "Permission denied!", nil, nil))
			return
		}
		ctx = context.WithValue(ctx, "user", userData)
		ctx = context.WithValue(ctx, "access_scope", accessScope)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
