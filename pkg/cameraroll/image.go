package cameraroll

import (
	"context"
	"time"
)

type Image struct {
	ID          int64     `json:"id"`
	Path        string    `json:"path"`
	Thumbnail   string    `json:"thumbnail"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type ImageService interface {
	AddImage(ctx context.Context, image *Image) error
	GetImages(ctx context.Context, start uint64, count uint64) ([]*Image, error)
	GetImageByID(ctx context.Context, id int64) (*Image, error)
	UpdateImageByID(ctx context.Context, id int64, newImg *Image) error
	DeleteImageByID(ctx context.Context, id int64) error
}
