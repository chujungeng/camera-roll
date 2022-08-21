package mysql

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// timeouts
const modelInitTimeout = 7 * time.Second

// keys for prepared sql statements
const (
	keyQueryGetImages          = "GetImages"
	keyQueryGetImageByID       = "GetImageByID"
	keyQueryGetTags            = "GetTags"
	keyQueryGetTagByID         = "GetTagByID"
	keyQueryGetAlbums          = "GetAlbums"
	keyQueryGetAlbumByID       = "GetAlbumByID"
	keyQueryGetImagesFromAlbum = "GetImagesFromAlbum"
	keyQueryGetAlbumsOfImage   = "GetAlbumsOfImage"
	keyQueryGetAlbumsWithTag   = "GetAlbumsWithTag"
	keyQueryGetTagsOfAlbum     = "GetTagsOfAlbum"
	keyQueryGetImagesWithTag   = "GetImagesWithTag"
	keyQueryGetTagsOfImage     = "GetTagsOfImage"
)

type Service struct {
	// database connection
	db *sql.DB

	// prepared sql statements, safe for concurrent READs
	// unsafe for concurrent READ and WRITE or concurrent WRITEs
	preparedStmts map[string]*sql.Stmt
}

// set up prepared sql statements for faster repeated execution
func (service *Service) createPreparedStmts() error {
	// sql templates
	queries := map[string]string{
		keyQueryGetImages:    `SELECT * FROM images ORDER BY created_at DESC LIMIT ?, ?`,
		keyQueryGetImageByID: `SELECT * FROM images WHERE id=?`,
		keyQueryGetTags:      `SELECT * FROM tags ORDER BY id`,
		keyQueryGetTagByID:   `SELECT * FROM tags WHERE id=?`,
		keyQueryGetAlbums:    `SELECT * FROM albums ORDER BY created_at DESC LIMIT ?, ?`,
		keyQueryGetAlbumByID: `SELECT * FROM albums WHERE id=?`,
		keyQueryGetImagesFromAlbum: `SELECT images.id, images.path, images.width, images.height, images.thumbnail, images.width_thumb, images.height_thumb, images.title, images.description, images.created_at
									FROM albums JOIN image_albums 
									ON albums.id=image_albums.album_id 
									JOIN images 
									ON image_albums.image_id=images.id 
									WHERE albums.id=?
									ORDER BY image_albums.id DESC`,
		keyQueryGetAlbumsOfImage: `SELECT albums.id, albums.title, albums.description, albums.created_at, albums.cover_id
									FROM images JOIN image_albums 
									ON images.id=image_albums.image_id 
									JOIN albums 
									ON image_albums.album_id=albums.id 
									WHERE images.id=?
									ORDER BY image_albums.id DESC`,
		keyQueryGetAlbumsWithTag: `SELECT albums.id, albums.title, albums.description, albums.created_at, albums.cover_id
									FROM tags JOIN album_tags
									ON tags.id=album_tags.tag_id
									JOIN albums
									ON albums.id=album_tags.album_id
									WHERE tags.id=?
									ORDER BY album_tags.id DESC
									LIMIT ?, ?`,
		keyQueryGetTagsOfAlbum: `SELECT tags.id, tags.name
								FROM albums JOIN album_tags
								ON albums.id=album_tags.album_id
								JOIN tags
								ON album_tags.tag_id=tags.id
								WHERE albums.id=?
								ORDER BY tags.id DESC`,
		keyQueryGetImagesWithTag: `SELECT images.id, images.path, images.width, images.height, images.thumbnail, images.width_thumb, images.height_thumb, images.title, images.description, images.created_at
									FROM tags JOIN image_tags
									ON tags.id=image_tags.tag_id
									JOIN images
									ON images.id=image_tags.image_id
									WHERE tags.id=?
									ORDER BY image_tags.id DESC
									LIMIT ?, ?`,
		keyQueryGetTagsOfImage: `SELECT tags.id, tags.name
								FROM images JOIN image_tags
								ON images.id=image_tags.image_id
								JOIN tags
								ON image_tags.tag_id=tags.id
								WHERE images.id=?
								ORDER BY tags.id DESC`,
	}

	var err error
	ctx, failure := context.WithTimeout(context.Background(), modelInitTimeout)
	defer failure()

	for k, q := range queries {
		if service.preparedStmts[k] == nil {
			if service.preparedStmts[k], err = service.db.PrepareContext(ctx, q); err != nil {
				return err
			}
		}
	}

	return nil
}

// close all prepared statements and remove them from preparedStmts map
func (service *Service) closePreparedStmts() {
	if service.preparedStmts == nil {
		return
	}

	for key, stmt := range service.preparedStmts {
		if stmt != nil {
			stmt.Close()
		}

		delete(service.preparedStmts, key)
	}
}

// Cleanup closes all the reusable resources
func (service *Service) Cleanup() {
	service.closePreparedStmts()
}

func NewService(db *sql.DB) *Service {
	if db == nil {
		log.Fatalln("NewService: Null pointer error")
		return nil
	}

	service := Service{
		db:            db,
		preparedStmts: make(map[string]*sql.Stmt),
	}

	service.createPreparedStmts()

	return &service
}
