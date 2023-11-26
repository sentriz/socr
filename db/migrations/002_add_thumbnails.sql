create table thumbnails (
    id serial primary key,
    media_id integer not null references medias (id) on delete cascade,
    mime text not null,
    dim_width int not null,
    dim_height int not null,
    timestamp timestamptz not null,
    data bytea not null
);

create unique index idx_thumbnails_media_id on thumbnails (media_id);

