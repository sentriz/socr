package db

import (
	"context"
	_ "embed"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed schema.pgsql
var schema string

type DB struct {
	*pgxpool.Pool
	sq.StatementBuilderType
}

func New(dsn string) (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("create and connect pool: %w", err)
	}

	if _, err := pool.Exec(context.Background(), schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return &DB{
		Pool:                 pool,
		StatementBuilderType: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}, nil
}

func (db *DB) CreateScreenshot(screenshot *Screenshot) (*Screenshot, error) {
	q := db.
		Insert("screenshots").
		Columns("hash", "timestamp", "dim_width", "dim_height", "dominant_colour", "blurhash").
		Values(screenshot.Hash, screenshot.Timestamp, screenshot.DimWidth, screenshot.DimHeight, screenshot.DominantColour, screenshot.Blurhash).
		Suffix("returning *")

	sql, args, _ := q.ToSql()
	var result Screenshot
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

var sortFields = map[string]struct{}{
	"timestamp":  {},
	"similarity": {},
}
var sortOrders = map[string]struct{}{
	"asc":  {},
	"desc": {},
}

type SearchScreenshotsOptions struct {
	Body      string
	Directory string
	Limit     int
	Offset    int
	SortField string
	SortOrder string
}

func (db *DB) GetScreenshotByID(id int) (*Screenshot, error) {
	q := db.
		Select("screenshots.*").
		From("screenshots").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, _ := q.ToSql()
	var result Screenshot
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

func (db *DB) GetScreenshotByHash(hash string) (*Screenshot, error) {
	q := db.
		Select("*").
		From("screenshots").
		Where(sq.Eq{"hash": hash}).
		Limit(1)

	sql, args, _ := q.ToSql()
	var result Screenshot
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

func (db *DB) GetScreenshotByHashWithRelations(hash string) (*Screenshot, error) {
	colAggBlocks := db.
		Select("json_agg(blocks order by index)").
		From("blocks").
		Where("screenshot_id = screenshots.id")
	colAggAliases := db.
		Select("json_agg(dir_infos.directory_alias)").
		From("dir_infos").
		Where("screenshot_id = screenshots.id")

	q := db.
		Select("screenshots.*").
		Column(sq.Alias(colAggBlocks, "blocks")).
		Column(sq.Alias(colAggAliases, "directories")).
		From("screenshots").
		Where(sq.Eq{"hash": hash}).
		Limit(1)

	sql, args, _ := q.ToSql()
	var result Screenshot
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

func (db *DB) SearchScreenshots(options SearchScreenshotsOptions) ([]*Screenshot, error) {
	if _, ok := sortFields[options.SortField]; !ok {
		return nil, fmt.Errorf("invalid sort field %q provided", options.SortField)
	}
	if _, ok := sortOrders[options.SortOrder]; !ok {
		return nil, fmt.Errorf("invalid sort order %q provided", options.SortOrder)
	}

	q := db.
		Select("screenshots.*").
		From("screenshots").
		Limit(uint64(options.Limit)).
		Offset(uint64(options.Offset)).
		OrderBy(fmt.Sprintf("%s %s", options.SortField, options.SortOrder))
	if options.Directory != "" {
		q = q.
			Join("dir_infos on dir_infos.screenshot_id = screenshots.id").
			Where(sq.Eq{"dir_infos.directory_alias": options.Directory})
	}
	if options.Body != "" {
		colAggBlocks := sq.Expr("json_agg(blocks order by blocks.index)")
		colSimilarity := sq.Expr("avg(similarity(blocks.body, ?))", options.Body)
		q = q.
			Column(sq.Alias(colAggBlocks, "highlighted_blocks")).
			Column(sq.Alias(colSimilarity, "similarity")).
			LeftJoin("blocks on blocks.screenshot_id = screenshots.id").
			Where("blocks.body % ?", options.Body).
			GroupBy("screenshots.id")
	}

	sql, args, _ := q.ToSql()
	var results []*Screenshot
	return results, pgxscan.Select(context.Background(), db, &results, sql, args...)
}

func (db *DB) CreateBlocks(blocks []*Block) error {
	q := db.
		Insert("blocks").
		Columns("screenshot_id", "index", "min_x", "min_y", "max_x", "max_y", "body")
	for _, block := range blocks {
		q = q.Values(block.ScreenshotID, block.Index, block.MinX, block.MinY, block.MaxX, block.MaxY, block.Body)
	}

	sql, args, _ := q.ToSql()
	_, err := db.Exec(context.Background(), sql, args...)
	return err
}

func (db *DB) CreateDirInfo(dirInfo *DirInfo) (*DirInfo, error) {
	q := db.
		Insert("dir_infos").
		Columns("screenshot_id", "filename", "directory_alias").
		Values(dirInfo.ScreenshotID, dirInfo.Filename, dirInfo.DirectoryAlias).
		Suffix("returning *")

	sql, args, _ := q.ToSql()
	var result DirInfo
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

func (db *DB) GetDirInfo(directoryAlias string, filename string) (*DirInfo, error) {
	q := db.
		Select("*").
		From("dir_infos").
		Where(sq.Eq{
			"directory_alias": directoryAlias,
			"filename":        filename,
		}).
		Limit(1)

	sql, args, _ := q.ToSql()
	var result DirInfo
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

func (db *DB) GetDirInfoByScreenshotHash(hash string) (*DirInfo, error) {
	q := db.
		Select("dir_infos.*").
		From("dir_infos").
		Join("screenshots on screenshots.id = dir_infos.screenshot_id").
		Where(sq.Eq{"screenshots.hash": hash}).
		Limit(1)

	sql, args, _ := q.ToSql()
	var result DirInfo
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

func (db *DB) CountDirectories() ([]*DirectoryCount, error) {
	q := db.
		Select(
			"directory_alias",
			"count(1) as count",
		).
		From("dir_infos").
		GroupBy("directory_alias")

	sql, args, _ := q.ToSql()
	var result []*DirectoryCount
	return result, pgxscan.Select(context.Background(), db, &result, sql, args...)
}
