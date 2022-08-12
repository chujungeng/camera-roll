package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// context keys
const (
	albumKey key = iota
	tagKey
	imageKey
	pageIDKey
)

// ApiRouterProtected contains secured routes that require admin access
func (handler Handler) ApiRouterProtected() chi.Router {
	r := chi.NewRouter()

	// sub-routes
	r.Mount("/albums", handler.AlbumRouterProtected())
	r.Mount("/albumImages", handler.AlbumImageRouter())
	r.Mount("/albumTags", handler.AlbumTagRouter())
	r.Mount("/tags", handler.TagRouterProtected())
	r.Mount("/images", handler.ImageRouterProtected())
	r.Mount("/imageTags", handler.ImageTagRouter())

	return r
}

// ApiRouter handles RESTful API requests
func (handler Handler) ApiRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(render.SetContentType(render.ContentTypeJSON))

	// public routes
	r.Group(func(r chi.Router) {
		r.Mount("/albums", handler.AlbumRouterPublic())
		r.Mount("/tags", handler.TagRouterPublic())
		r.Mount("/images", handler.ImageRouterPublic())
		r.Mount("/auth", handler.AuthRouter())
	})

	// protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(handler.jwtTokenAuth))

		// Handle valid / invalid tokens
		r.Use(AdminOnly)

		r.Mount("/admin", handler.ApiRouterProtected())
	})

	return r
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

	r.Mount("/", handler.ApiRouter())

	// Create a route along /assets that will serve contents from
	// the ./public/ folder.
	FileServer(r, staticFileURL, http.Dir(StaticFileDir()))

	return r
}
