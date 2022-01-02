-- in database robozdb

DROP DATABASE IF EXISTS postgres;

CREATE SCHEMA roboz;

SET SCHEMA 'roboz';

DROP SCHEMA IF EXISTS public;

ALTER DATABASE robozdb SET search_path TO roboz;

CREATE TABLE IF NOT EXISTS executions (
  "id" SERIAL PRIMARY KEY,
  "timestamp" timestamptz NOT NULL, -- would be better to avoid reserved words
  commands integer CONSTRAINT positive_commands CHECK (commands > 0) NOT NULL, -- specification says commmands
  "result" integer CONSTRAINT positive_or_zero_result CHECK (result >= 0) NOT NULL, -- would be better to avoid reserved words
  duration double precision CONSTRAINT positive_or_zero_duration CHECK (duration >= 0.0) NOT NULL -- for go float64 - a bit extravagant
);

/*
--  TODO
--  - application users with retstricted rights, rather than using superuser
--  - find out how to source new user password from environment

--  @see https://github.com/docker-library/postgres/issues/130
--  @see https://graspingtech.com/docker-compose-postgresql/

CREATE USER robozreadonly WITH LOGIN PASSWORD 'changeme';
GRANT USAGE ON SCHEMA roboz TO robozreadonly;
GRANT SELECT ON executions TO robozreadonly;
GRANT SELECT ON executions_id_seq TO robozreadonly;

CREATE USER robozupdate WITH LOGIN PASSWORD 'changeme';
GRANT USAGE ON SCHEMA roboz TO robozupdate;
GRANT SELECT, INSERT ON ALL TABLES IN SCHEMA roboz TO robozupdate;
GRANT SELECT, USAGE, UPDATE ON ALL SEQUENCES IN SCHEMA roboz TO robozupdate;
*/