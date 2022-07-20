package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"chujungeng/camera-roll/pkg/cameraroll"
)

const (
	ParamImageID = "imageID"
)

// ImageRouter specifies all the routes related to images
func (handler Handler) ImageRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.GetImages) // GET /images
	r.Post("/", handler.AddImage) // POST /images

	r.Route("/{imageID}", func(r chi.Router) {
		r.Use(handler.ImageCtx)            // Load the *Image on the request context
		r.Get("/", handler.GetImage)       // GET /images/123
		r.Put("/", handler.UpdateImage)    // PUT /images/123
		r.Delete("/", handler.DeleteImage) // DELETE /images/123

		r.Delete("/tags/{tagID}", handler.RemoveTagFromImage) // DELETE /images/123/tags/789
	})

	return r
}

// ImageRequest is the request body of images' CRUD operations
type ImageRequest struct {
	*cameraroll.Image
}

// Bind preprocesses the request for some basic error checking
func (req *ImageRequest) Bind(r *http.Request) error {
	// Return an error to avoid a nil pointer dereference.
	if req.Image == nil {
		return errors.New("missing required Image fields")
	}

	return nil
}

// ImageImagesResponse is the response body of imageImages' GET method
type ImageResponse struct {
	*cameraroll.Image
}

// Render preprocess the response before it's sent to the wire
func (rsp *ImageResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// do nothing
	return nil
}

// NewImageResponse is the constructor method for the ImageResponse type
func NewImageResponse(img *cameraroll.Image) *ImageResponse {
	resp := ImageResponse{
		Image: img,
	}

	return &resp
}

// NewImageListResponse is the constructor method for a list of images
func NewImageListResponse(images []*cameraroll.Image) []render.Renderer {
	list := []render.Renderer{}

	for _, img := range images {
		list = append(list, NewImageResponse(img))
	}

	return list
}

// ImageCtx middleware is used to load an Image object from
// the URL parameters passed through as the request. In case
// the Image could not be found, we stop here and return a 404.
func (handler Handler) ImageCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var image *cameraroll.Image
		var imageID int64
		var err error

		// find the imageID from URL params
		if param := chi.URLParam(r, ParamImageID); len(param) > 0 {
			imageID, err = strconv.ParseInt(param, ParamNumberBase, ParamNumberBit)
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
			image, err = handler.Service.GetImageByID(r.Context(), imageID)
		} else {
			render.Render(w, r, ErrNotFound())
			return
		}

		if err != nil {
			render.Render(w, r, ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), imageKey, image)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RemoveTagFromImage removes a tag from the image
func (handler Handler) RemoveTagFromImage(w http.ResponseWriter, r *http.Request) {
	var tagID int64

	// find tagID from URL param
	if param := chi.URLParam(r, ParamTagID); len(param) > 0 {
		num, err := strconv.ParseInt(param, ParamNumberBase, ParamNumberBit)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		tagID = num
	}

	// find the image from context
	image := r.Context().Value(imageKey).(*cameraroll.Image)

	// remove the relationship from database
	if err := handler.Service.RemoveTagFromImage(r.Context(), image.ID, tagID); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// DeleteImage removes the image in the context
func (handler Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	image := r.Context().Value(imageKey).(*cameraroll.Image)

	if err := handler.Service.DeleteImageByID(r.Context(), image.ID); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// UpdateImage updates the image in the context
func (handler Handler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	image := r.Context().Value(imageKey).(*cameraroll.Image)

	imageReq := ImageRequest{}

	// unmarshal new image from request
	if err := render.Bind(r, &imageReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// add the new image to database
	newImage := imageReq.Image
	if err := handler.Service.UpdateImageByID(r.Context(), image.ID, newImage); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// GetImage returns the image in the context
func (handler Handler) GetImage(w http.ResponseWriter, r *http.Request) {
	image := r.Context().Value(imageKey).(*cameraroll.Image)

	if err := render.Render(w, r, NewImageResponse(image)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

}

// GetImages returns a list of images with pagination available
func (handler Handler) GetImages(w http.ResponseWriter, r *http.Request) {
	offset := PaginationDefaultOffset
	limit := PaginationDefaultLimit

	// try read offset from URL param
	if param := chi.URLParam(r, ParamOffset); len(param) > 0 {
		num, err := strconv.ParseUint(param, ParamNumberBase, ParamNumberBit)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		offset = num
	}

	// try read limit from URL param
	if param := chi.URLParam(r, ParamLimit); len(param) > 0 {
		num, err := strconv.ParseUint(param, ParamNumberBase, ParamNumberBit)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		limit = num
	}

	// query the database for list of images
	images, err := handler.Service.GetImages(r.Context(), offset, limit)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	// render response
	if err := render.RenderList(w, r, NewImageListResponse(images)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// AddImage adds a new image to the database
func (handler Handler) AddImage(w http.ResponseWriter, r *http.Request) {
	imageReq := ImageRequest{}

	// unmarshal new image from request
	if err := render.Bind(r, &imageReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// add the new image to database
	image := imageReq.Image
	if err := handler.Service.AddImage(r.Context(), image); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// render response
	render.Status(r, http.StatusOK)
	render.Render(w, r, NewImageResponse(image))
}
