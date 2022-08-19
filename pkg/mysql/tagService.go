package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"chujungeng/camera-roll/pkg/cameraroll"
)

// DeleteTagByID removes a tag from the database
func (service Service) DeleteTagByID(ctx context.Context, id int64) error {
	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("DeleteTagByID [%d]: %v", id, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`DELETE FROM tags 
		WHERE id=?`,
		id)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("DeleteTagByID [%d]: %v", id, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteTagByID [%d]: %v", id, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("DeleteTagByID [%d]: %v", id, err)
	}

	return nil
}

// UpdateTagByID updates an tag's name
func (service Service) UpdateTagByID(ctx context.Context, id int64, newTag *cameraroll.Tag) error {
	if newTag == nil {
		return fmt.Errorf("UpdateTagByID [%d]: null pointer error", id)
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("UpdateTagByID [%d]: %v", id, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`UPDATE tags 
		SET name=? 
		WHERE id=?`,
		newTag.Name,
		id)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("UpdateTagByID [%d]: %v", id, err)
	}

	_, err = result.RowsAffected()
	if err != nil {
		return fmt.Errorf("UpdateTagByID [%d]: %v", id, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("UpdateTagByID [%d]: %v", id, err)
	}

	return nil
}

// GetTagByID returns the tag given its ID
func (service Service) GetTagByID(ctx context.Context, id int64) (*cameraroll.Tag, error) {
	tag := cameraroll.Tag{}

	// find prepared statement
	stmt := service.preparedStmts[keyQueryGetTagByID]
	if stmt == nil {
		return nil, fmt.Errorf("GetTagByID: Cannot find prepared sql query")
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("GetTagByID[%d]: %v", id, err)
	}
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, stmt)

	// execute the query
	row := txStmt.QueryRowContext(ctx, id)

	// parse response
	if err := row.Scan(&tag.ID, &tag.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("GetTagByID[%d]: no such tag", id)
		}

		return nil, fmt.Errorf("GetTagByID[%d]: %v", id, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("GetTagByID[%d]: %v", id, err)
	}

	return &tag, nil
}

// GetTags queries the database for all the tags
func (service Service) GetTags(ctx context.Context) ([]*cameraroll.Tag, error) {
	// tag slice to hold the data from database query
	tags := []*cameraroll.Tag{}

	// find prepared statement
	stmt := service.preparedStmts[keyQueryGetTags]
	if stmt == nil {
		return nil, fmt.Errorf("GetTags: Cannot find prepared sql query")
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("GetTags: %v", err)
	}
	defer tx.Rollback()

	txStmt := tx.StmtContext(ctx, stmt)

	// execute the query
	rows, err := txStmt.QueryContext(ctx)

	// check if the query failed
	if err != nil {
		return nil, fmt.Errorf("GetTags: %v", err)
	}

	defer rows.Close()

	// parse response
	for rows.Next() {
		tag := cameraroll.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name); err != nil {
			return nil, fmt.Errorf("GetTags: %v", err)
		}

		// add tag to the return slice
		tags = append(tags, &tag)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("GetTags: %v", err)
	}

	return tags, nil
}

// AddTag adds 1 tag to the database,
// updating the tag's ID upon success
func (service Service) AddTag(ctx context.Context, tag *cameraroll.Tag) error {
	if tag == nil {
		return fmt.Errorf("AddTag : null pointer error")
	}

	// start a transaction
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AddTag [%s]: %v", tag.Name, err)
	}
	defer tx.Rollback()

	// execute the query
	result, err := tx.ExecContext(ctx,
		`INSERT INTO tags (name) 
		VALUES (?)`,
		tag.Name)

	// check if the query failed
	if err != nil {
		return fmt.Errorf("AddTag [%s]: %v", tag.Name, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("AddTag [%s]: %v", tag.Name, err)
	}

	// commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("AddTag [%s]: %v", tag.Name, err)
	}

	tag.ID = id

	return nil
}
