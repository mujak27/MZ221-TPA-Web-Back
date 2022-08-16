package auth

import (
	"context"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if auth == "" {
			next.ServeHTTP(w, r)
			return
		}

		authToken := auth[len("Bearer "):]
		if authToken == "" {
			http.Error(w, "token not exists", http.StatusForbidden)
			return
		}

		validate, err := JwtValidate(context.Background(), authToken)
		if err != nil || !validate.Valid {
			http.Error(w, "invalid jwt token", http.StatusForbidden)
			return
		}

		customClaim, _ := validate.Claims.(*JwtCustomClaim)

		ctx := context.WithValue(r.Context(), string("tokenValue"), customClaim)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
