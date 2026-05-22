CREATE TABLE categories (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    inactive BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE UNIQUE INDEX udx_categories_name ON categories (name) WHERE inactive = FALSE;