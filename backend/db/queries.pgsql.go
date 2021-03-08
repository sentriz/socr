// Code generated by pggen. DO NOT EDIT.

package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"time"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	GetDirInfo(ctx context.Context, directoryAlias string, filename string) (int, error)
	// GetDirInfoBatch enqueues a GetDirInfo query into batch to be executed
	// later by the batch.
	GetDirInfoBatch(batch *pgx.Batch, directoryAlias string, filename string)
	// GetDirInfoScan scans the result of an executed GetDirInfoBatch query.
	GetDirInfoScan(results pgx.BatchResults) (int, error)

	GetScreenshotByID(ctx context.Context, id int) (GetScreenshotByIDRow, error)
	// GetScreenshotByIDBatch enqueues a GetScreenshotByID query into batch to be executed
	// later by the batch.
	GetScreenshotByIDBatch(batch *pgx.Batch, id int)
	// GetScreenshotByIDScan scans the result of an executed GetScreenshotByIDBatch query.
	GetScreenshotByIDScan(results pgx.BatchResults) (GetScreenshotByIDRow, error)

	GetScreenshotByHash(ctx context.Context, hash string) (GetScreenshotByHashRow, error)
	// GetScreenshotByHashBatch enqueues a GetScreenshotByHash query into batch to be executed
	// later by the batch.
	GetScreenshotByHashBatch(batch *pgx.Batch, hash string)
	// GetScreenshotByHashScan scans the result of an executed GetScreenshotByHashBatch query.
	GetScreenshotByHashScan(results pgx.BatchResults) (GetScreenshotByHashRow, error)

	CreateScreenshot(ctx context.Context, params CreateScreenshotParams) (CreateScreenshotRow, error)
	// CreateScreenshotBatch enqueues a CreateScreenshot query into batch to be executed
	// later by the batch.
	CreateScreenshotBatch(batch *pgx.Batch, params CreateScreenshotParams)
	// CreateScreenshotScan scans the result of an executed CreateScreenshotBatch query.
	CreateScreenshotScan(results pgx.BatchResults) (CreateScreenshotRow, error)

	GetAllScreenshots(ctx context.Context) ([]GetAllScreenshotsRow, error)
	// GetAllScreenshotsBatch enqueues a GetAllScreenshots query into batch to be executed
	// later by the batch.
	GetAllScreenshotsBatch(batch *pgx.Batch)
	// GetAllScreenshotsScan scans the result of an executed GetAllScreenshotsBatch query.
	GetAllScreenshotsScan(results pgx.BatchResults) ([]GetAllScreenshotsRow, error)

	CreateBlock(ctx context.Context, params CreateBlockParams) (pgconn.CommandTag, error)
	// CreateBlockBatch enqueues a CreateBlock query into batch to be executed
	// later by the batch.
	CreateBlockBatch(batch *pgx.Batch, params CreateBlockParams)
	// CreateBlockScan scans the result of an executed CreateBlockBatch query.
	CreateBlockScan(results pgx.BatchResults) (pgconn.CommandTag, error)

	CountDirectoriesByAlias(ctx context.Context) ([]CountDirectoriesByAliasRow, error)
	// CountDirectoriesByAliasBatch enqueues a CountDirectoriesByAlias query into batch to be executed
	// later by the batch.
	CountDirectoriesByAliasBatch(batch *pgx.Batch)
	// CountDirectoriesByAliasScan scans the result of an executed CountDirectoriesByAliasBatch query.
	CountDirectoriesByAliasScan(results pgx.BatchResults) ([]CountDirectoriesByAliasRow, error)

	// https://www.postgresql.org/docs/current/pgtrgm.html
	SearchScreenshots(ctx context.Context, params SearchScreenshotsParams) ([]SearchScreenshotsRow, error)
	// SearchScreenshotsBatch enqueues a SearchScreenshots query into batch to be executed
	// later by the batch.
	SearchScreenshotsBatch(batch *pgx.Batch, params SearchScreenshotsParams)
	// SearchScreenshotsScan scans the result of an executed SearchScreenshotsBatch query.
	SearchScreenshotsScan(results pgx.BatchResults) ([]SearchScreenshotsRow, error)

	CreateDirInfo(ctx context.Context, params CreateDirInfoParams) (pgconn.CommandTag, error)
	// CreateDirInfoBatch enqueues a CreateDirInfo query into batch to be executed
	// later by the batch.
	CreateDirInfoBatch(batch *pgx.Batch, params CreateDirInfoParams)
	// CreateDirInfoScan scans the result of an executed CreateDirInfoBatch query.
	CreateDirInfoScan(results pgx.BatchResults) (pgconn.CommandTag, error)

	GetScreenshotPathByHash(ctx context.Context, hash string) (GetScreenshotPathByHashRow, error)
	// GetScreenshotPathByHashBatch enqueues a GetScreenshotPathByHash query into batch to be executed
	// later by the batch.
	GetScreenshotPathByHashBatch(batch *pgx.Batch, hash string)
	// GetScreenshotPathByHashScan scans the result of an executed GetScreenshotPathByHashBatch query.
	GetScreenshotPathByHashScan(results pgx.BatchResults) (GetScreenshotPathByHashRow, error)
}

