package routes

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"chujungeng/camera-roll/pkg/cameraroll"
)

// ImageTagRouter specifies all the routes related to imageTag
func (handler Handler) ImageTagRouter() chi.Router {
	r := chi.NewRouter()

	r.Post("/", handler.AddTagToImage) // POST /admin/imageTags

	return r
}

// ImageTagRequest is the request body of imageTag's POST method
type ImageTagRequest struct {
	*cameraroll.ImageTag
}

// Bind preprocesses the request for some basic error checking
func (req *ImageTagRequest) Bind(r *http.Request) error {
	// do nothing

	return nil
}

// ImageTagResponse is the response body of imageTag's CRUD methods
type ImageTagResponse struct {
	*cameraroll.ImageTag
}

// Render preprocess the response before it's sent to the wire
func (rsp *ImageTagResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// do nothing
	return nil
}

// NewImageTagResponse is the constructor method for ImageTagResponse
func NewImageTagResponse(imageTag *cameraroll.ImageTag) *ImageTagResponse {
	rsp := ImageTagResponse{
		ImageTag: imageTag,
	}

	return &rsp
}

// AddTagToImage adds a tag to the image
func (handler Handler) AddTagToImage(w http.ResponseWriter, r *http.Request) {
	imageTagReq := ImageTagRequest{}

	// unmarshal new image from request
	if err := render.Bind(r, &imageTagReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// add the new relationship to database
	if err := handler.Service.AddTagToImage(r.Context(), imageTagReq.ImageID, imageTagReq.TagID); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// render response
	render.Status(r, http.StatusOK)
	render.Render(w, r, NewImageTagResponse(imageTagReq.ImageTag))
}
