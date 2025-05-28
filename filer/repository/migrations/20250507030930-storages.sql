-- +migrate Up

CREATE TABLE storages
(
    id   SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    kind VARCHAR(255) NOT NULL
);

INSERT INTO storages (name, kind) VALUES ('avatars', 'avatar');

-- +migrate Down

DROP TABLE IF EXISTS storages;