alter table integration rename to integration_old;
alter table integration_old drop constraint pk_integration;

create table integration (
    id         uuid         not null,
    driver     varchar(64)  not null,
    name       varchar(256) not null,
    enabled    boolean      not null,
    parameters text         not null,
    constraint pk_integration primary key (id),
    constraint uk_integration unique (name)
);

insert into integration (id, driver, name, enabled, parameters)
select
    lower(hex(randomblob(4))) || '-' ||
    lower(hex(randomblob(2))) || '-' ||
    '4' || substr(lower(hex(randomblob(2))), 2) || '-' ||
    substr('89ab', 1 + (abs(random()) % 4), 1) || substr(lower(hex(randomblob(2))), 2) || '-' ||
    lower(hex(randomblob(6))) as id,
    case id
        when 'DOCKER' then 'DOCKER'
        when 'TRUENAS_SCALE' then 'TRUENAS'
        end as driver,
    case id
        when 'DOCKER' then 'Docker'
        when 'TRUENAS_SCALE' then 'TrueNAS'
        end as name,
    enabled,
    parameters
from integration_old;

drop table integration_old;
create index idx_integration_name on integration (name);
create index idx_integration_enabled on integration (enabled);

alter table host_route rename column integration_id to integration_id_old;
alter table host_route add column integration_id uuid;
update host_route
    set integration_id = case integration_id_old
        when 'DOCKER' then (select id from integration where driver = 'DOCKER')
        when 'TRUENAS_SCALE' then (select id from integration where driver = 'TRUENAS')
    end
where integration_id_old is not null;

alter table host_route drop column integration_id_old;
alter table host_route add constraint fk_host_route_integration foreign key (integration_id) references integration (id);
create index idx_host_route_integration_id on host_route (integration_id);

