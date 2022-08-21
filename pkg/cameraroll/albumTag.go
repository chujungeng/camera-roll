package cameraroll

import "context"

type AlbumTag struct {
	ID      int64 `json:"id,omitempty"`
	AlbumID int64 `json:"albumID"`
	TagID   int64 `json:"tagID"`
}

type AlbumTagService interface {
	AddTagToAlbum(ctx context.Context, albumID int64, tagID int64) error
	GetAlbumsWithTag(ctx context.Context, tagID int64, start uint64, count uint64) ([]*Album, error)
	GetTagsOfAlbum(ctx context.Context, albumID int64) ([]*Tag, error)
	RemoveTagFromAlbum(ctx context.Context, albumID int64, tagID int64) error
}
