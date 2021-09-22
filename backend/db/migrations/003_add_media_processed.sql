alter table medias
    add column processed boolean not null default false;

update
    medias
set
    processed = true;

