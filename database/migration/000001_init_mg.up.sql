CREATE TABLE
    users (
        id SERIAL PRIMARY KEY,
        email VARCHAR(254) UNIQUE NOT NULL,
        password_hash VARCHAR(72) NOT NULL,
        created_at DATE DEFAULT NOW (),
        updated_at DATE
    );