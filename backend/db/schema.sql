create extension if not exists pg_trgm;

create table if not exists screenshots (
    id serial primary key,
    hash text not null,
    timestamp timestamptz not null,
    dim_width int not null,
    dim_height int not null,
    dominant_colour text not null,
    blurhash text not null
);

create table if not exists videos (
    id serial primary key,
    hash text not null,
    timestamp timestamptz not null,
    dim_width int not null,
    dim_height int not null
);

create table if not exists dir_infos (
    id serial primary key,
    screenshot_id int references screenshots (id),
    video_id int references videos (id),
    filename text not null,
    directory_alias text not null,
    unique (screenshot_id, filename, directory_alias),
    unique (video_id, filename, directory_alias)
);

create unique index if not exists idx_dir_infos_screenshot_id on dir_infos (screenshot_id);

create unique index if not exists idx_dir_infos_video_id on dir_infos (video_id);

create unique index if not exists idx_dir_infos_path on dir_infos (filename, directory_alias);

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

