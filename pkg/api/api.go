package api

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

const (
	staticFileFolder = "public"
)

// context keys
const (
	albumKey key = iota
	tagKey
	imageKey
)

func staticFilePath() string {
	workDir, _ := os.Getwd()
	fileDirPath := filepath.Join(workDir, staticFileFolder)
	if _, err := os.Stat(fileDirPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fileDirPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	return fileDirPath
}

// ApiRouter handles RESTful API requests at /api
func (handler Handler) ApiRouter() chi.Router {
	r := chi.NewRouter()

	r.Mount("/albums", handler.AlbumRouter())
	r.Mount("/albumImages", handler.AlbumImageRouter())
	r.Mount("/albumTags", handler.AlbumTagRouter())
	r.Mount("/tags", handler.TagRouter())
	r.Mount("/images", handler.ImageRouter())
	r.Mount("/imageTags", handler.ImageTagRouter())

	return r
}

// Routes is the collection of all routes being served
func (handler Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Mount("/api", handler.ApiRouter())

	// Create a route along /assets that will serve contents from
	// the ./public/ folder.
	FileServer(r, "/assets", "public")

	return r
}
