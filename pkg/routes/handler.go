package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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

// Routes is the collection of all routes being served
func (handler Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// CORS policies
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   handler.corsOrigin,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	RootServer(r)

	// Create a route along /assets that will serve contents from
	// the ./public/ folder.
	FileServer(r, StaticAssetURL(), http.Dir(StaticFileDir()))

	r.Mount("/api/", handler.ApiRouter())
	r.Mount("/auth/", handler.AuthRouter())

	return r
}
