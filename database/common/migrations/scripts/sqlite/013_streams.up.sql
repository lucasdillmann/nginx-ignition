create table stream (
    id uuid not null,
    enabled boolean not null,
    description varchar(256) not null,
    binding_protocol varchar(32) not null,
    binding_address varchar(256) not null,
    binding_port integer,
    backend_protocol varchar(32) not null,
    backend_address varchar(256) not null,
    backend_port integer,
    use_proxy_protocol boolean not null,
    ssl boolean not null,
    tcp_keep_alive boolean not null,
    tcp_no_delay boolean not null,
    tcp_deferred boolean not null,
    constraint pk_stream primary key (id)
);

create index idx_stream_enabled on stream (enabled);
create index idx_stream_description on stream (description);
