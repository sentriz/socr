-- name: GetDirInfo :one
select
    1
from
    dir_infos
where
    directory_alias = pggen.arg ('directory_alias')
    and filename = pggen.arg ('filename')
limit 1;

-- name: GetScreenshotByID :one
select
    *
from
    screenshots
where
    id = pggen.arg ('id')
limit 1;

-- name: GetScreenshotByHash :one
select
    *
from
    screenshots
where
    hash = pggen.arg ('hash')
limit 1;

-- name: GetScreenshotWithBlocksByHash :one
select
    screenshots.*,
    array_agg(blocks order by blocks.index) as blocks
from
    screenshots
    join blocks on blocks.screenshot_id = screenshots.id
where
    hash = pggen.arg ('hash')
group by
    screenshots.id
limit 1;

-- name: CreateScreenshot :one
insert into screenshots (hash, timestamp, dim_width, dim_height, dominant_colour, blurhash)
    values (pggen.arg ('hash'), pggen.arg ('timestamp'), pggen.arg ('dim_width'), pggen.arg ('dim_height'), pggen.arg ('dominant_colour'), pggen.arg ('blurhash'))
returning
    *;

-- name: GetAllScreenshots :many
select
    *
from
    screenshots;

-- name: CreateBlock :exec
insert into blocks (screenshot_id, index, min_x, min_y, max_x, max_y, body)
        values (pggen.arg ('screenshot_id'), pggen.arg ('index'), pggen.arg ('min_x'), pggen.arg ('min_y'), pggen.arg ('max_x'), pggen.arg ('max_y'), pggen.arg ('body'));

-- name: CountDirectoriesByAlias :many
select
    directory_alias,
    count(1)
from
    dir_infos
group by
    directory_alias;

-- https://www.postgresql.org/docs/current/pgtrgm.html
-- name: SearchScreenshots :many
select
    screenshots.*,
    array_agg(blocks order by blocks.index) as highlighted_blocks,
    coalesce(avg(similarity (blocks.body, pggen.arg ('body'))), 1.0) as similarity,
    count(1) over () as total
from
    screenshots
    left join blocks on pggen.arg ('body') != ''
        and blocks.screenshot_id = screenshots.id
where
    pggen.arg ('body') = ''
    or blocks.body % pggen.arg ('body')
group by
    screenshots.id
order by
    (case when pggen.arg('sort_field') = 'timestamp' and pggen.arg('sort_order') = 'asc' then timestamp end) asc,
    (case when pggen.arg('sort_field') = 'timestamp' and pggen.arg('sort_order') = 'desc' then timestamp end) desc
limit pggen.arg ('limit') offset pggen.arg ('offset');

-- name: CreateDirInfo :exec
insert into dir_infos (screenshot_id, filename, directory_alias)
    values (pggen.arg ('screenshot_id'), pggen.arg ('filename'), pggen.arg ('directory_alias'))
on conflict
    do nothing;

-- name: GetScreenshotPathByHash :one
select
    dir_infos.filename,
    dir_infos.directory_alias
from
    dir_infos
    join screenshots on screenshots.id = dir_infos.screenshot_id
where
    screenshots.hash = pggen.arg ('hash')
limit 1;

