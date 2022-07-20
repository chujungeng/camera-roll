package cameraroll

import (
	"context"
)

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type TagService interface {
	AddTag(ctx context.Context, tag *Tag) error
	GetTags(ctx context.Context) ([]*Tag, error)
	GetTagByID(ctx context.Context, id int64) (*Tag, error)
	UpdateTagByID(ctx context.Context, id int64, newTag *Tag) error
	DeleteTagByID(ctx context.Context, id int64) error
}
