create extension if not exists pg_trgm;

create table if not exists screenshots (
    id serial primary key,
    hash text,
    timestamp timestamptz not null,
    directory_alias text not null,
    filename text not null,
    dim_width int not null,
    dim_height int not null,
    dominant_colour text not null,
    blurhash text not null
);

create unique index if not exists idx_screenshots_directory_alias_filename on screenshots (directory_alias, filename);

create table if not exists tags (
    id int primary key,
    name text not null
);

create table if not exists tag_screenshots (
    tag_id int,
    screenshot_id integer,
    primary key (tag_id, screenshot_id)
);

create table if not exists blocks (
    id serial primary key,
    screenshot_id integer not null references screenshots (id),
    index smallint not null,
    min_x smallint not null,
    min_y smallint not null,
    max_x smallint not null,
    max_y smallint not null,
    body text not null
);

create index if not exists idx_blocks_body on blocks using gin (body gin_trgm_ops);

