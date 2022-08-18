package api

import (
	"github.com/go-chi/jwtauth/v5"
	"golang.org/x/oauth2"

	"chujungeng/camera-roll/pkg/cameraroll"
)

// Handler handles all API requests to camera roll
type Handler struct {
	Service cameraroll.Service

	rootURL           string
	corsOrigin        []string
	jwtTokenAuth      *jwtauth.JWTAuth
	adminID           string
	googleOAuthConfig *oauth2.Config
}

// NewHandler is the contructor method for the Handler
func NewHandler(service cameraroll.Service, rootURL string, corsOrigin []string, jwtSecret string, admin string, oauthGoogleConfig *oauth2.Config) *Handler {
	handler := Handler{
		Service:           service,
		rootURL:           rootURL,
		corsOrigin:        corsOrigin,
		jwtTokenAuth:      jwtauth.New("HS256", []byte(jwtSecret), nil),
		adminID:           admin,
		googleOAuthConfig: oauthGoogleConfig,
	}

	return &handler
}
