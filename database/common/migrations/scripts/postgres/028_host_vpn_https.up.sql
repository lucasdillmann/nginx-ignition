alter table host_vpn add column enable_https boolean not null default true;
alter table host_vpn add column certificate_id uuid;

alter table host_vpn add constraint fk_host_vpn_certificate foreign key (certificate_id) references certificate (id);
create index idx_host_vpn_certificate_id on host_vpn (certificate_id);
