CREATE TABLE IF NOT EXISTS swift_codes (
    id SERIAL PRIMARY KEY,
    country_iso2 CHAR(2) NOT NULL CHECK (LENGTH(country_iso2) = 2),
    swift_code VARCHAR(11) NOT NULL CHECK (LENGTH(swift_code) IN (8, 11)),
    code_type VARCHAR(5) NOT NULL,
    name TEXT NOT NULL,
    address TEXT,
    town_name VARCHAR(60),
    country_name VARCHAR(50),
    time_zone VARCHAR(50),
    UNIQUE (swift_code)
);

CREATE UNIQUE INDEX swift_code_idx ON swift_codes (swift_code);
CREATE INDEX country_iso2_idx ON swift_codes (country_iso2);

COPY swift_codes (country_iso2, swift_code, code_type, name, address, town_name, country_name, time_zone)
FROM '/docker-entrypoint-initdb.d/data.csv'
DELIMITER ','
CSV HEADER;
