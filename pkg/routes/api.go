package routes

import (
	"github.com/go-chi/chi"
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
	r.Mount("/verify", handler.AdminRouter())

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
		r.Mount("/token", handler.TokenRouter())
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
