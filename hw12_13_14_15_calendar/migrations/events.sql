-- +goose Up
create table if not exists events
(
    id          uuid      not null,
    title       text      not null,
    started_at  timestamp not null,
    finished_at timestamp not null,
    description text      not null default '',
    user_id     uuid      not null,
    notify      int       not null default 0
);

-- +goose Down
drop table if exists events;