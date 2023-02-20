CREATE SCHEMA log;

CREATE TABLE log.tbl_auth_user_login_attempt (
    auth_user_login_attempt_id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255),
    message VARCHAR(255),
    success BOOL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL
);