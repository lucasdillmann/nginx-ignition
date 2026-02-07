create table settings_nginx_stats (
    id uuid not null,
    enabled boolean not null,
    persistent boolean not null,
    maximum_size_mb integer not null,
    database_location varchar(128),
    constraint pk_settings_nginx_stats primary key (id)
);

insert into settings_nginx_stats (id, enabled, persistent, maximum_size_mb)
values ('32f9a2e6-815c-4b53-b924-11887e74880b', false, false, 64);

alter table nginx_settings_buffers rename to settings_nginx_buffers;
