create table if not exists host_global_binding_certificate_override (
    host_id uuid not null,
    global_binding_id uuid not null,
    certificate_id uuid,
    primary key (host_id, global_binding_id),
    constraint fk_host_global_binding_certificate_override_host
        foreign key (host_id)
        references host (id)
        on delete cascade,
    constraint fk_host_global_binding_certificate_override_certificate
        foreign key (certificate_id)
        references certificate (id)
        on delete set null
);
