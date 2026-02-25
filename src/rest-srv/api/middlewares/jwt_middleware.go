package middlewares

import (
	"context"
	"net/http"
	"rest-srv/utility"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenClaims, err := utility.VerifyToken(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utility.ContextKey("role"), tokenClaims["role"])
		ctx = context.WithValue(r.Context(), utility.ContextKey("exp"), tokenClaims["exp"])
		ctx = context.WithValue(r.Context(), utility.ContextKey("username"), tokenClaims["user"])
		ctx = context.WithValue(r.Context(), utility.ContextKey("userId"), tokenClaims["uid"])
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
