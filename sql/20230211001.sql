CREATE TABLE auth.tbl_token(
    token_id uuid DEFAULT uuid_generate_v4 (),
    email VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    token_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (token_id)
);

CREATE INDEX idx_auth_t_token ON auth.tbl_token USING GIN (token);

CREATE TABLE auth.tbl_user_token(
    user_token_id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid NOT NULL,
    token_id uuid NOT NULL,
    deactivated_at TIMESTAMP,
    PRIMARY KEY (token_id),
    CONSTRAINT fk_user_token_user_id
      FOREIGN KEY (user_id)
        REFERENCES auth.tbl_user(user_id)
        ON UPDATE 
            CASCADE 
        ON DELETE 
            CASCADE,
    CONSTRAINT fk_user_token_token_id
      FOREIGN KEY (token_id)
        REFERENCES auth.tbl_token(token_id),
    UNIQUE (user_id, token_id)
);
