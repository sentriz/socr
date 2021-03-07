-- name: GetScreenshotByPath :one
select
    *
from
    screenshots
where
    directory_alias = pggen.arg ('DirectoryAlias')
    and filename = pggen.arg ('Filename')
limit 1;

-- name: GetScreenshotByID :one
select
    *
from
    screenshots
where
    id = pggen.arg ('ID')
limit 1;

-- name: CreateScreenshot :one
insert into screenshots (id, timestamp, directory_alias, filename, dim_width, dim_height, dominant_colour, blurhash)
    values (pggen.arg ('Id'), pggen.arg ('Timestamp'), pggen.arg ('DirectoryAlias'), pggen.arg ('Filename'), pggen.arg ('DimWidth'), pggen.arg ('DimHeight'), pggen.arg ('DominantColour'), pggen.arg ('Blurhash'))
returning
    *;

-- name: GetAllScreenshots :many
select
    *
from
    screenshots;

-- name: CreateBlock :exec
insert into blocks (screenshot_id, index, min_x, min_y, max_x, max_y, body)
        values (pggen.arg ('ScreenshotId'), pggen.arg ('Index'), pggen.arg ('MinX'), pggen.arg ('MinY'), pggen.arg ('MaxX'), pggen.arg ('MaxY'), pggen.arg ('Body'));

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
    array_agg(blocks) as blocks
from
    screenshots
    join blocks on blocks.screenshot_id = screenshots.id
where
    pggen.arg ('Body') % blocks.body
group by
    screenshots.id
limit pggen.arg ('Limit') offset pggen.arg ('Offset');