type DBQuerier struct {
	conn genericConn
}

var _ Querier = &DBQuerier{}

// genericConn is a connection to a Postgres database. This is usually backed by
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	// Query executes sql with args. If there is an error the returned Rows will
	// be returned in an error state. So it is allowed to ignore the error
	// returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with pgx.ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes sql. sql can be either a prepared statement name or an SQL
	// string. arguments should be referenced positionally from the sql string
	// as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{
		conn: conn,
	}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

// ignoredOID means we don't know or care about the OID for a type. This is okay
// because pgx only uses the OID to encode values and lookup a decoder. We only
// use ignoredOID for decoding and we always specify a concrete decoder for scan
// methods.
const ignoredOID = 0

// Blocks represents the Postgres composite type "blocks".
type Blocks struct {
	ID           int    `json:"id"`
	ScreenshotID int    `json:"screenshot_id"`
	Index        int    `json:"index"`
	MinX         int    `json:"min_x"`
	MinY         int    `json:"min_y"`
	MaxX         int    `json:"max_x"`
	MaxY         int    `json:"max_y"`
	Body         string `json:"body"`
}

const getDirInfoSQL = `select
    1
from
    dir_infos
where
    directory_alias = $1
    and filename = $2
limit 1;`

// GetDirInfo implements Querier.GetDirInfo.
func (q *DBQuerier) GetDirInfo(ctx context.Context, directoryAlias string, filename string) (int, error) {
	row := q.conn.QueryRow(ctx, getDirInfoSQL, directoryAlias, filename)
	var item int
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query GetDirInfo: %w", err)
	}
	return item, nil
}

// GetDirInfoBatch implements Querier.GetDirInfoBatch.
func (q *DBQuerier) GetDirInfoBatch(batch *pgx.Batch, directoryAlias string, filename string) {
	batch.Queue(getDirInfoSQL, directoryAlias, filename)
}

// GetDirInfoScan implements Querier.GetDirInfoScan.
func (q *DBQuerier) GetDirInfoScan(results pgx.BatchResults) (int, error) {
	row := results.QueryRow()
	var item int
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan GetDirInfoBatch row: %w", err)
	}
	return item, nil
}

const getScreenshotByIDSQL = `select
    *
from
    screenshots
where
    id = $1
limit 1;`

type GetScreenshotByIDRow struct {
	ID             int       `json:"id"`
	Hash           string    `json:"hash"`
	Timestamp      time.Time `json:"timestamp"`
	DimWidth       int       `json:"dim_width"`
	DimHeight      int       `json:"dim_height"`
	DominantColour string    `json:"dominant_colour"`
	Blurhash       string    `json:"blurhash"`
}

