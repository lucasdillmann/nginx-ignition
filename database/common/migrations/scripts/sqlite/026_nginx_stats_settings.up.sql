create table nginx_settings_stats (
    id uuid not null,
    enabled boolean not null,
    persistent boolean not null,
    maximum_size_mb integer not null,
    database_location varchar(128),
    constraint pk_nginx_settings_stats primary key (id)
);

insert into nginx_settings_stats (id, enabled, persistent, maximum_size_mb)
values ('32f9a2e6-815c-4b53-b924-11887e74880b', false, false, 16);
