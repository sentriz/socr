-- name: GetScreenshotByID :one
SELECT
    *
FROM
    screenshots
WHERE
    id = $1
LIMIT
    1;


-- name: CreateScreenshot :exec
INSERT INTO
    screenshots (id, directory, filename, stamp)
VALUES
    ($1, $2, $3, $4);


-- name: GetAllScreenshots :many
SELECT
    *
FROM
    screenshots;


-- name: CreateBlock :exec
INSERT INTO
    blocks (min_x, min_y, max_x, max_y, body)
VALUES
    ($1, $2, $3, $4, $5);


-- name: SearchBlock :many
SELECT
    *
FROM
    screenshots
WHERE
    sqlc.arg(body) :: TEXT % body
LIMIT
    40;