// GetScreenshotByID implements Querier.GetScreenshotByID.
func (q *DBQuerier) GetScreenshotByID(ctx context.Context, id int) (GetScreenshotByIDRow, error) {
	row := q.conn.QueryRow(ctx, getScreenshotByIDSQL, id)
	var item GetScreenshotByIDRow
	if err := row.Scan(&item.ID, &item.Hash, &item.Timestamp, &item.DimWidth, &item.DimHeight, &item.DominantColour, &item.Blurhash); err != nil {
		return item, fmt.Errorf("query GetScreenshotByID: %w", err)
	}
	return item, nil
}

// GetScreenshotByIDBatch implements Querier.GetScreenshotByIDBatch.
func (q *DBQuerier) GetScreenshotByIDBatch(batch *pgx.Batch, id int) {
	batch.Queue(getScreenshotByIDSQL, id)
}

// GetScreenshotByIDScan implements Querier.GetScreenshotByIDScan.
func (q *DBQuerier) GetScreenshotByIDScan(results pgx.BatchResults) (GetScreenshotByIDRow, error) {
	row := results.QueryRow()
	var item GetScreenshotByIDRow
	if err := row.Scan(&item.ID, &item.Hash, &item.Timestamp, &item.DimWidth, &item.DimHeight, &item.DominantColour, &item.Blurhash); err != nil {
		return item, fmt.Errorf("scan GetScreenshotByIDBatch row: %w", err)
	}
	return item, nil
}

const getScreenshotByHashSQL = `select
    *
from
    screenshots
where
    hash = $1
limit 1;`

type GetScreenshotByHashRow struct {
	ID             int       `json:"id"`
	Hash           string    `json:"hash"`
	Timestamp      time.Time `json:"timestamp"`
	DimWidth       int       `json:"dim_width"`
	DimHeight      int       `json:"dim_height"`
	DominantColour string    `json:"dominant_colour"`
	Blurhash       string    `json:"blurhash"`
}

// GetScreenshotByHash implements Querier.GetScreenshotByHash.
func (q *DBQuerier) GetScreenshotByHash(ctx context.Context, hash string) (GetScreenshotByHashRow, error) {
	row := q.conn.QueryRow(ctx, getScreenshotByHashSQL, hash)
	var item GetScreenshotByHashRow
	if err := row.Scan(&item.ID, &item.Hash, &item.Timestamp, &item.DimWidth, &item.DimHeight, &item.DominantColour, &item.Blurhash); err != nil {
		return item, fmt.Errorf("query GetScreenshotByHash: %w", err)
	}
	return item, nil
}

// GetScreenshotByHashBatch implements Querier.GetScreenshotByHashBatch.
func (q *DBQuerier) GetScreenshotByHashBatch(batch *pgx.Batch, hash string) {
	batch.Queue(getScreenshotByHashSQL, hash)
}

// GetScreenshotByHashScan implements Querier.GetScreenshotByHashScan.
func (q *DBQuerier) GetScreenshotByHashScan(results pgx.BatchResults) (GetScreenshotByHashRow, error) {
	row := results.QueryRow()
	var item GetScreenshotByHashRow
	if err := row.Scan(&item.ID, &item.Hash, &item.Timestamp, &item.DimWidth, &item.DimHeight, &item.DominantColour, &item.Blurhash); err != nil {
		return item, fmt.Errorf("scan GetScreenshotByHashBatch row: %w", err)
	}
	return item, nil
}

const createScreenshotSQL = `insert into screenshots (hash, timestamp, dim_width, dim_height, dominant_colour, blurhash)
    values ($1, $2, $3, $4, $5, $6)
returning
    *;`

type CreateScreenshotParams struct {
	Hash           string
	Timestamp      time.Time
	DimWidth       int
	DimHeight      int
	DominantColour string
	Blurhash       string
}

