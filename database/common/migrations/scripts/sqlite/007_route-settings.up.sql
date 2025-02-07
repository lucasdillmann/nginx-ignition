alter table host_route
    add column include_forward_headers boolean not null default true;
alter table host_route
    add column proxy_ssl_server_name boolean not null default true;
alter table host_route
    add column keep_original_domain_name boolean not null default true;
