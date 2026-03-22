alter table host_route add column ignore_ssl_errors boolean not null default false;
alter table host_route add column integration_use_https boolean not null default false;
