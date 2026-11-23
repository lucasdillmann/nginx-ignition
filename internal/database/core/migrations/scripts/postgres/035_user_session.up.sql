create table user_session (
    id uuid not null,
    iteration_id uuid not null,
    user_id uuid not null,
    created_at timestamp with time zone not null,
    expires_at timestamp with time zone not null,
    constraint pk_user_session primary key (id),
    constraint fk_user_session_user foreign key (user_id) references "user" (id) on delete cascade
);

create index idx_user_session_user_id on user_session (user_id);

create index idx_user_session_expires_at on user_session (expires_at);
