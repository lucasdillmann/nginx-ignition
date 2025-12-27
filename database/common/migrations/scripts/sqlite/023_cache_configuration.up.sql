create table cache (
    id uuid not null,
    name varchar(256) not null,
    storage_path varchar(512),
    inactive_seconds integer,
    maximum_size_mb integer,
    allowed_methods varchar array not null,
    minimum_uses_before_caching integer not null,
    use_stale varchar array not null,
    background_update boolean not null,
    concurrency_lock_enabled boolean not null,
    concurrency_lock_timeout_seconds integer,
    concurrency_lock_age_seconds integer,
    revalidate boolean not null,
    bypass_rules varchar array not null,
    no_cache_rules varchar array not null,
    file_extensions varchar array not null,
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
    maximum_size_mb,
    allowed_methods,
    minimum_uses_before_caching,
    use_stale,
    background_update,
    concurrency_lock_enabled,
    revalidate,
    bypass_rules,
    no_cache_rules,
    file_extensions
) values (
    '08c8430a-661d-4034-893d-4c31278f99e8',
    'Static assets',
    null,
    259200,
    512,
    '["GET", "HEAD"]',
    1,
    '{"ERROR", "TIMEOUT", "UPDATING", "HTTP_500", "HTTP_502", "HTTP_503", "HTTP_504"}',
    true,
    false,
    true,
    '[]',
    '[]',
    '["ico", "css", "js", "gif", "jpg", "jpeg", "png", "svg", "svgz", "webp", "avif", "woff", "woff2", "ttf", "otf", "mp4", "webm", "wav", "mp3", "m4a", "aac", "ogg", "json", "xml", "html", "htm", "webmanifest"]'
);

alter table "user" add column caches_access_level
    varchar(32) not null default 'NO_ACCESS';

alter table host
    add column cache_id text references cache(id);

alter table host_route
    add column cache_id text references cache(id);

create index idx_host_cache_id
    on host (cache_id);
create index idx_host_route_cache_id
    on host_route (cache_id);
