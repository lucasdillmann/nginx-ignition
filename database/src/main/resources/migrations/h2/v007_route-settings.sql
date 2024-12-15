alter table host_route
    add column include_forward_headers boolean not null default true;
alter table host_route
    add column proxy_ssl_server_name boolean not null default true;
alter table host_route
    add column keep_original_domain_name boolean not null default true;
alter table host_route
    add column forward_query_params boolean not null default true;

alter table host_route
    alter column include_forward_headers drop default;
alter table host_route
    alter column proxy_ssl_server_name drop default;
alter table host_route
    alter column keep_original_domain_name drop default;
alter table host_route
    alter column forward_query_params drop default;
