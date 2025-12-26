create table cache (
    id uuid not null,
    name varchar(256) not null,
    storage_path text,
    inactive_seconds integer,
    max_size_mb integer,
    allowed_methods varchar array not null,
    minimum_uses_before_caching integer,
    use_stale varchar array not null,
    background_update boolean,
    concurrency_lock_enabled boolean not null,
    concurrency_lock_timeout_seconds integer,
    concurrency_lock_age_seconds integer,
    revalidate boolean,
    bypass_rules varchar array not null,
    no_cache_rules varchar array not null,
    constraint pk_cache primary key (id)
);

create table cache_duration (
    id uuid not null,
    cache_id uuid not null,
    status_codes integer array not null,
    valid_time_seconds integer not null,
    constraint pk_cache_duration primary key (id),
    constraint fk_cache_duration_cache foreign key (cache_id) references cache (id) on delete cascade
);

create index idx_cache_duration_cache_id on cache_duration (cache_id);

insert into cache (
    id,
    name,
    storage_path,
    inactive_seconds,
    max_size_mb,
    allowed_methods,
    minimum_uses_before_caching,
    use_stale,
    background_update,
    concurrency_lock_enabled,
    concurrency_lock_timeout_seconds,
    concurrency_lock_age_seconds,
    revalidate,
    bypass_rules,
    no_cache_rules
) values (
    '08c8430a-661d-4034-893d-4c31278f99e8',
    'Static files',
    null,
    259200,
    512,
    '["GET", "HEAD"]',
    1,
    '["error", "timeout", "updating", "http_500", "http_502", "http_503", "http_504"]',
    true,
    true,
    5,
    5,
    true,
    '[]',
    '[]'
);

insert into cache_duration (
    id,
    cache_id,
    status_codes,
    valid_time_seconds
) values (
    '8b2847a1-90be-459f-958b-308102e3b2e5',
    '08c8430a-661d-4034-893d-4c31278f99e8',
    '[200, 301, 302]',
    259200
);

alter table "user" add column caches_access_level
    varchar(32) not null default 'NO_ACCESS';

alter table host
    add column cache_id text references cache(id);

alter table host_route
    add column cache_id text references cache(id);

create index idx_host_cache_id on host (cache_id);
create index idx_host_route_cache_id on host_route (cache_id);
