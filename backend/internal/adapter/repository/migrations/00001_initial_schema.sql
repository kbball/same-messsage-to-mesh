-- +goose Up

-- Decoded SAME/EAS alerts
CREATE TABLE same_alerts (
    id          BIGSERIAL PRIMARY KEY,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    originator  VARCHAR(3)  NOT NULL,
    event_code  VARCHAR(3)  NOT NULL,
    fips_codes  TEXT[]      NOT NULL,
    purge_time  VARCHAR(4)  NOT NULL,
    issue_time  VARCHAR(7)  NOT NULL,
    call_sign   VARCHAR(8)  NOT NULL,
    raw_message TEXT        NOT NULL,
    published   BOOLEAN     NOT NULL DEFAULT FALSE
);

CREATE INDEX same_alerts_received_at_idx ON same_alerts (received_at DESC);

-- SAME/EAS event type reference data
CREATE TABLE event_codes (
    code        VARCHAR(3) PRIMARY KEY,
    description TEXT       NOT NULL,
    category    TEXT       NOT NULL DEFAULT '',
    is_warning  BOOLEAN    NOT NULL DEFAULT FALSE,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- State/county FIPS reference data used to match SAME messages
CREATE TABLE fips_codes (
    state_code  VARCHAR(2) NOT NULL,
    county_code VARCHAR(3) NOT NULL,
    state_name  TEXT       NOT NULL,
    county_name TEXT       NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (state_code, county_code)
);

-- Singleton filter configuration (id always 1)
CREATE TABLE filter_config (
    id          INT     PRIMARY KEY DEFAULT 1,
    state_codes TEXT[]  NOT NULL DEFAULT '{}',
    fips_codes  TEXT[]  NOT NULL DEFAULT '{}',
    event_codes TEXT[]  NOT NULL DEFAULT '{}',
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT filter_config_singleton CHECK (id = 1)
);

INSERT INTO filter_config (id) VALUES (1);

-- Singleton SDR hardware configuration (id always 1)
CREATE TABLE sdr_config (
    id          INT     PRIMARY KEY DEFAULT 1,
    device_path TEXT    NOT NULL DEFAULT '/dev/bus/usb',
    frequency   BIGINT  NOT NULL DEFAULT 162550000,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT sdr_config_singleton CHECK (id = 1)
);

INSERT INTO sdr_config (id) VALUES (1);

-- +goose Down
DROP TABLE sdr_config;
DROP TABLE filter_config;
DROP TABLE fips_codes;
DROP TABLE event_codes;
DROP INDEX same_alerts_received_at_idx;
DROP TABLE same_alerts;
