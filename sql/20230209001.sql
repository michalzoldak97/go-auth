CREATE SCHEMA auth;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE EXTENSION IF NOT EXISTS btree_gin;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE auth.tbl_user(
    user_id uuid DEFAULT uuid_generate_v4 (),
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    password VARCHAR(63) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    deactivated_at TIMESTAMP,
    PRIMARY KEY (user_id)
);


CREATE INDEX idx_auth_u_email ON auth.tbl_user USING GIN (email);

CREATE INDEX idx_auth_u_pass ON auth.tbl_user USING GIN (password);
