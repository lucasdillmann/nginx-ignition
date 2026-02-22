alter table host_vpn add column enable_https boolean not null default true;
alter table host_vpn add column certificate_id uuid;

create index idx_host_vpn_certificate_id on host_vpn (certificate_id);
