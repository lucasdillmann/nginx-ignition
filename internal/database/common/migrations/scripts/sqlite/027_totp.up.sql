alter table "user" add column totp_secret varchar(64);
alter table "user" add column totp_validated boolean default false;
