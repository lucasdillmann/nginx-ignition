create table nginx_settings_buffers (
    id uuid not null,
    client_body_kb integer not null,
    client_header_kb integer not null,
    large_client_header_size_kb integer not null,
    large_client_header_amount integer not null,
    output_size_kb integer not null,
    output_amount integer not null,
    constraint pk_nginx_settings_buffers primary key (id)
);

insert into nginx_settings_buffers values (
    '7e3f5a8b-2d4c-4e9b-a1f3-6c8d9e0f1a2b',
    16,
    1,
    8,
    4,
    32,
    4
);

alter table settings_nginx add column client_body_timeout integer not null default 60;
alter table settings_nginx add column tcp_nodelay_enabled boolean not null default true;
alter table settings_nginx add column custom text;
