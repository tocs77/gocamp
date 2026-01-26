package middlewares

import (
	"net/http"
	"strings"
)

func ExcludeRoutes(middlewareFunc func(http.Handler) http.Handler, excludeRoutes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, route := range excludeRoutes {
				if strings.HasPrefix(r.URL.Path, route) {
					next.ServeHTTP(w, r)
					return
				}
			}
			middlewareFunc(next).ServeHTTP(w, r)
		})
	}
}
