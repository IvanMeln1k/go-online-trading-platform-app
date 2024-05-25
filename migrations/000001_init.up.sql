CREATE TYPE roles AS ENUM ('user', 'seller', 'manager');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(1024) NOT NULL,
    name VARCHAR(1024) NOT NULL,
    email VARCHAR(1024) NOT NULL,
    hash_password VARCHAR(1024) NOT NULL,
    role roles,
    email_verified BOOL DEFAULT false
);

CREATE TABLE cards (
    id BIGSERIAL PRIMARY KEY,
    number VARCHAR(16) NOT NULL,
    data VARCHAR(5) NOT NULL,
    cvv VARCHAR(3) NOT NULL,
    user_id BIGINT NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE sellers (
    user_id BIGINT UNIQUE,
    name VARCHAR(1024),
    confirmed BOOL DEFAULT false,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE accounts (
    seller_id BIGINT UNIQUE,
    number VARCHAR(16) NOT NULL,
    bank VARCHAR(1024) NOT NULL,
    FOREIGN KEY (seller_id) REFERENCES sellers (user_id) ON DELETE CASCADE
);

CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    lat FLOAT NOT NULL,
    lon FLOAT NOT NULL,
    adress TEXT NOT NULL
)

CREATE TYPE place_type AS ENUM ('warehouse', 'point_of_issue', 'point_of_path');

CREATE TABLE places (
    id BIGSERIAL PRIMARY KEY,
    address_id BIGINT,
    place_type place_type,
    FOREIGN KEY (address_id) REFERENCES addresses (id) ON DELETE SET NULL
);

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    article VARCHAR(1024) NOT NULL,
    name VARCHAR(1024) NOT NULL,
    price INT NOT NULL,
    manufacturer VARCHAR(1024),
    seller_id BIGINT NOT NULL,
    deleted BOOL NOT NULL DEFAULT false,
    FOREIGN KEY (seller_id) REFERENCES sellers (user_id) ON DELETE CASCADE
);

CREATE TABLE products_warehouses (
    product_id BIGINT,
    warehouse_id BIGINT,
    PRIMARY KEY (product_id, warehouse_id),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (warehouse_id) REFERENCES places(id) ON DELETE CASCADE
);


CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGIN NOT NULL,
    product_id BIGINT NOT NULL,
    received BOOL NOT NULL DEFAULT false,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE TYPE shipping_status AS ENUM ('on_the_way', 'at_the_point', 'at_the_point_of_issue');

CREATE TABLE shipping_information (
    order_id BIGINT UNIQUE,
    place_id BIGINT,
    place_of_issue_id BIGINT NOT NULL,
    status shipping_status,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    FOREIGN KEY (place_id) REFERENCES places (id) ON DELETE SET NULL,
    FOREIGN KEY (place_of_issue_id) REFERENCES places (id) ON DELETE SET NULL
);

CREATE TABLE order_info (
    order_id BIGINT UNIQUE,
    quantity INT NOT NULL,
    paidFor BOOL NOT NULL,
    amount INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE
);