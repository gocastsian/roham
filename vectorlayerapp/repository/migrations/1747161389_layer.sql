-- +migrate Up

CREATE TABLE styles
(
    id         BIGSERIAL PRIMARY KEY,
    file_path  VARCHAR(512) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE layers
(
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(255) UNIQUE NOT NULL,
    geom_type     VARCHAR(50)         NOT NULL,
    default_style BIGINT        NOT NULL REFERENCES styles (id) ON DELETE RESTRICT,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW()
);

CREATE TABLE layer_styles
(
    id         BIGSERIAL PRIMARY KEY,
    layer_id   BIGINT NOT NULL REFERENCES layers (id) ON DELETE CASCADE,
    style_id   BIGINT NOT NULL REFERENCES styles (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +migrate Down

DROP TABLE IF EXISTS layer_styles;
DROP TABLE IF EXISTS layers;
DROP TABLE IF EXISTS styles;
