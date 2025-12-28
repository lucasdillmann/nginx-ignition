alter table settings_nginx
    add column api_enabled boolean not null default true,
    add column api_address varchar(32) not null default '127.0.0.1',
    add column api_port integer not null default 8091,
    add column api_write_enabled boolean not null default false;

alter table settings_nginx
    alter column api_enabled drop default,
    alter column api_address drop default,
    alter column api_port drop default,
    alter column api_write_enabled drop default;
