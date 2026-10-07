create table user_token (
    id uuid not null,
    user_id uuid not null,
    name varchar(256) not null,
    expiration timestamp with time zone,
    created_at timestamp with time zone not null,
    constraint pk_user_token primary key (id),
    constraint uk_user_token_name unique (user_id, name),
    constraint fk_user_token_user foreign key (user_id) references "user" (id) on delete cascade
);

create index idx_user_token_user_id on user_token (user_id);