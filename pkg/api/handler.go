package api

import (
	"time"

	"github.com/go-chi/jwtauth/v5"

	"chujungeng/camera-roll/pkg/cameraroll"
)

const (
	JWTClaimUserRole      string = "user_role"
	JWTClaimUserRoleAdmin string = "admin"
)

// Handler handles all API requests to camera roll
type Handler struct {
	Service cameraroll.Service

	jwtTokenAuth *jwtauth.JWTAuth
}

// NewHandler is the contructor method for the Handler
func NewHandler(service cameraroll.Service, jwtSecret string) *Handler {
	handler := Handler{
		Service:      service,
		jwtTokenAuth: jwtauth.New("HS256", []byte(jwtSecret), nil),
	}

	return &handler
}

// generateAdminJWT creates a JWT token whith admin claim
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
