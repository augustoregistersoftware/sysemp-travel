CREATE TABLE IF NOT EXISTS users (
    id_user SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE users
ADD CONSTRAINT users_email_unique UNIQUE (email);

ALTER TABLE users
ADD CONSTRAINT users_username_unique UNIQUE (username);

CREATE TABLE IF NOT EXISTS approved_users (
    id_approved_users SERIAL PRIMARY KEY,
    id_user INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE approved_users
ADD COLUMN negated BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE approved_users
ADD COLUMN email_user VARCHAR(255) NOT NULL DEFAULT '';