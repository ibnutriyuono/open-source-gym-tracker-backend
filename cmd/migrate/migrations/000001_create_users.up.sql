CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_banned BOOLEAN NOT NULL DEFAULT FALSE,
    banned_at TIMESTAMPTZ NULL,
    avatar VARCHAR(255) NULL,
    phone_number VARCHAR(20) NULL,
    year_of_birth CHAR(4) NULL,
    month_of_birth CHAR(2) NULL,
    date_of_birth CHAR(2) NULL,
    gender VARCHAR(20) NULL CHECK (gender IS NULL OR gender IN ('male', 'female')),
    verify_at TIMESTAMPTZ NULL
);
