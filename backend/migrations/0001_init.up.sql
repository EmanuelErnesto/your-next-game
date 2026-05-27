create type game_status as enum ('played', 'not_played', 'playing', 'paused', 'dropped');

create table users (
  id uuid primary key,
  display_name text not null,
  avatar_url text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table steam_accounts (
  user_id uuid primary key references users(id) on delete cascade,
  steam_id varchar(20) not null unique,
  persona_name text not null,
  avatar_url text,
  profile_url text,
  linked_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table games_catalog (
  id uuid primary key,
  provider text not null check (provider in ('steam', 'local', 'manual')),
  provider_game_id text not null,
  title text not null,
  cover_url text,
  description text not null default '',
  developer text not null default '',
  genre text not null default '',
  release_date date,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (provider, provider_game_id)
);

create table user_games (
  id uuid primary key,
  user_id uuid not null references users(id) on delete cascade,
  catalog_game_id uuid not null references games_catalog(id) on delete restrict,
  steam_app_id bigint,
  status game_status not null default 'not_played',
  hours_played numeric(8,1) not null default 0,
  rating numeric(2,1) check (rating is null or (rating >= 0 and rating <= 5)),
  source text not null check (source in ('seed', 'steam_sync', 'manual')),
  last_synced_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (user_id, catalog_game_id)
);

create table recommendation_sessions (
  id uuid primary key,
  user_id uuid not null references users(id) on delete cascade,
  mood text not null,
  time_available_hours integer not null check (time_available_hours between 1 and 24),
  selected_user_game_id uuid references user_games(id) on delete set null,
  inputs_json jsonb not null,
  result_json jsonb not null,
  created_at timestamptz not null default now()
);

create index idx_user_games_user on user_games (user_id);
create index idx_user_games_status on user_games (status);
create index idx_catalog_provider on games_catalog (provider, provider_game_id);
create index idx_steam_accounts_steam_id on steam_accounts (steam_id);
