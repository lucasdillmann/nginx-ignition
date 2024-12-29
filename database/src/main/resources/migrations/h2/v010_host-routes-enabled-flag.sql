alter table host_route add column enabled boolean not null default true;
alter table host_route alter column enabled drop default;
