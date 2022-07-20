package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"chujungeng/camera-roll/pkg/cameraroll"
)

// AlbumImageRouter specifies all the routes related to albumImages
func (handler Handler) AlbumImageRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", handler.AddImageToAlbum) // POST /albumImages

	return r
}

// AlbumImageRequest is the request body of albumImage's POST method
type AlbumImageRequest struct {
	*cameraroll.AlbumImage
}

// Bind preprocesses the request for some basic error checking
func (req *AlbumImageRequest) Bind(r *http.Request) error {
	// do nothing

	return nil
}

// AlbumImageResponse is the response body of albumImage's CRUD methods
type AlbumImageResponse struct {
	*cameraroll.AlbumImage
}

// Render preprocess the response before it's sent to the wire
func (rsp *AlbumImageResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// do nothing
	return nil
}

// NewAlbumImageResponse is the constructor method for AlbumImageResponse
func NewAlbumImageResponse(albumImage *cameraroll.AlbumImage) *AlbumImageResponse {
	rsp := AlbumImageResponse{
		AlbumImage: albumImage,
	}

	return &rsp
}

// AddImageToAlbum adds an image to the album
func (handler Handler) AddImageToAlbum(w http.ResponseWriter, r *http.Request) {
	albumImageReq := AlbumImageRequest{}

	// unmarshal new album from request
	if err := render.Bind(r, &albumImageReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// add the new relationship to database
	if err := handler.Service.AddImageToAlbum(r.Context(), albumImageReq.AlbumID, albumImageReq.ImageID); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// render response
	render.Status(r, http.StatusOK)
	render.Render(w, r, NewAlbumImageResponse(albumImageReq.AlbumImage))
}
