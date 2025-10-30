alter table integration rename to integration_old;

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
    case id
        when 'DOCKER' then '83fa2c8c-8f53-495f-9d38-5b1c619e4c17'
        when 'TRUENAS_SCALE' then '6a14e7c9-88ee-4025-8554-de368071b0a9'
        end as id,
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
        when 'DOCKER' then '83fa2c8c-8f53-495f-9d38-5b1c619e4c17'
        when 'TRUENAS_SCALE' then '6a14e7c9-88ee-4025-8554-de368071b0a9'
    end
where integration_id_old is not null;

alter table host_route drop column integration_id_old;
alter table host_route add constraint fk_host_route_integration foreign key (integration_id) references integration (id);
create index idx_host_route_integration_id on host_route (integration_id);
