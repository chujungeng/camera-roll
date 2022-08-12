package api

import (
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/nfnt/resize"

	"chujungeng/camera-roll/pkg/cameraroll"
	"chujungeng/camera-roll/pkg/url"
)

const (
	ParamImageID          = "imageID"
	ParamImageTitle       = "title"
	ParamImageDescription = "description"
	ParamImageFile        = "image"
)

const (
	MaxImageSize = 2 << 20
)

// ImageRouterPublic specifies all the public routes related to images
func (handler Handler) ImageRouterPublic() chi.Router {
	r := chi.NewRouter()

	r.With(Pagination).Get("/", handler.GetImages) // GET /images

	r.Route("/{imageID}", func(r chi.Router) {
		r.Use(handler.ImageCtx)      // Load the *Image on the request context
		r.Get("/", handler.GetImage) // GET /images/123
	})

	return r
}

// ImageRouterProtected specifies all the protected routes related to images
func (handler Handler) ImageRouterProtected() chi.Router {
	r := chi.NewRouter()

	r.With(Pagination).Get("/", handler.GetImages) // GET /admin/images
	r.Post("/", handler.AddImage)                  // POST /admin/images

	r.Route("/{imageID}", func(r chi.Router) {
		r.Use(handler.ImageCtx)            // Load the *Image on the request context
		r.Get("/", handler.GetImage)       // GET /admin/images/123
		r.Put("/", handler.UpdateImage)    // PUT /admin/images/123
		r.Delete("/", handler.DeleteImage) // DELETE /admin/images/123

		r.Delete("/tags/{tagID}", handler.RemoveTagFromImage) // DELETE /admin/images/123/tags/789
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

	DeleteAssetFromFilesystem(image.Path)
	DeleteAssetFromFilesystem(image.Thumbnail)

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

	// find the pageID from context
	page := r.Context().Value(pageIDKey).(int)
	if page > 1 {
		offset = PaginationDefaultLimit * (uint64(page) - 1)
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

func verifyImageFile(imageFile multipart.File) error {
	// read 512 bytes from the file for image validation
	buff := make([]byte, 512)
	if _, err := imageFile.Read(buff); err != nil {
		return err
	}

	// check if the uploaded file really is an image
	if !strings.HasPrefix(http.DetectContentType(buff), "image") {
		return errors.New("invalid file type")
	}

	// rewind the imageFile
	imageFile.Seek(0, 0)

	return nil
}

func saveImageFile(imageFile multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// find the image's file extension
	fileNameSlice := strings.Split(fileHeader.Filename, ".")
	fileType := fileNameSlice[len(fileNameSlice)-1]

	// get a uuid for the file's new name
	fileNameNew := fmt.Sprintf("%s.%s", uuid.New().String(), fileType)

	// construct the file path
	fileDest := filepath.Join(StaticFileDir(), fileNameNew)

	// copy the file to static asset folder
	f, err := os.Create(fileDest)
	if err != nil {
		return fileNameNew, err
	}
	defer f.Close()
	io.Copy(f, imageFile)

	// rewind the imageFile
	imageFile.Seek(0, 0)

	return fileNameNew, nil
}

type ThumbnailStats struct {
	Path            string
	Width           int
	Height          int
	ThumbnailWidth  int
	ThumbnailHeight int
}

func createImageThumbnail(imageFile multipart.File, fileHeader *multipart.FileHeader) (ThumbnailStats, error) {
	const (
		ThumbnailMaxWidthPx  = 400
		ThumbnailMaxHeightPx = 400
	)

	// find the image's file extension
	fileNameSlice := strings.Split(fileHeader.Filename, ".")
	fileType := fileNameSlice[len(fileNameSlice)-1]

	// get a uuid for the thumbnail's filename
	fileNameNew := fmt.Sprintf("%s.%s", uuid.New().String(), fileType)

	// construct the file path
	fileDest := filepath.Join(StaticFileDir(), fileNameNew)

	// open the destination file
	f, err := os.Create(fileDest)
	if err != nil {
		return ThumbnailStats{Path: fileNameNew}, err
	}
	defer f.Close()

	// decode the original image
	img, err := jpeg.Decode(imageFile)
	if err != nil {
		return ThumbnailStats{Path: fileNameNew}, err
	}

	// create the thumbnail
	thumb := resize.Thumbnail(ThumbnailMaxWidthPx, ThumbnailMaxHeightPx, img, resize.Lanczos3)

	// save the thumbnail to static folder
	jpeg.Encode(f, thumb, nil)

	// rewind the imageFile
	imageFile.Seek(0, 0)

	// get all the stats
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	thumbWidth := thumb.Bounds().Dx()
	thumbHeight := thumb.Bounds().Dy()

	return ThumbnailStats{
		Path:            fileNameNew,
		Width:           width,
		Height:          height,
		ThumbnailWidth:  thumbWidth,
		ThumbnailHeight: thumbHeight}, nil
}

// AddImage adds a new image to the database
func (handler Handler) AddImage(w http.ResponseWriter, r *http.Request) {
	imageReq := ImageRequest{&cameraroll.Image{}}

	// // unmarshal new image from request
	// if err := render.Bind(r, &imageReq); err != nil {
	// 	render.Render(w, r, ErrInvalidRequest(err))
	// 	return
	// }

	// parse the form from request
	if err := r.ParseMultipartForm(MaxImageSize); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// find title from form data
	title := r.Form.Get(ParamImageTitle)
	if len(title) > 0 {
		imageReq.Title = title
	}

	// find description from form data
	desc := r.Form.Get(ParamImageDescription)
	if len(desc) > 0 {
		imageReq.Description = desc
	}

	// find the image file from form data
	imageFile, fileHeader, err := r.FormFile(ParamImageFile)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(errors.New("image file was empty")))
		return
	}
	defer imageFile.Close()

	// check if the uploaded file is an image
	if err := verifyImageFile(imageFile); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// save the image to static folder
	fileNameNew, err := saveImageFile(imageFile, fileHeader)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// create a thumbnail
	thumbnail, err := createImageThumbnail(imageFile, fileHeader)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// update the image's path
	imagePath := url.Join(handler.rootURL, staticFileURL, fileNameNew)
	imageReq.Path = imagePath

	// update thumbnail path
	thumbnailPath := url.Join(handler.rootURL, staticFileURL, thumbnail.Path)
	imageReq.Thumbnail = thumbnailPath

	// update image stats
	imageReq.Width = thumbnail.Width
	imageReq.Height = thumbnail.Height
	imageReq.ThumbnailWidth = thumbnail.ThumbnailWidth
	imageReq.ThumbnailHeight = thumbnail.ThumbnailHeight

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
