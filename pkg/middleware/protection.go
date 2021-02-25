package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func WithCors(allowedMethods []string, allowedOrigins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowed := false
		for _, m := range allowedMethods {
			if strings.ToUpper(m) == strings.ToUpper(r.Method) {
				allowed = true
				break
			}
		}
		if !allowed {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = fmt.Fprintln(w, "")
			return
		}
		ref := r.Referer()
		if len(ref) < 1 {
			w.WriteHeader(http.StatusNotAcceptable)
			_, _ = fmt.Fprintln(w, "")
			return
		}
		uri, err := url.Parse(ref)
		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			_, _ = fmt.Fprintln(w, "")
			return
		}
		allowed = false
		for _, o := range allowedOrigins {
			if uri.Host == o {
				allowed = true
				break
			}
		}
		if !allowed {
			w.WriteHeader(http.StatusForbidden)
			_, _ = fmt.Fprintln(w, "")
			return
		}
		next.ServeHTTP(w, r)
	})
}

//when using session or jwe
func WithProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithAccessControl(path, method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

