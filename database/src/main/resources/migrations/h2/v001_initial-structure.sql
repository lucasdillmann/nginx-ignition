create table certificate (
    id uuid not null,
    domain_names varchar array not null,
    provider_id varchar(64) not null,
    issued_at timestamp with time zone not null,
    valid_until timestamp with time zone not null,
    valid_from timestamp with time zone not null,
    renew_after timestamp with time zone,
    private_key varchar(2048) not null,
    public_key varchar(2048) not null,
    certification_chain text array not null,
    parameters text not null,
    metadata text,
    constraint pk_certificate primary key (id)
);

create index idx_certificate_renew_after on certificate (renew_after);

create table host (
    id uuid not null,
    enabled boolean not null,
    "default" boolean not null,
    domain_names varchar(512) not null,
    websocket_support boolean not null,
    http2_support boolean not null,
    redirect_http_to_https boolean not null,
    constraint pk_host primary key (id)
);

create index idx_host_enabled on host (enabled);
create index idx_host_default on host ("default");

create table host_binding (
    id uuid not null,
    host_id uuid not null,
    type varchar(64) not null,
    ip varchar(256) not null,
    port integer not null,
    certificate_id uuid,
    constraint pk_host_binding primary key (id),
    constraint fk_host_binding_host_id foreign key (host_id) references host (id),
    constraint fk_host_binding_certificate_id foreign key (certificate_id) references certificate (id)
);

create index idx_host_binding_host_id on host_binding (host_id);

create table host_route (
    id uuid not null,
    host_id uuid not null,
    priority integer not null,
    type varchar(64) not null,
    source_path varchar(512) not null,
    target_uri varchar(512),
    custom_settings text,
    static_response_code integer,
    static_response_payload text,
    static_response_headers text,
    redirect_code integer,
    constraint pk_host_route primary key (id),
    constraint fk_host_route_host_id foreign key (host_id) references host (id)
);

create index idx_host_route_host_id on host_route (host_id);

create table "user" (
    id uuid not null,
    enabled boolean not null,
    name varchar(256) not null,
    username varchar(256) not null,
    password_hash varchar(2048) not null,
    password_salt varchar(512) not null,
    role varchar(32) not null,
    constraint pk_user primary key (id)
);

create index idx_user_username on "user" (username);
