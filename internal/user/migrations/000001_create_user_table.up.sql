CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE UNIQUE INDEX udx_users_email ON users (email) WHERE is_active = TRUE;
CREATE UNIQUE INDEX udx_users_phone ON users (phone) WHERE is_active = TRUE;

