package db

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

//go:embed schema.pgsql
var schema string

type DB struct {
	*pgxpool.Pool
	*DBQuerier
	squirrel.StatementBuilderType
}

func New(dsn string) (*DB, error) {
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("create and connect pool: %w", err)
	}

	if _, err := pool.Exec(context.Background(), schema); err != nil {
		return nil, fmt.Errorf("executing schema: %w", err)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	return &DB{
		Pool:                 pool,
		DBQuerier:            NewQuerier(pool),
		StatementBuilderType: builder,
	}, nil
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
		Offset(uint64(options.Offset))
	if options.Directory != "" {
		q = q.
			Join("dir_infos on dir_infos.screenshot_id = screenshots.id").
			Where("dir_infos.directory_alias = ?", options.Directory)
	}
	if options.SortField != "" && options.SortOrder != "" {
		q = q.OrderBy(fmt.Sprintf("%s %s", options.SortField, options.SortOrder))
	}
	if options.Body != "" {
		q = q.
			Column("json_agg(blocks order by blocks.index) as highlighted_blocks").
			Column("avg(similarity(blocks.body, ?)) as similarity", options.Body).
			LeftJoin("blocks on blocks.screenshot_id = screenshots.id").
			Where("blocks.body % ?", options.Body).
			GroupBy("screenshots.id")
	}

	sql, args := q.MustSql()
	var results []*Screenshot
	return results, pgxscan.Select(context.Background(), db, &results, sql, args...)
}
