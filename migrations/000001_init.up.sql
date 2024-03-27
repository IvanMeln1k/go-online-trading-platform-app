CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(1024) NOT NULL,
    name VARCHAR(1024) NOT NULL,
    email VARCHAR(1024) NOT NULL,
    hash_password VARCHAR(1024) NOT NULL,
    email_verified BOOL DEFAULT false
);