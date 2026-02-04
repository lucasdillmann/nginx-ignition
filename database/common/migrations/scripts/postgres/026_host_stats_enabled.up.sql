alter table host add column stats_enabled boolean not null default false;
alter table host alter column stats_enabled drop default;
