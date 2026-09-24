alter table host_route add column directory_listing_enabled boolean not null default false;
update host_route set type = 'STATIC_FILES', directory_listing_enabled = true where type = 'DIRECTORY';
update host_route set type = 'EXECUTE_CODE' where type = 'SOURCE_CODE';
