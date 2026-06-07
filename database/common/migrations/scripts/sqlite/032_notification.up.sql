alter table "user" add column notification_language varchar(16) not null default 'en';

create table notification (
    id uuid not null,
    user_id uuid not null,
    title text not null,
    summary text not null,
    category varchar(64) not null,
    payload text not null,
    read_at timestamp,
    created_at timestamp not null,
    delivery_completed boolean not null,
    constraint pk_notification primary key (id),
    constraint fk_notification_user foreign key (user_id) references "user" (id)
);

create index idx_notification_user_id on notification (user_id);
create index idx_notification_category on notification (category);
create index idx_notification_created_at on notification (created_at);
create index idx_notification_read_at on notification (read_at);
create index idx_notification_delivery_completed on notification (delivery_completed);

create table notification_configuration (
    id uuid not null,
    user_id uuid not null,
    name varchar(256) not null,
    provider varchar(64) not null,
    enabled boolean not null,
    parameters text not null,
    categories text,
    constraint pk_notification_configuration primary key (id),
    constraint fk_notification_configuration_user foreign key (user_id) references "user" (id),
    constraint uq_notification_configuration_user_name unique (user_id, name)
);

create table notification_provider_submission (
    id uuid not null,
    notification_id uuid not null,
    configuration_id uuid not null,
    provider varchar(64) not null,
    status varchar(32) not null,
    attempt_count integer not null,
    last_error text,
    last_attempt_at timestamp,
    succeeded_at timestamp,
    constraint pk_notification_provider_submission primary key (id),
    constraint fk_notification_provider_submission_notification
        foreign key (notification_id) references notification (id) on delete cascade,
    constraint fk_notification_provider_submission_configuration
        foreign key (configuration_id) references notification_configuration (id) on delete cascade
);

create index idx_notification_provider_submission_status on notification_provider_submission (status);

create table notification_related_entity (
    notification_id uuid not null,
    entity_type varchar(64) not null,
    entity_id uuid not null,
    name varchar(256),
    constraint fk_notification_related_entity_notification
        foreign key (notification_id) references notification (id) on delete cascade
);

create index idx_notification_related_entity_entity
    on notification_related_entity (entity_type, entity_id, notification_id);
create index idx_notification_related_entity_notification_id
    on notification_related_entity (notification_id);
