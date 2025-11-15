create table host_vpn_new (
    host_id integer not null,
    vpn_id integer not null,
    name varchar(255) not null,
    primary key (host_id, vpn_id, name),
    foreign key (host_id) references host (id) on delete cascade,
    foreign key (vpn_id) references integration (id) on delete cascade
);

insert into host_vpn_new (host_id, vpn_id, name)
select host_id, vpn_id, name from host_vpn;

drop table host_vpn;
alter table host_vpn_new rename to host_vpn;
