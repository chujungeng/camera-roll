package routes

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
)

// AdminOnly verifies if user_role is admin from the JWT claims
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, claims, err := FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if userRole := claims[JWTClaimUserRole]; userRole != JWTClaimUserRoleAdmin {
			http.Error(w, "authentication failed", http.StatusUnauthorized)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

func FromContext(ctx context.Context) (jwt.Token, map[string]interface{}, error) {
	token, _ := ctx.Value(jwtauth.TokenCtxKey).(jwt.Token)

	var err error
	var claims map[string]interface{}

	if token != nil {
		claims, err = token.AsMap(context.Background())
		if err != nil {
			return token, nil, err
		}
	} else {
		claims = map[string]interface{}{}
	}

	err, _ = ctx.Value(jwtauth.ErrorCtxKey).(error)

	return token, claims, err
}
