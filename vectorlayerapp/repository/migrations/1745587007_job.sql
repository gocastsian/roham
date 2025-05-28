-- +migrate Up
create type status as Enum ('completed' , 'processing' , 'failed' , 'pending');
create table jobs
(
    id         bigserial primary key,
    token      varchar(199) not null unique,
    status     status       not null,
    error      text,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


-- +migrate Down

drop table jobs;