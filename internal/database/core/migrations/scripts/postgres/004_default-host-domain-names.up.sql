alter table host alter column domain_names drop not null;
update host set domain_names = null where default_server;
