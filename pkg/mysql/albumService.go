package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"chujungeng/camera-roll/pkg/cameraroll"
)

// DeleteAlbumByID removes an album from database
func (service Service) DeleteAlbumByID(ctx context.Context, id int64) error {
	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("DeleteAlbumByID [%d]: %v", id, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`DELETE FROM albums 
		WHERE id=?`,
		id)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("DeleteAlbumByID [%d]: %v", id, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteAlbumByID [%d]: %v", id, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("DeleteAlbumByID [%d]: %v", id, err)
	}

	return nil
}

// UpdateAlbumByID updates an album's title, description and cover_id
func (service Service) UpdateAlbumByID(ctx context.Context, id int64, newAlb *cameraroll.Album) error {
	if newAlb == nil {
		return fmt.Errorf("UpdateAlbumByID [%d]: null pointer error", id)
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("UpdateAlbumByID [%d]: %v", id, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`UPDATE albums 
		SET title=?, description=?, cover_id=? 
		WHERE id=?`,
		newAlb.Title,
		newAlb.Description,
		newAlb.CoverID,
		id)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("UpdateAlbumByID [%d]: %v", id, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateAlbumByID [%d]: %v", id, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("UpdateAlbumByID [%d]: %v", id, err)
	}

	return nil
}

// GetAlbumByID queries the database for the album specified by its ID
func (service Service) GetAlbumByID(ctx context.Context, id int64) (*cameraroll.Album, error) {
	alb := cameraroll.Album{}

	// find prepared statement
	stmt := service.preparedStmts[keyQueryGetAlbumByID]
	if stmt == nil {
		return nil, fmt.Errorf("GetAlbumByID [%d]: Cannot find prepared sql query", id)
	}

	// execute the query
	row := stmt.QueryRow(id)

	// parse response
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Description, &alb.CreatedAt, &alb.CoverID); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("GetAlbumByID[%d]: no such album", id)
		}

		return nil, fmt.Errorf("GetAlbumByID[%d]: %v", id, err)
	}

	// query database for album cover
	if alb.CoverID.Valid {
		img, _ := service.GetImageByID(ctx, alb.CoverID.Int64)
		if img != nil {
			alb.Cover = img
		}
	}

	return &alb, nil
}

// GetAlbums queries the database for certain amount of albums from a starting index
func (service Service) GetAlbums(ctx context.Context, start uint64, count uint64) ([]*cameraroll.Album, error) {
	// Album slice to hold the data from database query
	albums := []*cameraroll.Album{}

	// find prepared statement
	stmt := service.preparedStmts[keyQueryGetAlbums]
	if stmt == nil {
		return nil, fmt.Errorf("GetAlbums start[%d] count[%d]: Cannot find prepared sql query", start, count)
	}

	// execute the query
	rows, err := stmt.Query(start, count)

	// check if the query failed
	if err != nil {
		return nil, fmt.Errorf("GetAlbums start[%d] count[%d]: %v", start, count, err)
	}

	defer rows.Close()

	// parse response
	for rows.Next() {
		alb := cameraroll.Album{}
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Description, &alb.CreatedAt, &alb.CoverID); err != nil {
			return nil, fmt.Errorf("GetAlbums start[%d] count[%d]: %v", start, count, err)
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

// AddAlbum adds 1 album to the database,
// updating the album's ID upon success
func (service Service) AddAlbum(ctx context.Context, album *cameraroll.Album) error {
	if album == nil {
		return fmt.Errorf("AddAlbum : null pointer error")
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AddAlbum [%s]: %v", album.Title, err)
	}
	defer tx.Rollback()

	// execute the query
	var result sql.Result

	if album.CoverID.Valid {
		result, err = tx.ExecContext(ctx,
			`INSERT INTO albums (title, description, cover_id) 
			VALUES (?, ?, ?)`,
			album.Title,
			album.Description,
			album.CoverID)
	} else {
		result, err = tx.ExecContext(ctx,
			`INSERT INTO albums (title, description)
			VALUES (?, ?)`,
			album.Title,
			album.Description)
	}

	// check if the query failed
	if err != nil {
		return fmt.Errorf("AddAlbum [%s]: %v", album.Title, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("AddAlbum [%s]: %v", album.Title, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("AddAlbum [%s]: %v", album.Title, err)
	}

	album.ID = id

	return nil
}
