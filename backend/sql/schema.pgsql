CREATE EXTENSION IF NOT EXISTS pg_trgm;


CREATE TABLE IF NOT EXISTS directories (
    id INT PRIMARY KEY,
    alias TEXT NOT NULL,
    rel_path TEXT NOT NULL
);


CREATE TYPE filetype AS ENUM ('GIF', 'PNG', 'JPG');


CREATE TABLE IF NOT EXISTS tags (id INT PRIMARY KEY, NAME TEXT NOT NULL);


CREATE TABLE IF NOT EXISTS screenshot_properties (
    id BIGINT PRIMARY KEY REFERENCES directories(id),
    dim_width INT,
    dim_height INT,
    dominant_colour TEXT,
    blurhash TEXT
);


CREATE TABLE IF NOT EXISTS screenshots (
    id BIGINT PRIMARY KEY,
    stamp TIMESTAMP NOT NULL,
    directory INT NOT NULL REFERENCES directories(id),
    filename TEXT NOT NULL,
    filetype filetype NOT NULL
);


CREATE TABLE IF NOT EXISTS tag_screenshots (
    tag_id BIGINT,
    screenshot_id INT,
    PRIMARY KEY (tag_id, screenshot_id)
);


CREATE TABLE IF NOT EXISTS blocks (
    screenshot_id BIGINT NOT NULL REFERENCES screenshot(id),
    min_x INT NOT NULL,
    min_y INT NOT NULL,
    max_x INT NOT NULL,
    max_y INT NOT NULL,
    body TEXT NOT NULL
);


CREATE INDEX IF NOT EXISTS blocks_body ON blocks USING GIN(body gin_trgm_ops);

