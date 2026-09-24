alter table host_vpn drop constraint pk_host_vpn;
alter table host_vpn add constraint pk_host_vpn primary key (host_id, vpn_id, name);
