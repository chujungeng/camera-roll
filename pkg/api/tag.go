package api

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"chujungeng/camera-roll/pkg/cameraroll"
)

const (
	ParamTagID = "tagID"
)

// TagRouterPublic specifies all the public routes related to tags
func (handler Handler) TagRouterPublic() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.GetTags) // GET /tags

	r.Route("/{tagID}", func(r chi.Router) {
		r.Use(handler.TagCtx)                      // Load the *Tag on the request context
		r.Get("/", handler.GetTag)                 // GET /tags/123
		r.Get("/albums", handler.GetAlbumsWithTag) // GET /tags/123/albums
		r.Get("/images", handler.GetImagesWithTag) // GET /tags/123/images
	})

	return r
}

// TagRouterProtected specifies all the protected routes related to tags
func (handler Handler) TagRouterProtected() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handler.GetTags) // GET /admin/tags
	r.Post("/", handler.AddTag) // POST /admin/tags

	r.Route("/{tagID}", func(r chi.Router) {
		r.Use(handler.TagCtx)            // Load the *Tag on the request context
		r.Get("/", handler.GetTag)       // GET /admin/tags/123
		r.Put("/", handler.UpdateTag)    // PUT /admin/tags/123
		r.Delete("/", handler.DeleteTag) // DELETE /admin/tags/123

		r.Get("/albums", handler.GetAlbumsWithTag) // GET /admin/tags/123/albums
		r.Get("/images", handler.GetImagesWithTag) // GET /admin/tags/123/images
	})

	return r
}

// TagRequest is the request body of tags' CRUD operations
type TagRequest struct {
	*cameraroll.Tag
}

// Bind preprocesses the request for some basic error checking
func (req *TagRequest) Bind(r *http.Request) error {
	// Return an error to avoid a nil pointer dereference.
	if req.Tag == nil {
		return errors.New("missing required Tag fields")
	}

	return nil
}

// TagResponse is the response body of tags' CRUD operations
type TagResponse struct {
	*cameraroll.Tag
}

// Render preprocess the response before it's sent to the wire
func (rsp *TagResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// do nothing
	return nil
}

// NewTagResponse is the constructor method for TagResponse type
func NewTagResponse(tag *cameraroll.Tag) *TagResponse {
	resp := TagResponse{Tag: tag}

	return &resp
}

// NewTagListResponse is the constructor method for a list of TagResponses
func NewTagListResponse(tags []*cameraroll.Tag) []render.Renderer {
	list := []render.Renderer{}

	for _, tag := range tags {
		list = append(list, NewTagResponse(tag))
	}

	return list
}

// TagCtx middleware is used to load an Tag object from
// the URL parameters passed through as the request. In case
// the Tag could not be found, we stop here and return a 404.
func (handler Handler) TagCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tag *cameraroll.Tag
		var tagID int64
		var err error

		// find the tagID from URL params
		if param := chi.URLParam(r, ParamTagID); len(param) > 0 {
			tagID, err = strconv.ParseInt(param, ParamNumberBase, ParamNumberBit)
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
			tag, err = handler.Service.GetTagByID(r.Context(), tagID)
		} else {
			render.Render(w, r, ErrNotFound())
			return
		}

		if err != nil {
			render.Render(w, r, ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), tagKey, tag)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetImagesWithTag returns all the images under specified tag
func (handler Handler) GetImagesWithTag(w http.ResponseWriter, r *http.Request) {
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

	tag := r.Context().Value(tagKey).(*cameraroll.Tag)

	images, err := handler.Service.GetImagesWithTag(r.Context(), tag.ID, offset, limit)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// render response
	if err := render.RenderList(w, r, NewImageListResponse(images)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// GetAlbumsWithTag returns all the albums under specified tag
func (handler Handler) GetAlbumsWithTag(w http.ResponseWriter, r *http.Request) {
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

	tag := r.Context().Value(tagKey).(*cameraroll.Tag)

	albums, err := handler.Service.GetAlbumsWithTag(r.Context(), tag.ID, offset, limit)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// render response
	if err := render.RenderList(w, r, NewAlbumListResponse(albums)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// DeleteTag removes the tag in the context
func (handler Handler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	tag := r.Context().Value(tagKey).(*cameraroll.Tag)

	if err := handler.Service.DeleteTagByID(r.Context(), tag.ID); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// UpdateTag updates the tag in the context
func (handler Handler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	tag := r.Context().Value(tagKey).(*cameraroll.Tag)

	tagReq := TagRequest{}

	// unmarshal new tag from request
	if err := render.Bind(r, &tagReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// add the new tag to database
	newTag := tagReq.Tag
	if err := handler.Service.UpdateTagByID(r.Context(), tag.ID, newTag); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// GetTag returns the tag in the context
func (handler Handler) GetTag(w http.ResponseWriter, r *http.Request) {
	tag := r.Context().Value(tagKey).(*cameraroll.Tag)

	if err := render.Render(w, r, NewTagResponse(tag)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

}

// GetTags returns a list of tags with pagination available
func (handler Handler) GetTags(w http.ResponseWriter, r *http.Request) {

	// query the database for list of tags
	tags, err := handler.Service.GetTags(r.Context())
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	// render response
	if err := render.RenderList(w, r, NewTagListResponse(tags)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// AddTag adds a new tag to the database
func (handler Handler) AddTag(w http.ResponseWriter, r *http.Request) {
	tagReq := TagRequest{}

	// unmarshal new tag from request
	if err := render.Bind(r, &tagReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// add the new tag to database
	tag := tagReq.Tag
	if err := handler.Service.AddTag(r.Context(), tag); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// render response
	render.Status(r, http.StatusOK)
	render.Render(w, r, NewTagResponse(tag))
}
