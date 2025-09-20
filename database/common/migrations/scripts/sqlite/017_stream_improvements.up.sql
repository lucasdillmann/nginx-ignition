create table stream_route (
    id uuid not null,
    stream_id uuid not null,
    domain_name varchar(256) not null,
    constraint stream_route_pk primary key (id),
    constraint stream_route_stream_id_fk foreign key (stream_id) references stream (id),
    constraint stream_domain_name_uk unique (stream_id, domain_name)
);

create table stream_backend (
    id uuid not null,
    stream_id uuid,
    stream_route_id uuid,
    protocol varchar(64) not null,
    address varchar(512) not null,
    port int,
    weight int,
    max_failures int,
    open_seconds int,
    constraint stream_backend_pk primary key (id),
    constraint stream_backend_stream_id_fk foreign key (stream_id) references stream (id),
    constraint stream_backend_stream_route_id_fk foreign key (stream_route_id) references stream_route (id),
    constraint stream_backend_reference_chk check (stream_id is not null or stream_route_id is not null)
);

create index stream_route_stream_id_idx on stream_route (stream_id);
create index stream_backend_stream_id_idx on stream_backend (stream_id);
create index stream_backend_stream_route_id_idx on stream_backend (stream_route_id);

alter table stream add column type varchar(64) not null default 'SIMPLE';

do $$
declare
    row record;
begin
    for row in select * from stream
    loop
        insert into stream_backend (id, stream_id, protocol, address, port)
        values (
            row.id,
            row.id,
            row.backend_protocol,
            row.backend_address,
            row.backend_port
        );
    end loop;
end $$;

alter table stream drop column backend_protocol;
alter table stream drop column backend_address;
alter table stream drop column backend_port;
