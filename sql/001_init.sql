-- buat database (jalankan sekali di psql/pgAdmin)
-- CREATE DATABASE bioskop_db;

-- jalankan di database bioskop_db
CREATE TABLE IF NOT EXISTS bioskop (
  id     SERIAL PRIMARY KEY,
  nama   TEXT   NOT NULL,
  lokasi TEXT   NOT NULL,
  rating REAL   NOT NULL DEFAULT 0
);
