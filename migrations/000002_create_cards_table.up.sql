CREATE TABLE cards (
    id BIGSERIAL PRIMARY KEY,
    number VARCHAR(16) NOT NULL,
    data VARCHAR(5) NOT NULL,
    cvv VARCHAR(3) NOT NULL,
    user_id BIGINT NOT NULL,
    FOREIGN KEY(users_id) REFERENCES users (id)
)