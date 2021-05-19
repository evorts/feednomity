package middleware

import (
	"context"
	"github.com/evorts/feednomity/pkg/api"
	"github.com/evorts/feednomity/pkg/view"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
)

func WithWebMethodFilter(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.ToUpper(r.Method) != method {
			render(http.StatusMethodNotAllowed, "not allowed", w, forbiddenTemplate, nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func evalFilters(r *http.Request, method string, allowedMethods, allowedOrigins []string) (status int, err error, origin string) {
	var (
		err2 *api.ResponseError
		uri *url.URL
	)
	status, err2, uri = evalCors(r, allowedMethods, allowedOrigins)
	origin = uri.String()
	if err2 != nil {
		return status, errors.New(err2.Message), origin
	}
	if r.Method == http.MethodOptions {
		return http.StatusOK, nil, origin
	}
	if strings.ToUpper(r.Method) != method {
		return http.StatusMethodNotAllowed, errors.New("not allowed"), origin
	}
	return http.StatusContinue, nil, origin
}

func WithFiltersForApi(method string, allowedMethods, allowedOrigins []string, vm view.IManager, next http.Handler) http.Handler {
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
		next.ServeHTTP(w, r)
	})
}

func WithInjection(next http.Handler, contextInjection map[string]interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if len(contextInjection) > 0 {
			for k, v := range contextInjection {
				ctx = context.WithValue(ctx, k, v)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