type CreateScreenshotRow struct {
	ID             int       `json:"id"`
	Hash           string    `json:"hash"`
	Timestamp      time.Time `json:"timestamp"`
	DimWidth       int       `json:"dim_width"`
	DimHeight      int       `json:"dim_height"`
	DominantColour string    `json:"dominant_colour"`
	Blurhash       string    `json:"blurhash"`
}

// CreateScreenshot implements Querier.CreateScreenshot.
func (q *DBQuerier) CreateScreenshot(ctx context.Context, params CreateScreenshotParams) (CreateScreenshotRow, error) {
	row := q.conn.QueryRow(ctx, createScreenshotSQL, params.Hash, params.Timestamp, params.DimWidth, params.DimHeight, params.DominantColour, params.Blurhash)
	var item CreateScreenshotRow
	if err := row.Scan(&item.ID, &item.Hash, &item.Timestamp, &item.DimWidth, &item.DimHeight, &item.DominantColour, &item.Blurhash); err != nil {
		return item, fmt.Errorf("query CreateScreenshot: %w", err)
	}
	return item, nil
}

// CreateScreenshotBatch implements Querier.CreateScreenshotBatch.
func (q *DBQuerier) CreateScreenshotBatch(batch *pgx.Batch, params CreateScreenshotParams) {
	batch.Queue(createScreenshotSQL, params.Hash, params.Timestamp, params.DimWidth, params.DimHeight, params.DominantColour, params.Blurhash)
}

// CreateScreenshotScan implements Querier.CreateScreenshotScan.
func (q *DBQuerier) CreateScreenshotScan(results pgx.BatchResults) (CreateScreenshotRow, error) {
	row := results.QueryRow()
	var item CreateScreenshotRow
	if err := row.Scan(&item.ID, &item.Hash, &item.Timestamp, &item.DimWidth, &item.DimHeight, &item.DominantColour, &item.Blurhash); err != nil {
		return item, fmt.Errorf("scan CreateScreenshotBatch row: %w", err)
	}
	return item, nil
}

const getAllScreenshotsSQL = `select
    *
from
    screenshots;`

type GetAllScreenshotsRow struct {
	ID             int       `json:"id"`
	Hash           string    `json:"hash"`
	Timestamp      time.Time `json:"timestamp"`
	DimWidth       int       `json:"dim_width"`
	DimHeight      int       `json:"dim_height"`
	DominantColour string    `json:"dominant_colour"`
	Blurhash       string    `json:"blurhash"`
}

// GetAllScreenshots implements Querier.GetAllScreenshots.
func (q *DBQuerier) GetAllScreenshots(ctx context.Context) ([]GetAllScreenshotsRow, error) {
	rows, err := q.conn.Query(ctx, getAllScreenshotsSQL)
	if err != nil {
		return nil, fmt.Errorf("query GetAllScreenshots: %w", err)
	}
	defer rows.Close()
	items := []GetAllScreenshotsRow{}
	for rows.Next() {
		var item GetAllScreenshotsRow
		if err := rows.Scan(&item.ID, &item.Hash, &item.Timestamp, &item.DimWidth, &item.DimHeight, &item.DominantColour, &item.Blurhash); err != nil {
			return nil, fmt.Errorf("scan GetAllScreenshots row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close GetAllScreenshots rows: %w", err)
	}
	return items, err
}

// GetAllScreenshotsBatch implements Querier.GetAllScreenshotsBatch.
func (q *DBQuerier) GetAllScreenshotsBatch(batch *pgx.Batch) {
	batch.Queue(getAllScreenshotsSQL)
}

// GetAllScreenshotsScan implements Querier.GetAllScreenshotsScan.
func (q *DBQuerier) GetAllScreenshotsScan(results pgx.BatchResults) ([]GetAllScreenshotsRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query GetAllScreenshotsBatch: %w", err)
	}
	defer rows.Close()
	items := []GetAllScreenshotsRow{}
	for rows.Next() {
		var item GetAllScreenshotsRow
		if err := rows.Scan(&item.ID, &item.Hash, &item.Timestamp, &item.DimWidth, &item.DimHeight, &item.DominantColour, &item.Blurhash); err != nil {
			return nil, fmt.Errorf("scan GetAllScreenshotsBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close GetAllScreenshotsBatch rows: %w", err)
	}
	return items, err
}

const createBlockSQL = `insert into blocks (screenshot_id, index, min_x, min_y, max_x, max_y, body)
        values ($1, $2, $3, $4, $5, $6, $7);`

type CreateBlockParams struct {
	ScreenshotID int
	Index        int
	MinX         int
	MinY         int
	MaxX         int
	MaxY         int
	Body         string
}

// CreateBlock implements Querier.CreateBlock.
func (q *DBQuerier) CreateBlock(ctx context.Context, params CreateBlockParams) (pgconn.CommandTag, error) {
	cmdTag, err := q.conn.Exec(ctx, createBlockSQL, params.ScreenshotID, params.Index, params.MinX, params.MinY, params.MaxX, params.MaxY, params.Body)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query CreateBlock: %w", err)
	}
	return cmdTag, err
}

