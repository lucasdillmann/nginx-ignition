drop index idx_host_vpn_host_id;
drop index idx_host_vpn_vpn_id;
alter table host_vpn rename to host_vpn_old;

create table host_vpn (
    host_id uuid not null,
    vpn_id  uuid not null,
    name varchar(256) not null,
    host varchar(256),
    constraint pk_host_vpn primary key (host_id, vpn_id, name),
    constraint fk_host_vpn_host foreign key (host_id) references host (id),
    constraint fk_host_vpn_vpn foreign key (vpn_id) references vpn (id)
);

create index idx_host_vpn_host_id on host_vpn (host_id);
create index idx_host_vpn_vpn_id on host_vpn (vpn_id);

insert into host_vpn (host_id, vpn_id, name, host)
    select host_id, vpn_id, name, host from host_vpn_old;

drop table host_vpn_old;
