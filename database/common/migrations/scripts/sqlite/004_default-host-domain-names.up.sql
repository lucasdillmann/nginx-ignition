drop index idx_hosts_search;

alter table host add column domain_names_v2 varchar array;
update host set domain_names_v2 = domain_names;

alter table host drop column domain_names;
alter table host rename column domain_names_v2 to domain_names;

update host set domain_names = null where default_server;

create index idx_hosts_search on host (domain_names);
