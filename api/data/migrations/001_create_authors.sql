CREATE TABLE IF NOT EXISTS authors (
    key             VARCHAR PRIMARY KEY,
    name            VARCHAR,
    type_key        VARCHAR,
    source_records  VARCHAR[],
    latest_revision INTEGER,
    revision        INTEGER,
    created         TIMESTAMP,
    last_modified   TIMESTAMP
);