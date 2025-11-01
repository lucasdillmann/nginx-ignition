create table vpn (
    id         uuid         not null,
    driver     varchar(64)  not null,
    name       varchar(256) not null,
    enabled    boolean      not null,
    parameters text         not null,
    constraint pk_vpn primary key (id),
    constraint uk_vpn unique (name)
);

create table host_vpn (
    host_id uuid not null,
    vpn_id  uuid not null,
    name varchar(256) not null,
    constraint pk_host_vpn primary key (host_id, vpn_id),
    constraint fk_host_vpn_host foreign key (host_id) references host (id),
    constraint fk_host_vpn_vpn foreign key (vpn_id) references vpn (id)
);

create index idx_host_vpn_host_id on host_vpn (host_id);
create index idx_host_vpn_vpn_id on host_vpn (vpn_id);

alter table "user" add column vpns_access_level varchar(255) default 'NO_ACCESS';
