create table access_list (
    id uuid not null,
    "name" varchar(256) not null,
    realm varchar(256),
    satisfy_all boolean not null,
    default_outcome varchar(8) not null,
    forward_authentication_header boolean not null,
    constraint pk_access_list primary key (id),
    constraint pk_access_list_name unique ("name")
);

create table access_list_credentials (
    id uuid not null,
    access_list_id uuid not null,
    username varchar(256) not null,
    password varchar(256) not null,
    constraint pk_access_list_credentials
        primary key (id),
    constraint fk_access_list_credentials_access_list
        foreign key (access_list_id) references access_list (id)
);

create index idx_access_list_credentials_access_list_id on access_list_credentials (access_list_id);

create table access_list_entry_set (
    id uuid not null,
    access_list_id uuid not null,
    priority integer not null,
    outcome varchar(8) not null,
    source_addresses varchar array not null,
    constraint pk_access_list_entry_set
        primary key (id),
    constraint fk_access_list_entry_set_access_list
        foreign key (access_list_id) references access_list (id)
);

create index idx_access_list_entry_set_access_list_id on access_list_entry_set (access_list_id);

alter table host add column access_list_id uuid;
alter table host_route add column access_list_id uuid;

create index idx_host_access_list_id on host (access_list_id);
create index idx_host_route_access_list_id on host_route (access_list_id);
