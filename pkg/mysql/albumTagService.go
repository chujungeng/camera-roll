package mysql

import (
	"context"
	"fmt"

	"chujungeng/camera-roll/pkg/cameraroll"
)

// RemoveTagFromAlbum removes a tag from an album
func (service Service) RemoveTagFromAlbum(ctx context.Context, albumID int64, tagID int64) error {
	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("RemoveTagFromAlbum albumID[%d] tagID[%d]: %v", albumID, tagID, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`DELETE FROM album_tags
		WHERE tag_id=? AND album_id=?`,
		tagID,
		albumID)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("RemoveTagFromAlbum albumID[%d] tagID[%d]: %v", albumID, tagID, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("RemoveTagFromAlbum albumID[%d] tagID[%d]: %v", albumID, tagID, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("RemoveTagFromAlbum albumID[%d] tagID[%d]: %v", albumID, tagID, err)
	}

	return nil
}

// GetAlbumsWithTag queries the database for certain amount of albums under a tag specified by tagID
// returns a slice of albums on success
func (service Service) GetAlbumsWithTag(ctx context.Context, tagID int64, start uint64, count uint64) ([]*cameraroll.Album, error) {
	// Album slice to hold the data from database query
	albums := []*cameraroll.Album{}

	// find prepared statement
	stmt := service.preparedStmts[keyQueryGetAlbumsWithTag]
	if stmt == nil {
		return nil, fmt.Errorf("GetAlbumsWithTag[%d] start[%d] count[%d]: Cannot find prepared sql query", tagID, start, count)
	}

	// execute the query
	rows, err := stmt.Query(tagID, start, count)

	// check if the query failed
	if err != nil {
		return nil, fmt.Errorf("GetAlbumsWithTag[%d] start[%d] count[%d]: %v", tagID, start, count, err)
	}

	defer rows.Close()

	// parse response
	for rows.Next() {
		alb := cameraroll.Album{}
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Description, &alb.CoverID); err != nil {
			return nil, fmt.Errorf("GetAlbumsWithTag[%d] start[%d] count[%d]: %v", tagID, start, count, err)
		}

		// add album to the return slice
		albums = append(albums, &alb)
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

// AddTagToAlbum adds a tag to an album
func (service Service) AddTagToAlbum(ctx context.Context, albumID int64, tagID int64) error {
	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AddTagToAlbum albumID[%d] tagID[%d]: %v", albumID, tagID, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`INSERT INTO album_tags(album_id, tag_id) 
		VALUES(?, ?)`,
		albumID,
		tagID)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("AddTagToAlbum albumID[%d] tagID[%d]: %v", albumID, tagID, err)
	}

	_, err = result.LastInsertId()
	if err != nil {
		return fmt.Errorf("AddTagToAlbum albumID[%d] tagID[%d]: %v", albumID, tagID, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("AddTagToAlbum albumID[%d] tagID[%d]: %v", albumID, tagID, err)
	}

	return nil
}
