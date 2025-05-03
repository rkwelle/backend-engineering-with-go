CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    email CITEXT NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);