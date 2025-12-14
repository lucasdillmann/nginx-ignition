alter table certificate
    rename to server_certificate;
alter table host_binding
    rename column certificate_id to server_certificate_id;
alter table settings_global_binding
    rename column certificate_id to server_certificate_id;
alter table "user"
    rename column certificates_access_level to server_certificates_access_level;
alter table "user"
    add column client_certificates_access_level varchar(32) not null default 'NO_ACCESS';

create table client_certificate (
    id text not null primary key,
    name varchar(255) not null,
    type varchar(32) not null,
    validation_mode varchar(32) not null,
    stapling_enabled boolean not null,
    stapling_verify boolean not null,
    stapling_responder_url varchar(500),
    stapling_responder_file_path varchar(500),
    ca_public_key blob,
    ca_private_key blob,
    ca_send_to_clients boolean
);

create table client_certificate_item (
    id text not null primary key,
    client_certificate_id text not null references client_certificate(id),
    dn varchar(255),
    public_key blob,
    private_key blob,
    issued_at datetime,
    expires_at datetime,
    revoked boolean
);

create index idx_client_certificate_item_certificate_id
    on client_certificate_item (client_certificate_id);

alter table host
    add column client_certificate_id text references client_certificate(id);
alter table host_route
    add column client_certificate_id text references client_certificate(id);
