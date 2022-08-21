package mysql

import (
	"context"
	"fmt"

	"chujungeng/camera-roll/pkg/cameraroll"
)

// RemoveImageFromAlbum removes an image from album
func (service Service) RemoveImageFromAlbum(ctx context.Context, albumID int64, imageID int64) error {
	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("RemoveImageFromAlbum albumID[%d] imageID[%d]: %v", albumID, imageID, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`DELETE FROM image_albums 
		WHERE album_id=? AND image_id=?`,
		albumID,
		imageID)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("RemoveImageFromAlbum albumID[%d] imageID[%d]: %v", albumID, imageID, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("RemoveImageFromAlbum albumID[%d] imageID[%d]: %v", albumID, imageID, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("RemoveImageFromAlbum albumID[%d] imageID[%d]: %v", albumID, imageID, err)
	}

	return nil
}

// GetImagesFromAlbum gets all the images from an album
func (service Service) GetImagesFromAlbum(ctx context.Context, id int64) ([]*cameraroll.Image, error) {
	// Image slice to hold the data from database query
	images := []*cameraroll.Image{}

	// find prepared statement
	stmt := service.preparedStmts[keyQueryGetImagesFromAlbum]
	if stmt == nil {
		return nil, fmt.Errorf("GetImagesFromAlbum id[%d]: Cannot find prepared sql query", id)
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("GetImagesFromAlbum id[%d]: Cannot start a database transaction. err[%v]", id, err)
	}
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, stmt)

	// execute the query
	rows, err := txStmt.QueryContext(ctx, id)

	// check if the query failed
	if err != nil {
		return nil, fmt.Errorf("GetImagesFromAlbum id[%d]: %v", id, err)
	}

	defer rows.Close()

	// parse response
	for rows.Next() {
		img := cameraroll.Image{}
		if err := rows.Scan(
			&img.ID,
			&img.Path,
			&img.Width,
			&img.Height,
			&img.Thumbnail,
			&img.ThumbnailWidth,
			&img.ThumbnailHeight,
			&img.Title,
			&img.Description,
			&img.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetImagesFromAlbum id[%d]: %v", id, err)
		}

		// add image to the return slice
		images = append(images, &img)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("GetImagesFromAlbum id[%d]: %v", id, err)
	}

	return images, nil
}

// GetAlbumsOfImage gets all the albums that an image belongs to
func (service Service) GetAlbumsOfImage(ctx context.Context, id int64) ([]*cameraroll.Album, error) {
	// Album slice to hold the data from database query
	albums := []*cameraroll.Album{}

	// find prepared statement
	stmt := service.preparedStmts[keyQueryGetAlbumsOfImage]
	if stmt == nil {
		return nil, fmt.Errorf("GetAlbumsOfImage id[%d]: Cannot find prepared sql query", id)
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("GetAlbumsOfImage id[%d]: Cannot start a database transaction. err[%v]", id, err)
	}
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, stmt)

	// execute the query
	rows, err := txStmt.QueryContext(ctx, id)

	// check if the query failed
	if err != nil {
		return nil, fmt.Errorf("GetAlbumsOfImage id[%d]: %v", id, err)
	}

	defer rows.Close()

	// parse response
	for rows.Next() {
		alb := cameraroll.Album{}
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Description, &alb.CreatedAt, &alb.CoverID); err != nil {
			return nil, fmt.Errorf("GetAlbumsOfImage id[%d]: %v", id, err)
		}

		// add album to the return slice
		albums = append(albums, &alb)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("GetAlbumsOfImage id[%d]: %v", id, err)
	}

	// query database for album covers
	for _, alb := range albums {
		if alb.CoverID.Valid {
			img, _ := service.GetImageByID(ctx, alb.CoverID.Int64)
			if img != nil {
				alb.Cover = img
			}
		}
	}

	return albums, nil
}

// AddImageToAlbum adds an image to an album
func (service Service) AddImageToAlbum(ctx context.Context, albumID int64, imageID int64) error {
	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AddImageToAlbum imageID[%d] albumID[%d]: %v", imageID, albumID, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`INSERT INTO image_albums(album_id, image_id) 
		VALUES(?, ?)`,
		albumID,
		imageID)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("AddImageToAlbum imageID[%d] albumID[%d]: %v", imageID, albumID, err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("AddImageToAlbum imageID[%d] albumID[%d]: %v", imageID, albumID, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("AddImageToAlbum imageID[%d] albumID[%d]: %v", imageID, albumID, err)
	}

	return nil
}
