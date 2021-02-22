// Code generated by sqlc. DO NOT EDIT.
// source: queries.pgsql

package db

import (
	"context"
	"time"
)

const createBlock = `-- name: CreateBlock :exec
INSERT INTO
    blocks (min_x, min_y, max_x, max_y, body)
VALUES
    ($1, $2, $3, $4, $5)
`

type CreateBlockParams struct {
	MinX int32  `json:"min_x"`
	MinY int32  `json:"min_y"`
	MaxX int32  `json:"max_x"`
	MaxY int32  `json:"max_y"`
	Body string `json:"body"`
}

func (q *Queries) CreateBlock(ctx context.Context, arg CreateBlockParams) error {
	_, err := q.exec(ctx, q.createBlockStmt, createBlock,
		arg.MinX,
		arg.MinY,
		arg.MaxX,
		arg.MaxY,
		arg.Body,
	)
	return err
}

const createScreenshot = `-- name: CreateScreenshot :exec
INSERT INTO
    screenshots (id, directory, filename, stamp)
VALUES
    ($1, $2, $3, $4)
`

type CreateScreenshotParams struct {
	ID        int64     `json:"id"`
	Directory int32     `json:"directory"`
	Filename  string    `json:"filename"`
	Stamp     time.Time `json:"stamp"`
}

func (q *Queries) CreateScreenshot(ctx context.Context, arg CreateScreenshotParams) error {
	_, err := q.exec(ctx, q.createScreenshotStmt, createScreenshot,
		arg.ID,
		arg.Directory,
		arg.Filename,
		arg.Stamp,
	)
	return err
}

const getAllScreenshots = `-- name: GetAllScreenshots :many
SELECT
    id, stamp, directory, filename, filetype
FROM
    screenshots
`

func (q *Queries) GetAllScreenshots(ctx context.Context) ([]Screenshot, error) {
	rows, err := q.query(ctx, q.getAllScreenshotsStmt, getAllScreenshots)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Screenshot
	for rows.Next() {
		var i Screenshot
		if err := rows.Scan(
			&i.ID,
			&i.Stamp,
			&i.Directory,
			&i.Filename,
			&i.Filetype,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getScreenshotByID = `-- name: GetScreenshotByID :one
SELECT
    id, stamp, directory, filename, filetype
FROM
    screenshots
WHERE
    id = $1
LIMIT
    1
`

func (q *Queries) GetScreenshotByID(ctx context.Context, id int64) (Screenshot, error) {
	row := q.queryRow(ctx, q.getScreenshotByIDStmt, getScreenshotByID, id)
	var i Screenshot
	err := row.Scan(
		&i.ID,
		&i.Stamp,
		&i.Directory,
		&i.Filename,
		&i.Filetype,
	)
	return i, err
}

const searchBlock = `-- name: SearchBlock :many
SELECT
    id, stamp, directory, filename, filetype
FROM
    screenshots
WHERE
    $1 :: TEXT % body
LIMIT
    40
`

func (q *Queries) SearchBlock(ctx context.Context, body string) ([]Screenshot, error) {
	rows, err := q.query(ctx, q.searchBlockStmt, searchBlock, body)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Screenshot
	for rows.Next() {
		var i Screenshot
		if err := rows.Scan(
			&i.ID,
			&i.Stamp,
			&i.Directory,
			&i.Filename,
			&i.Filetype,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
