create table settings_global_binding (
    id uuid not null,
    "type" varchar(64) not null,
    ip varchar(256) not null,
    port integer not null,
    certificate_id uuid,
    constraint pk_settings_global_binding primary key (id),
    constraint fk_settings_global_bindings_certificate_id foreign key (certificate_id) references certificate (id)
);

insert into settings_global_binding values (
    'ccdb59e0-8641-4816-9c21-472a93c4095a',
    'HTTP',
    '0.0.0.0',
    80,
    null
);

create table settings_nginx (
    id uuid not null,
    worker_processes integer not null,
    worker_connections integer not null,
    server_tokens_enabled boolean not null,
    sendfile_enabled boolean not null,
    gzip_enabled boolean not null,
    default_content_type varchar(128) not null,
    maximum_body_size_mb integer not null,
    read_timeout integer not null,
    connect_timeout integer not null,
    send_timeout integer not null,
    keepalive_timeout integer not null,
    server_logs_enabled boolean not null,
    server_logs_level varchar(8) not null,
    access_logs_enabled boolean not null,
    error_logs_enabled boolean not null,
    error_logs_level varchar(8) not null
);

insert into settings_nginx values (
    '5fb063d8-343c-42de-bf48-823e71cfb86b',
    2,
    1024,
    false,
    true,
    true,
    'application/octet-stream',
    1024,
    300,
    5,
    300,
    30,
    true,
    'ERROR',
    true,
    true,
    'ERROR'
);

create table settings_log_rotation (
    id uuid not null,
    enabled boolean not null,
    maximum_lines integer not null,
    interval_unit varchar(32) not null,
    interval_unit_count integer not null
);

insert into settings_log_rotation values (
    '80d8875f-e15a-42a4-a5a7-b987ade7d8dc',
    true,
    10000,
    'HOURS',
    1
);

create table settings_certificate_auto_renew (
    id uuid not null,
    enabled boolean not null,
    interval_unit varchar(32) not null,
    interval_unit_count integer not null
);

insert into settings_certificate_auto_renew values (
    '9459be07-b915-4d90-ad2d-f842dca2af59',
    true,
    'HOURS',
    1
);

alter table host add column use_global_bindings boolean not null default false;
alter table host alter column use_global_bindings drop default;
