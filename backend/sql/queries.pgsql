-- name: GetScreenshotByPath :one
select
    *
from
    screenshots
where
    directory_alias = $1
    and filename = $2
limit 1;

-- name: GetScreenshotByID :one
select
    *
from
    screenshots
where
    id = $1
limit 1;

-- name: CreateScreenshot :one
insert into screenshots (id, timestamp, directory_alias, filename, dim_width, dim_height, dominant_colour, blurhash)
    values ($1, $2, $3, $4, $5, $6, $7, $8)
returning
    *;

-- name: GetAllScreenshots :many
select
    *
from
    screenshots;

-- name: CreateBlock :exec
insert into blocks (screenshot_id, index, min_x, min_y, max_x, max_y, body)
        values ($1, $2, $3, $4, $5, $6, $7);

-- name: SearchBlock :many
select
    *
from
    screenshots
where (@body::text) % body
limit 40;

