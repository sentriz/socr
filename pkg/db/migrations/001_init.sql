create extension pg_trgm;

create table schema_version (
    version smallint primary key
);

insert into schema_version
    values (0);

create type media_type as enum (
    'image',
    'video'
);

create table medias (
    id serial primary key,
    type media_type not null,
    mime text not null,
    hash text not null,
    timestamp timestamptz not null,
    dim_width int not null,
    dim_height int not null,
    dominant_colour text not null,
    blurhash text not null
);

create unique index idx_medias_hash on medias (hash);

create table dir_infos (
    media_id int references medias (id),
    filename text not null,
    directory_alias text not null,
    primary key (media_id, filename, directory_alias)
);

create table blocks (
    id serial primary key,
    media_id integer not null references medias (id),
    index int not null,
    min_x int not null,
    min_y int not null,
    max_x int not null,
    max_y int not null,
    body text not null
);

create index idx_blocks_body on blocks using gin (body gin_trgm_ops);

