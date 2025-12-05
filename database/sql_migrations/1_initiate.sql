-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE IF NOT EXISTS bioskop (
    id     SERIAL PRIMARY KEY,
    nama   VARCHAR(256) NOT NULL,
    lokasi VARCHAR(256) NOT NULL,
    rating REAL DEFAULT 0
);

-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin

DROP TABLE IF EXISTS bioskop;

-- +migrate StatementEnd
