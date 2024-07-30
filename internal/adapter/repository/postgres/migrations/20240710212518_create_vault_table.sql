-- +goose Up
create type vault_kind as enum ('creds', 'note', 'card', 'file');
create table public.vault (
  id uuid primary key not null default gen_random_uuid(),
  user_id uuid not null,
  kind vault_kind not null,
  data bytea not null,
  created_at timestamp without time zone not null default current_timestamp,
  updated_at timestamp without time zone not null default current_timestamp,
  foreign key (user_id) references public.users (id)
    on update cascade on delete cascade
);
create index vault_user_id_index on vault (user_id);
comment on table public.vault is 'Vault';

-- +goose Down
drop index public.vault_user_id_index;
drop table public.vault cascade;
drop type vault_kind;
