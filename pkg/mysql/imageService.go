package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"chujungeng/camera-roll/pkg/cameraroll"
)

// DeleteImageByID removes an image from database
func (service Service) DeleteImageByID(ctx context.Context, id int64) error {
	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("DeleteImageByID [%d]: %v", id, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`DELETE FROM images 
		WHERE id=?`,
		id)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("DeleteImageByID [%d]: %v", id, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteImageByID [%d]: %v", id, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("DeleteImageByID [%d]: %v", id, err)
	}

	return nil
}

// UpdateImageByID updates an image's path, title and description
func (service Service) UpdateImageByID(ctx context.Context, id int64, newImg *cameraroll.Image) error {
	if newImg == nil {
		return fmt.Errorf("UpdateImageByID [%d]: null pointer error", id)
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("UpdateImageByID [%d]: %v", id, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`UPDATE images 
		SET path=?, title=?, description=? 
		WHERE id=?`,
		newImg.Path,
		newImg.Title,
		newImg.Description,
		id)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("UpdateImageByID [%d]: %v", id, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateImageByID [%d]: %v", id, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("UpdateImageByID [%d]: %v", id, err)
	}

	return nil
}

// GetImageByID queries the database for the image specified by its ID
func (service Service) GetImageByID(ctx context.Context, id int64) (*cameraroll.Image, error) {
	img := cameraroll.Image{}

	// find prepared statement
	stmt := service.preparedStmts[keyQueryGetImageByID]
	if stmt == nil {
		return nil, fmt.Errorf("GetImageByID [%d]: Cannot find prepared sql query", id)
	}

	// execute the query
	row := stmt.QueryRow(id)

	// parse response
	if err := row.Scan(&img.ID, &img.Path, &img.Title, &img.Description, &img.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("GetImageByID[%d]: no such image", id)
		}

		return nil, fmt.Errorf("GetImageByID[%d]: %v", id, err)
	}

	return &img, nil
}

// GetImages queries the database for certain amount of images from a starting index
func (service Service) GetImages(ctx context.Context, start uint64, count uint64) ([]*cameraroll.Image, error) {
	// Image slice to hold the data from database query
	images := []*cameraroll.Image{}

	// find prepared statement
	stmt := service.preparedStmts[keyQueryGetImages]
	if stmt == nil {
		return nil, fmt.Errorf("GetImages start[%d] count[%d]: Cannot find prepared sql query", start, count)
	}

	// execute the query
	rows, err := stmt.Query(start, count)

	// check if the query failed
	if err != nil {
		return nil, fmt.Errorf("GetImages start[%d] count[%d]: %v", start, count, err)
	}

	defer rows.Close()

	// parse response
	for rows.Next() {
		img := cameraroll.Image{}
		if err := rows.Scan(&img.ID, &img.Path, &img.Title, &img.Description, &img.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetImages start[%d] count[%d]: %v", start, count, err)
		}

		// add image to the return slice
		images = append(images, &img)
	}

	return images, nil
}

// AddImage adds 1 image to the database,
// updating the image's ID upon success
func (service Service) AddImage(ctx context.Context, image *cameraroll.Image) error {
	if image == nil {
		return fmt.Errorf("AddImage : null pointer error")
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AddImage [%s]: %v", image.Path, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`INSERT INTO images (path, title, description) 
		VALUES (?, ?, ?)`,
		image.Path,
		image.Title,
		image.Description)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("AddImage [%s]: %v", image.Path, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("AddImage [%s]: %v", image.Path, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("AddImage [%s]: %v", image.Path, err)
	}

	image.ID = id

	return nil
}
