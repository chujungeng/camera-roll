package cameraroll

import "context"

type AlbumImage struct {
	ID      int64 `json:"id,omitempty"`
	AlbumID int64 `json:"albumID"`
	ImageID int64 `json:"imageID"`
}

type AlbumImageService interface {
	AddImageToAlbum(ctx context.Context, albumID int64, imageID int64) error
	GetImagesFromAlbum(ctx context.Context, id int64) ([]*Image, error)
	RemoveImageFromAlbum(ctx context.Context, albumID int64, imageID int64) error
}
