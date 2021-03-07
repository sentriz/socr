-- name: GetScreenshotByPath :one
select
    *
from
    screenshots
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

-- name: CreateScreenshot :one
insert into screenshots (id, timestamp, directory_alias, filename, dim_width, dim_height, dominant_colour, blurhash)
    values (pggen.arg ('id'), pggen.arg ('timestamp'), pggen.arg ('directory_alias'), pggen.arg ('filename'), pggen.arg ('dim_width'), pggen.arg ('dim_height'), pggen.arg ('dominant_colour'), pggen.arg ('blurhash'))
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
    screenshots
group by
    directory_alias;

-- name: SearchScreenshots :many
select
    screenshots.*,
    array_agg(blocks order by blocks.index) as blocks,
    avg(similarity (blocks.body, pggen.arg ('body'))) as similarity
from
    screenshots
    join blocks on blocks.screenshot_id = screenshots.id
where
    blocks.body % pggen.arg ('body')
group by
    screenshots.id
order by
    similarity desc
limit pggen.arg ('limit') offset pggen.arg ('offset');

