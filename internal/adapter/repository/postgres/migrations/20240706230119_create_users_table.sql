-- +goose Up

-- users
create table public.users (
    id uuid primary key not null default gen_random_uuid(),
    login varchar not null,
    password varchar(60) not null,
    created_at timestamp without time zone not null default CURRENT_TIMESTAMP
);
create unique index users_login_unique on public.users (login);
comment on table public.users is 'Users';

-- +goose Down
drop index public.users_login_unique;
drop table public.users cascade;
