package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"google.golang.org/api/idtoken"
)

const (
	JWTClaimUserRole      string = "user_role"
	JWTClaimUserRoleAdmin string = "admin"
	JWTClaimEmail         string = "email"
)

const (
	AdminAccessTimeLimit = 24 * time.Hour
)

// TokenResponse is the response body of oauth login requests
type TokenResponse struct {
	Token string `json:"token"`
}

// Render preprocess the response before it's sent to the wire
func (rsp *TokenResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// do nothing
	return nil
}

// NewAlbumResponse is the constructor method for AlbumResponse type
func NewTokenResponse(token string) *TokenResponse {
	resp := TokenResponse{Token: token}

	return &resp
}

// TokenRouter exchanges OAuth tokens for JWT token
func (handler Handler) TokenRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/google", handler.VerifyGoogleIdToken)

	return r
}

func (handler Handler) VerifyGoogleIdToken(w http.ResponseWriter, r *http.Request) {
	// get token from Authorization header
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) == 0 {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	token := splitToken[1]

	// validate token with Google
	payload, err := idtoken.Validate(context.Background(), token, handler.googleOAuthConfig.ClientID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// check if email is admin
	if email := payload.Claims[JWTClaimEmail]; email != handler.adminID {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// respond with an admin JWT on success
	jwt, _ := handler.generateAdminJWT(time.Now().Add(AdminAccessTimeLimit))
	render.Render(w, r, NewTokenResponse(jwt))
}

// generateAdminJWT creates a JWT token that has a userRole of admin
func (handler Handler) generateAdminJWT(expiresAt time.Time) (string, error) {
	claims := map[string]interface{}{
		JWTClaimUserRole: JWTClaimUserRoleAdmin,
	}

	jwtauth.SetExpiry(claims, expiresAt)
	_, tokenString, err := handler.jwtTokenAuth.Encode(claims)

	return tokenString, err
}

// GenerateTestJWT creates a JWT token for debugging purposes
func (handler Handler) GenerateTestJWT() string {
	const (
		testTokenExpires = 30 * time.Minute
	)

	tokenString, _ := handler.generateAdminJWT(time.Now().Add(testTokenExpires))

	return tokenString
}
