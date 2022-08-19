package mysql

import (
	"context"
	"fmt"

	"chujungeng/camera-roll/pkg/cameraroll"
)

// RemoveTagFromImage removes a tag from an image
func (service Service) RemoveTagFromImage(ctx context.Context, imageID int64, tagID int64) error {
	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("RemoveTagFromImage imageID[%d] tagID[%d]: %v", imageID, tagID, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`DELETE FROM image_tags
		WHERE tag_id=? AND image_id=?`,
		tagID,
		imageID)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("RemoveTagFromImage imageID[%d] tagID[%d]: %v", imageID, tagID, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("RemoveTagFromImage imageID[%d] tagID[%d]: %v", imageID, tagID, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("RemoveTagFromImage imageID[%d] tagID[%d]: %v", imageID, tagID, err)
	}

	return nil
}

// GetImagesWithTag queries the database for certain amount of images under a tag specified by tagID
// returns a slice of images on success
func (service Service) GetImagesWithTag(ctx context.Context, tagID int64, start uint64, count uint64) ([]*cameraroll.Image, error) {
	// Image slice to hold the data from database query
	images := []*cameraroll.Image{}

	// find prepared statement
	stmt := service.preparedStmts[keyQueryGetImagesWithTag]
	if stmt == nil {
		return nil, fmt.Errorf("GetImagesWithTag[%d] start[%d] count[%d]: Cannot find prepared sql query", tagID, start, count)
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("GetImagesWithTag[%d] start[%d] count[%d]: %v", tagID, start, count, err)
	}
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, stmt)

	// execute the query
	rows, err := txStmt.QueryContext(ctx, tagID, start, count)

	// check if the query failed
	if err != nil {
		return nil, fmt.Errorf("GetImagesWithTag[%d] start[%d] count[%d]: %v", tagID, start, count, err)
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
			return nil, fmt.Errorf("GetImagesWithTag[%d] start[%d] count[%d]: %v", tagID, start, count, err)
		}

		// add image to the return slice
		images = append(images, &img)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("GetImagesWithTag[%d] start[%d] count[%d]: %v", tagID, start, count, err)
	}

	return images, nil
}

// AddTagToImage adds a tag to an image
func (service Service) AddTagToImage(ctx context.Context, imageID int64, tagID int64) error {
	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AddTagToImage imageID[%d] tagID[%d]: %v", imageID, tagID, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`INSERT INTO image_tags(image_id, tag_id)
        VALUES(?, ?)`,
		imageID,
		tagID)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("AddTagToImage imageID[%d] tagID[%d]: %v", imageID, tagID, err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("AddTagToImage imageID[%d] tagID[%d]: %v", imageID, tagID, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("AddTagToImage imageID[%d] tagID[%d]: %v", imageID, tagID, err)
	}

	return nil
}
