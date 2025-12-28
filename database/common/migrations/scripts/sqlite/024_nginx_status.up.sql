alter table settings_nginx
    add column api_enabled boolean not null default false;

alter table settings_nginx
    add column api_address varchar(32) not null default '127.0.0.1';

alter table settings_nginx
    add column api_port integer not null default 8091;

alter table settings_nginx
    add column api_write_enabled boolean not null default false;