// CreateBlockBatch implements Querier.CreateBlockBatch.
func (q *DBQuerier) CreateBlockBatch(batch *pgx.Batch, params CreateBlockParams) {
	batch.Queue(createBlockSQL, params.ScreenshotID, params.Index, params.MinX, params.MinY, params.MaxX, params.MaxY, params.Body)
}

// CreateBlockScan implements Querier.CreateBlockScan.
func (q *DBQuerier) CreateBlockScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec CreateBlockBatch: %w", err)
	}
	return cmdTag, err
}

const countDirectoriesByAliasSQL = `select
    directory_alias,
    count(1)
from
    dir_infos
group by
    directory_alias;`

type CountDirectoriesByAliasRow struct {
	DirectoryAlias string `json:"directory_alias"`
	Count          int    `json:"count"`
}

// CountDirectoriesByAlias implements Querier.CountDirectoriesByAlias.
func (q *DBQuerier) CountDirectoriesByAlias(ctx context.Context) ([]CountDirectoriesByAliasRow, error) {
	rows, err := q.conn.Query(ctx, countDirectoriesByAliasSQL)
	if err != nil {
		return nil, fmt.Errorf("query CountDirectoriesByAlias: %w", err)
	}
	defer rows.Close()
	items := []CountDirectoriesByAliasRow{}
	for rows.Next() {
		var item CountDirectoriesByAliasRow
		if err := rows.Scan(&item.DirectoryAlias, &item.Count); err != nil {
			return nil, fmt.Errorf("scan CountDirectoriesByAlias row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close CountDirectoriesByAlias rows: %w", err)
	}
	return items, err
}

// CountDirectoriesByAliasBatch implements Querier.CountDirectoriesByAliasBatch.
func (q *DBQuerier) CountDirectoriesByAliasBatch(batch *pgx.Batch) {
	batch.Queue(countDirectoriesByAliasSQL)
}

// CountDirectoriesByAliasScan implements Querier.CountDirectoriesByAliasScan.
func (q *DBQuerier) CountDirectoriesByAliasScan(results pgx.BatchResults) ([]CountDirectoriesByAliasRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query CountDirectoriesByAliasBatch: %w", err)
	}
	defer rows.Close()
	items := []CountDirectoriesByAliasRow{}
	for rows.Next() {
		var item CountDirectoriesByAliasRow
		if err := rows.Scan(&item.DirectoryAlias, &item.Count); err != nil {
			return nil, fmt.Errorf("scan CountDirectoriesByAliasBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close CountDirectoriesByAliasBatch rows: %w", err)
	}
	return items, err
}

const searchScreenshotsSQL = `select
    screenshots.*,
    array_agg(blocks order by blocks.index) as blocks,
    avg(similarity (blocks.body, $1)) as similarity,
    count(1) over () as total
from
    screenshots
    join blocks on blocks.screenshot_id = screenshots.id
where
    blocks.body % $1
group by
    screenshots.id
order by
    similarity desc
limit $2 offset $3;`

type SearchScreenshotsParams struct {
	Body   string
	Limit  int
	Offset int
}

type SearchScreenshotsRow struct {
	ID             int       `json:"id"`
	Hash           string    `json:"hash"`
	Timestamp      time.Time `json:"timestamp"`
	DimWidth       int       `json:"dim_width"`
	DimHeight      int       `json:"dim_height"`
	DominantColour string    `json:"dominant_colour"`
	Blurhash       string    `json:"blurhash"`
	Blocks         []Blocks  `json:"blocks"`
	Similarity     float64   `json:"similarity"`
	Total          int       `json:"total"`
}

// SearchScreenshots implements Querier.SearchScreenshots.
func (q *DBQuerier) SearchScreenshots(ctx context.Context, params SearchScreenshotsParams) ([]SearchScreenshotsRow, error) {
	rows, err := q.conn.Query(ctx, searchScreenshotsSQL, params.Body, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query SearchScreenshots: %w", err)
	}
	defer rows.Close()
	items := []SearchScreenshotsRow{}
	blocksRow, _ := pgtype.NewCompositeTypeValues("blocks", []pgtype.CompositeTypeField{
		{Name: "ID", OID: ignoredOID},
		{Name: "ScreenshotID", OID: ignoredOID},
		{Name: "Index", OID: ignoredOID},
		{Name: "MinX", OID: ignoredOID},
		{Name: "MinY", OID: ignoredOID},
		{Name: "MaxX", OID: ignoredOID},
		{Name: "MaxY", OID: ignoredOID},
		{Name: "Body", OID: ignoredOID},
	}, []pgtype.ValueTranscoder{
		&pgtype.Int4{},
		&pgtype.Int4{},
		&pgtype.Int2{},
		&pgtype.Int2{},
		&pgtype.Int2{},
		&pgtype.Int2{},
		&pgtype.Int2{},
		&pgtype.Text{},
	})
	blocksArray := pgtype.NewArrayType("_blocks", ignoredOID, func() pgtype.ValueTranscoder {
		return blocksRow.NewTypeValue().(*pgtype.CompositeType)
	})
	for rows.Next() {
		var item SearchScreenshotsRow
		if err := rows.Scan(&item.ID, &item.Hash, &item.Timestamp, &item.DimWidth, &item.DimHeight, &item.DominantColour, &item.Blurhash, blocksArray, &item.Similarity, &item.Total); err != nil {
			return nil, fmt.Errorf("scan SearchScreenshots row: %w", err)
		}
		blocksArray.AssignTo(&item.Blocks)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close SearchScreenshots rows: %w", err)
	}
	return items, err
}

// SearchScreenshotsBatch implements Querier.SearchScreenshotsBatch.
func (q *DBQuerier) SearchScreenshotsBatch(batch *pgx.Batch, params SearchScreenshotsParams) {
	batch.Queue(searchScreenshotsSQL, params.Body, params.Limit, params.Offset)
}

// SearchScreenshotsScan implements Querier.SearchScreenshotsScan.
func (q *DBQuerier) SearchScreenshotsScan(results pgx.BatchResults) ([]SearchScreenshotsRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query SearchScreenshotsBatch: %w", err)
	}
	defer rows.Close()
	items := []SearchScreenshotsRow{}
	blocksRow, _ := pgtype.NewCompositeTypeValues("blocks", []pgtype.CompositeTypeField{
		{Name: "ID", OID: ignoredOID},
		{Name: "ScreenshotID", OID: ignoredOID},
		{Name: "Index", OID: ignoredOID},
		{Name: "MinX", OID: ignoredOID},
		{Name: "MinY", OID: ignoredOID},
		{Name: "MaxX", OID: ignoredOID},
		{Name: "MaxY", OID: ignoredOID},
		{Name: "Body", OID: ignoredOID},
	}, []pgtype.ValueTranscoder{
		&pgtype.Int4{},
		&pgtype.Int4{},
		&pgtype.Int2{},
		&pgtype.Int2{},
		&pgtype.Int2{},
		&pgtype.Int2{},
		&pgtype.Int2{},
		&pgtype.Text{},
	})
	blocksArray := pgtype.NewArrayType("_blocks", ignoredOID, func() pgtype.ValueTranscoder {
		return blocksRow.NewTypeValue().(*pgtype.CompositeType)
	})
	for rows.Next() {
		var item SearchScreenshotsRow
		if err := rows.Scan(&item.ID, &item.Hash, &item.Timestamp, &item.DimWidth, &item.DimHeight, &item.DominantColour, &item.Blurhash, blocksArray, &item.Similarity, &item.Total); err != nil {
			return nil, fmt.Errorf("scan SearchScreenshotsBatch row: %w", err)
		}
		blocksArray.AssignTo(&item.Blocks)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close SearchScreenshotsBatch rows: %w", err)
	}
	return items, err
}

