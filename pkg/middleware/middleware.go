package middleware

import (
	"context"
	"net/http"
)

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
