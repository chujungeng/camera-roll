package cameraroll

import (
	"context"
	"database/sql"
	"time"
)

type Album struct {
	ID          int64         `json:"id"`
	CreatedAt   time.Time     `json:"created_at,omitempty"`
	Title       string        `json:"title,omitempty"`
	Description string        `json:"description,omitempty"`
	CoverID     sql.NullInt64 `json:"cover_id,omitempty"`
	Cover       *Image        `json:"cover,omitempty"`
}

type AlbumService interface {
	AddAlbum(ctx context.Context, album *Album) error
	GetAlbums(ctx context.Context, start uint64, count uint64) ([]*Album, error)
	GetAlbumByID(ctx context.Context, id int64) (*Album, error)
	UpdateAlbumByID(ctx context.Context, id int64, newAlb *Album) error
	DeleteAlbumByID(ctx context.Context, id int64) error
}