const createDirInfoSQL = `insert into dir_infos (screenshot_id, filename, directory_alias)
    values ($1, $2, $3)
on conflict
    do nothing;`

type CreateDirInfoParams struct {
	ScreenshotID   int
	Filename       string
	DirectoryAlias string
}

// CreateDirInfo implements Querier.CreateDirInfo.
func (q *DBQuerier) CreateDirInfo(ctx context.Context, params CreateDirInfoParams) (pgconn.CommandTag, error) {
	cmdTag, err := q.conn.Exec(ctx, createDirInfoSQL, params.ScreenshotID, params.Filename, params.DirectoryAlias)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query CreateDirInfo: %w", err)
	}
	return cmdTag, err
}

// CreateDirInfoBatch implements Querier.CreateDirInfoBatch.
func (q *DBQuerier) CreateDirInfoBatch(batch *pgx.Batch, params CreateDirInfoParams) {
	batch.Queue(createDirInfoSQL, params.ScreenshotID, params.Filename, params.DirectoryAlias)
}

// CreateDirInfoScan implements Querier.CreateDirInfoScan.
func (q *DBQuerier) CreateDirInfoScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec CreateDirInfoBatch: %w", err)
	}
	return cmdTag, err
}

const getScreenshotPathByHashSQL = `select
    dir_infos.filename,
    dir_infos.directory_alias
from
    dir_infos
    join screenshots on screenshots.id = dir_infos.screenshot_id
where
    screenshots.hash = $1
limit 1;`

