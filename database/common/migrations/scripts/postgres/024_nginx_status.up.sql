alter table settings_nginx
    add column api_enabled boolean,
    add column api_address varchar(32),
    add column api_port integer,
    add column api_write_enabled boolean;

update settings_nginx set
    api_enabled = true,
    api_address = '127.0.0.1',
    api_port = 8091,
    api_write_enabled = false;

alter table settings_nginx
    alter column api_enabled set not null,
    alter column api_address set not null,
    alter column api_port set not null,
    alter column api_write_enabled set not null;
