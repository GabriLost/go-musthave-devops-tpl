-- +goose Up
CREATE TABLE IF NOT EXISTS counters
(
    id serial PRIMARY KEY,
    name VARCHAR(128) UNIQUE NOT NULL,
    value BIGINT NOT NULL
);
CREATE TABLE IF NOT EXISTS gauges
(
    id serial PRIMARY KEY,
    name VARCHAR(128) UNIQUE NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

-- +goose Down
DROP TABLE counters;
DROP TABLE gauges;