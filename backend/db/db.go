package db

import (
	"context"
	_ "embed" //nolint:golint
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

//nolint:gochecknoglobals
//go:embed schema.sql
var schema string

type DB struct {
	*pgxpool.Pool
	sq.StatementBuilderType
}

func New(dsn string) (*DB, error) {
	pool, err := waitConnect(context.Background(), dsn, 500*time.Millisecond, 10)
	if err != nil {
		return nil, fmt.Errorf("create and connect pool: %w", err)
	}

	log.Println("executing schema")
	if _, err := pool.Exec(context.Background(), schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	return &DB{
		Pool:                 pool,
		StatementBuilderType: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}, nil
}

func waitConnect(ctx context.Context, dsn string, interval time.Duration, times int) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error
	for i := 0; i < times; i++ {
		if pool, err = pgxpool.Connect(ctx, dsn); err == nil {
			return pool, nil
		}
		time.Sleep(interval)
	}
	return nil, fmt.Errorf("failed after %d tries: %w", times, err)
}

func (db *DB) CreateMedia(media *Media) (*Media, error) {
	q := db.
		Insert("medias").
		Columns("hash", "type", "timestamp", "dim_width", "dim_height", "dominant_colour", "blurhash").
		Values(media.Hash, media.Type, media.Timestamp, media.DimWidth, media.DimHeight, media.DominantColour, media.Blurhash).
		Suffix("returning *")

	sql, args, _ := q.ToSql()
	var result Media
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

type SearchMediasOptions struct {
	Body      string
	Directory string
	Limit     int
	Offset    int
	SortField string
	SortOrder string
}

func (db *DB) GetMediaByID(id int) (*Media, error) {
	q := db.
		Select("medias.*").
		From("medias").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, _ := q.ToSql()
	var result Media
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

func (db *DB) GetMediaByHash(hash string) (*Media, error) {
	q := db.
		Select("*").
		From("medias").
		Where(sq.Eq{"hash": hash}).
		Limit(1)

	sql, args, _ := q.ToSql()
	var result Media
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

func (db *DB) GetMediaByHashWithRelations(hash string) (*Media, error) {
	colAggBlocks := db.
		Select("json_agg(blocks order by index)").
		From("blocks").
		Where("media_id = medias.id")
	colAggAliases := db.
		Select("json_agg(distinct dir_infos.directory_alias)").
		From("dir_infos").
		Where("media_id = medias.id")

	q := db.
		Select("medias.*").
		Column(sq.Alias(colAggBlocks, "blocks")).
		Column(sq.Alias(colAggAliases, "directories")).
		From("medias").
		Where(sq.Eq{"hash": hash}).
		Limit(1)

	sql, args, _ := q.ToSql()
	var result Media
	return &result, pgxscan.Get(context.Background(), db, &result, sql, args...)
}

func (db *DB) SearchMedias(options SearchMediasOptions) ([]*Media, error) {
	if !isSortField(options.SortField) {
		return nil, fmt.Errorf("invalid sort field %q provided", options.SortField)
	}
	if !isSortOrder(options.SortOrder) {
		return nil, fmt.Errorf("invalid sort order %q provided", options.SortOrder)
	}

	q := db.
		Select("medias.*").
		From("medias").
		Limit(uint64(options.Limit)).
		Offset(uint64(options.Offset)).
		OrderBy(fmt.Sprintf("%s %s", options.SortField, options.SortOrder))
	if options.Directory != "" {
		q = q.
			Join("dir_infos on dir_infos.media_id = medias.id").
			Where(sq.Eq{"dir_infos.directory_alias": options.Directory})
	}
	if options.Body != "" {
		colAggBlocks := sq.Expr("json_agg(blocks order by blocks.index)")
		colSimilarity := sq.Expr("avg(similarity(blocks.body, ?))", options.Body)
		q = q.
			Column(sq.Alias(colAggBlocks, "highlighted_blocks")).
			Column(sq.Alias(colSimilarity, "similarity")).
			LeftJoin("blocks on blocks.media_id = medias.id").
			Where("blocks.body %> ?", options.Body).
			GroupBy("medias.id")
	}

	sql, args, _ := q.ToSql()
	var results []*Media
	return results, pgxscan.Select(context.Background(), db, &results, sql, args...)
}

func (db *DB) CreateBlocks(blocks []*Block) error {
	q := db.
		Insert("blocks").
		Columns("media_id", "index", "min_x", "min_y", "max_x", "max_y", "body")
	for _, block := range blocks {
		q = q.Values(block.MediaID, block.Index, block.MinX, block.MinY, block.MaxX, block.MaxY, block.Body)
	}

	sql, args, _ := q.ToSql()
	_, err := db.Exec(context.Background(), sql, args...)
	return err
}

func (db *DB) CreateDirInfo(dirInfo *DirInfo) (*DirInfo, error) {
	q := db.
		Insert("dir_infos").
		Columns("media_id", "filename", "directory_alias").
		Values(dirInfo.MediaID, dirInfo.Filename, dirInfo.DirectoryAlias).
		Suffix("on conflict do nothing").
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

func (db *DB) GetDirInfoByMediaHash(hash string) (*DirInfo, error) {
	q := db.
		Select("dir_infos.*").
		From("dir_infos").
		Join("medias on medias.id = dir_infos.media_id").
		Where(sq.Eq{"medias.hash": hash}).
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

func isSortField(f string) bool {
	fields := map[string]struct{}{
		"timestamp":  {},
		"similarity": {},
	}
	_, ok := fields[f]
	return ok
}

func isSortOrder(f string) bool {
	fields := map[string]struct{}{
		"asc":  {},
		"desc": {},
	}
	_, ok := fields[f]
	return ok
}
