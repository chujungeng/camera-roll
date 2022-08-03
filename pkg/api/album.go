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
	ParamAlbumID = "albumID"
)

// AlbumRouterPublic specifies all the public routes related to albums
func (handler Handler) AlbumRouterPublic() chi.Router {
	r := chi.NewRouter()

	r.With(Pagination).Get("/", handler.GetAlbums) // GET /albums

	r.Route("/{albumID}", func(r chi.Router) {
		r.Use(handler.AlbumCtx)      // Load the *Album on the request context
		r.Get("/", handler.GetAlbum) // GET /albums/123

		r.Get("/images", handler.GetImagesFromAlbum) // GET /albums/123/images
	})

	return r
}

// AlbumRouterProtected contains all the album routes that should be protected
func (handler Handler) AlbumRouterProtected() chi.Router {
	r := chi.NewRouter()

	r.With(Pagination).Get("/", handler.GetAlbums) // GET /admin/albums
	r.Post("/", handler.AddAlbum)                  // POST /admin/albums

	r.Route("/{albumID}", func(r chi.Router) {
		r.Use(handler.AlbumCtx)            // Load the *Album on the request context
		r.Get("/", handler.GetAlbum)       // GET /admin/albums/123
		r.Put("/", handler.UpdateAlbum)    // PUT /admin/albums/123
		r.Delete("/", handler.DeleteAlbum) // DELETE /admin/albums/123

		r.Get("/images", handler.GetImagesFromAlbum)                // GET /admin/albums/123/images
		r.Delete("/images/{imageID}", handler.RemoveImageFromAlbum) // DELETE /admin/albums/123/images/456

		r.Delete("/tags/{tagID}", handler.RemoveTagFromAlbum) // DELETE /admin/albums/123/tags/789
	})

	return r
}

// AlbumRequest is the request body of albums' CRUD operations
type AlbumRequest struct {
	*cameraroll.Album
}

// Bind preprocesses the request for some basic error checking
func (req *AlbumRequest) Bind(r *http.Request) error {
	// Return an error to avoid a nil pointer dereference.
	if req.Album == nil {
		return errors.New("missing required Album fields")
	}

	return nil
}

// AlbumResponse is the response body of albums' CRUD operations
type AlbumResponse struct {
	*cameraroll.Album
}

// Render preprocess the response before it's sent to the wire
func (rsp *AlbumResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// do nothing
	return nil
}

// NewAlbumResponse is the constructor method for AlbumResponse type
func NewAlbumResponse(album *cameraroll.Album) *AlbumResponse {
	resp := AlbumResponse{Album: album}

	return &resp
}

// NewAlbumListResponse is the constructor method for a list of AlbumResponses
func NewAlbumListResponse(albums []*cameraroll.Album) []render.Renderer {
	list := []render.Renderer{}

	for _, album := range albums {
		list = append(list, NewAlbumResponse(album))
	}

	return list
}

// AlbumCtx middleware is used to load an Album object from
// the URL parameters passed through as the request. In case
// the Album could not be found, we stop here and return a 404.
func (handler Handler) AlbumCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var album *cameraroll.Album
		var albumID int64
		var err error

		// find the albumID from URL params
		if param := chi.URLParam(r, ParamAlbumID); len(param) > 0 {
			albumID, err = strconv.ParseInt(param, ParamNumberBase, ParamNumberBit)
			if err != nil {
				render.Render(w, r, ErrInvalidRequest(err))
				return
			}
			album, err = handler.Service.GetAlbumByID(r.Context(), albumID)
		} else {
			render.Render(w, r, ErrNotFound())
			return
		}

		if err != nil {
			render.Render(w, r, ErrNotFound())
			return
		}

		ctx := context.WithValue(r.Context(), albumKey, album)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RemoveTagFromAlbum removes a tag from the album
func (handler Handler) RemoveTagFromAlbum(w http.ResponseWriter, r *http.Request) {
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

	// find the album from context
	album := r.Context().Value(albumKey).(*cameraroll.Album)

	// remove the relationship from database
	if err := handler.Service.RemoveTagFromAlbum(r.Context(), album.ID, tagID); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// RemoveImageFromAlbum removes an image from the album
func (handler Handler) RemoveImageFromAlbum(w http.ResponseWriter, r *http.Request) {
	var imageID int64

	// find imageID from URL param
	if param := chi.URLParam(r, ParamImageID); len(param) > 0 {
		num, err := strconv.ParseInt(param, ParamNumberBase, ParamNumberBit)
		if err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}

		imageID = num
	}

	// find the album from context
	album := r.Context().Value(albumKey).(*cameraroll.Album)

	// remove the relationship from database
	if err := handler.Service.RemoveImageFromAlbum(r.Context(), album.ID, imageID); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// GetAlbumImages returns all the images from an album
func (handler Handler) GetImagesFromAlbum(w http.ResponseWriter, r *http.Request) {
	album := r.Context().Value(albumKey).(*cameraroll.Album)

	images, err := handler.Service.GetImagesFromAlbum(r.Context(), album.ID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	if err := render.RenderList(w, r, NewImageListResponse(images)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// DeleteAlbum removes the album in the context
func (handler Handler) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	album := r.Context().Value(albumKey).(*cameraroll.Album)

	if err := handler.Service.DeleteAlbumByID(r.Context(), album.ID); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// UpdateAlbum updates the album in the context
func (handler Handler) UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	album := r.Context().Value(albumKey).(*cameraroll.Album)

	albumReq := AlbumRequest{}

	// unmarshal new album from request
	if err := render.Bind(r, &albumReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// add the new album to database
	newAlbum := albumReq.Album
	if err := handler.Service.UpdateAlbumByID(r.Context(), album.ID, newAlbum); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Status(r, http.StatusOK)
}

// GetAlbum returns the album in the context
func (handler Handler) GetAlbum(w http.ResponseWriter, r *http.Request) {
	album := r.Context().Value(albumKey).(*cameraroll.Album)

	if err := render.Render(w, r, NewAlbumResponse(album)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

}

// GetAlbums returns a list of albums with pagination available
func (handler Handler) GetAlbums(w http.ResponseWriter, r *http.Request) {
	offset := PaginationDefaultOffset
	limit := PaginationDefaultLimit

	// find the pageID from context
	page := r.Context().Value(pageIDKey).(int)
	if page > 1 {
		offset = PaginationDefaultLimit * (uint64(page) - 1)
	}

	// query the database for list of albums
	albums, err := handler.Service.GetAlbums(r.Context(), offset, limit)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	// render response
	if err := render.RenderList(w, r, NewAlbumListResponse(albums)); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}
}

// AddAlbum adds a new album to the database
func (handler Handler) AddAlbum(w http.ResponseWriter, r *http.Request) {
	albumReq := AlbumRequest{}

	// unmarshal new album from request
	if err := render.Bind(r, &albumReq); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// add the new album to database
	album := albumReq.Album
	if err := handler.Service.AddAlbum(r.Context(), album); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// render response
	render.Status(r, http.StatusOK)
	render.Render(w, r, NewAlbumResponse(album))
}
