CREATE TABLE "users" (
                        id SERIAL PRIMARY KEY,
                        username VARCHAR(50) NOT NULL,
                        email VARCHAR(100) UNIQUE NOT NULL,
                        password_hash VARCHAR(255) NOT NULL,
                        is_confirmed BOOLEAN DEFAULT false
);

CREATE TABLE email_confirmations (
                        id SERIAL PRIMARY KEY,
                        user_id INT REFERENCES "users"(id),
                        code TEXT NOT NULL,
                        expires_at TIMESTAMPTZ NOT NULL
);