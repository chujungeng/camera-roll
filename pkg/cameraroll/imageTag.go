package cameraroll

import "context"

type ImageTag struct {
	ID      int64 `json:"id,omitempty"`
	ImageID int64 `json:"imageID"`
	TagID   int64 `json:"tagID"`
}

type ImageTagService interface {
	AddTagToImage(ctx context.Context, imageID int64, tagID int64) error
	GetImagesWithTag(ctx context.Context, tagID int64, start uint64, count uint64) ([]*Image, error)
	GetTagsOfImage(ctx context.Context, imageID int64) ([]*Tag, error)
	RemoveTagFromImage(ctx context.Context, imageID int64, tagID int64) error
}