type GetScreenshotPathByHashRow struct {
	Filename       string `json:"filename"`
	DirectoryAlias string `json:"directory_alias"`
}

// GetScreenshotPathByHash implements Querier.GetScreenshotPathByHash.
func (q *DBQuerier) GetScreenshotPathByHash(ctx context.Context, hash string) (GetScreenshotPathByHashRow, error) {
	row := q.conn.QueryRow(ctx, getScreenshotPathByHashSQL, hash)
	var item GetScreenshotPathByHashRow
	if err := row.Scan(&item.Filename, &item.DirectoryAlias); err != nil {
		return item, fmt.Errorf("query GetScreenshotPathByHash: %w", err)
	}
	return item, nil
}

// GetScreenshotPathByHashBatch implements Querier.GetScreenshotPathByHashBatch.
func (q *DBQuerier) GetScreenshotPathByHashBatch(batch *pgx.Batch, hash string) {
	batch.Queue(getScreenshotPathByHashSQL, hash)
}

// GetScreenshotPathByHashScan implements Querier.GetScreenshotPathByHashScan.
func (q *DBQuerier) GetScreenshotPathByHashScan(results pgx.BatchResults) (GetScreenshotPathByHashRow, error) {
	row := results.QueryRow()
	var item GetScreenshotPathByHashRow
	if err := row.Scan(&item.Filename, &item.DirectoryAlias); err != nil {
		return item, fmt.Errorf("scan GetScreenshotPathByHashBatch row: %w", err)
	}
	return item, nil
}
