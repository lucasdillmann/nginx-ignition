alter table certificate rename to server_certificate;
alter table "user" rename column certificates_access_level to server_certificates_access_level;
alter table "user" add column client_certificates_access_level varchar(32) not null default 'NO_ACCESS';
