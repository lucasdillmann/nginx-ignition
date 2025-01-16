create index idx_hosts_search on host (domain_names);
create index idx_certificate_search on certificate (domain_names);
create index idx_user_search on "user" ("name", username);
