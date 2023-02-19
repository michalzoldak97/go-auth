CREATE TABLE auth.tbl_config(
    config_id SERIAL PRIMARY KEY,
    security_config_jsonb JSONB NOT NULL,
    is_active BOOL NOT NULL
);