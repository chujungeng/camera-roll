package api

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

// FsWithoutDirListing is Go's filesystem with directory listing turned off
type FsWithoutDirListing struct {
	http.Dir
}

func (m FsWithoutDirListing) Open(name string) (result http.File, err error) {
	f, err := m.Dir.Open(name)
	if err != nil {
		return
	}

	fi, err := f.Stat()
	if err != nil {
		return
	}
	if fi.IsDir() {
		// Return a response that would have been if directory would not exist:
		return m.Dir.Open("does-not-exist")
	}

	return f, nil
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.Dir) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(FsWithoutDirListing{root}))
		fs.ServeHTTP(w, r)
	})
}
