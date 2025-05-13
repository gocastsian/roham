-- +migrate Up
create table layers(
    id bigserial primary key ,
    name varchar(255) unique not null ,
    default_style varchar(255) not null ,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW()
);

-- +migrate Down
drop table layers;