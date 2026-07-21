CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    short_code VARCHAR(10) NOT NULL UNIQUE,

    original_url TEXT NOT NULL,

    user_id UUID NOT NULL,

    click_count INT NOT NULL DEFAULT 0,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_urls_users
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);