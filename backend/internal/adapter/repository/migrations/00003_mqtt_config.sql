-- +goose Up
CREATE TABLE mqtt_config (
    id            INT PRIMARY KEY DEFAULT 1,
    enabled       BOOLEAN NOT NULL DEFAULT FALSE,
    host          TEXT NOT NULL DEFAULT 'localhost',
    port          INT NOT NULL DEFAULT 1883,
    publish_topic TEXT NOT NULL DEFAULT 'same/alerts',
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
INSERT INTO mqtt_config (id) VALUES (1);

-- +goose Down
DROP TABLE mqtt_config;
