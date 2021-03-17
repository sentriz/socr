create extension if not exists pg_trgm;

create table if not exists screenshots (
    id serial primary key,
    hash text,
    timestamp timestamptz not null,
    dim_width int not null,
    dim_height int not null,
    dominant_colour text not null,
    blurhash text not null
);

create table if not exists dir_infos (
    screenshot_id int,
    filename text,
    directory_alias text,
    primary key (screenshot_id, filename, directory_alias)
);

create table if not exists blocks (
    id serial primary key,
    screenshot_id integer not null references screenshots (id),
    index int not null,
    min_x int not null,
    min_y int not null,
    max_x int not null,
    max_y int not null,
    body text not null
);

create index if not exists idx_blocks_body on blocks using gin (body gin_trgm_ops);

