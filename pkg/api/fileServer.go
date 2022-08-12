package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"

	"chujungeng/camera-roll/pkg/url"
)

const (
	staticFileFolder  = "public"
	deletedFileFolder = "deleted"
	staticFileURL     = "/assets"
)

func fileDirectoryPath(dirname string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	fileDirPath := filepath.Join(exPath, dirname)
	if _, err := os.Stat(fileDirPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fileDirPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	return fileDirPath
}

func StaticAssetURL() string {
	return staticFileURL
}

func StaticFileDir() string {
	return fileDirectoryPath(staticFileFolder)
}

func DeletedFileDir() string {
	return fileDirectoryPath(deletedFileFolder)
}

func DeleteAssetFromFilesystem(assetURL string) {
	relPath := url.GetPathFromURL(assetURL)
	if len(relPath) == 0 {
		return
	}
	file := strings.Replace(relPath, StaticAssetURL(), "", 1)

	absPath := filepath.Join(StaticFileDir(), file)
	newPath := filepath.Join(DeletedFileDir(), file)

	if err := os.Rename(absPath, newPath); err != nil {
		log.Println(err)
		return
	}
}

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
