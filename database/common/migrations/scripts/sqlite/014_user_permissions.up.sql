alter table "user"
    add column hosts_access_level varchar(32) not null default 'NO_ACCESS';
alter table "user"
    add column streams_access_level varchar(32) not null default 'NO_ACCESS';
alter table "user"
    add column certificates_access_level varchar(32) not null default 'NO_ACCESS';
alter table "user"
    add column settings_access_level varchar(32) not null default 'NO_ACCESS';
alter table "user"
    add column users_access_level varchar(32) not null default 'NO_ACCESS';
alter table "user"
    add column logs_access_level varchar(32) not null default 'NO_ACCESS';
alter table "user"
    add column integrations_access_level varchar(32) not null default 'NO_ACCESS';
alter table "user"
    add column access_lists_access_level varchar(32) not null default 'NO_ACCESS';
alter table "user"
    add column nginx_server_access_level varchar(32) not null default 'READ_ONLY';

update "user"
set hosts_access_level        = 'READ_WRITE',
    streams_access_level      = 'READ_WRITE',
    certificates_access_level = 'READ_WRITE',
    settings_access_level     = 'READ_WRITE',
    users_access_level        = 'READ_WRITE',
    logs_access_level         = 'READ_ONLY',
    integrations_access_level = 'READ_WRITE',
    access_lists_access_level = 'READ_WRITE',
    nginx_server_access_level = 'READ_WRITE'
where role = 'ADMIN';

update "user"
set hosts_access_level        = 'READ_WRITE',
    streams_access_level      = 'READ_WRITE',
    certificates_access_level = 'READ_WRITE',
    settings_access_level     = 'NO_ACCESS',
    users_access_level        = 'NO_ACCESS',
    logs_access_level         = 'READ_ONLY',
    integrations_access_level = 'NO_ACCESS',
    access_lists_access_level = 'READ_WRITE',
    nginx_server_access_level = 'READ_WRITE'
where role = 'USER';

alter table "user" drop column role;
