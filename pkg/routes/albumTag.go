package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"chujungeng/camera-roll/pkg/cameraroll"
)

// AlbumTagRouter specifies all the routes related to albumTag
func (handler Handler) AlbumTagRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", handler.AddTagToAlbum) // POST /admin/albumTags

	return r
}

// AlbumTagRequest is the request body of albumTag's POST method
type AlbumTagRequest struct {
	*cameraroll.AlbumTag
}

// Bind preprocesses the request for some basic error checking
func (req *AlbumTagRequest) Bind(r *http.Request) error {
	// do nothing

	return nil
}

// AlbumTagResponse is the response body of albumTag's CRUD methods
type AlbumTagResponse struct {
	*cameraroll.AlbumTag
}

// Render preprocess the response before it's sent to the wire
func (rsp *AlbumTagResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// do nothing
	return nil
}

// NewAlbumTagResponse is the constructor method for AlbumTagResponse
func NewAlbumTagResponse(albumTag *cameraroll.AlbumTag) *AlbumTagResponse {
	rsp := AlbumTagResponse{
		AlbumTag: albumTag,
	}

	return &rsp
}

// AddTagToAlbum adds a tag to the album
func (handler Handler) AddTagToAlbum(w http.ResponseWriter, r *http.Request) {
	albumTagReq := AlbumTagRequest{}

	// unmarshal new album from request
	if err := render.Bind(r, &albumTagReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// add the new relationship to database
	if err := handler.Service.AddTagToAlbum(r.Context(), albumTagReq.AlbumID, albumTagReq.TagID); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// render response
	render.Status(r, http.StatusOK)
	render.Render(w, r, NewAlbumTagResponse(albumTagReq.AlbumTag))
}
