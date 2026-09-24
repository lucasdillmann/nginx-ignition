create table integration (
    id varchar(128) not null,
    enabled boolean,
    parameters text not null,
    constraint pk_integration primary key (id)
);

alter table host_route add column integration_id varchar(128);
alter table host_route add column integration_option_id varchar(256);
