package cameraroll

type Service interface {
	AlbumService
	ImageService
	TagService
	AlbumImageService
	AlbumTagService
	ImageTagService
}
