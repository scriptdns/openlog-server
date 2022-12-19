-- Copy this into your psql to create the structure

CREATE DATABASE openlog;

\connect openlog

CREATE TABLE logs(
    "time" timestamp NOT NULL,
    stream text NOT NULL,
    data jsonb NOT NULL 
)

CREATE INDEX logs_stream_ix ON logs (stream, time DESC);

select create_hypertable('logs','time','stream',8);

--CREATE TABLE stream(
-- To keep some information about each stream?
--)