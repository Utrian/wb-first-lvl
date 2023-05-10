CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(19) PRIMARY KEY,
    track_number VARCHAR(14) UNIQUE NOT NULL,
    "entry" VARCHAR(4) NOT NULL,
    locale VARCHAR(2) NOT NULL,
    internal_signature VARCHAR(55) DEFAULT '',
    customer_id VARCHAR(55) NOT NULL,
    delivery_service VARCHAR(55) NOT NULL,
    shardkey VARCHAR(55) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP(),
    off_shard VARCHAR(55) NOT NULL
);

CREATE TABLE IF NOT EXISTS deliveries (
    order_uid VARCHAR(19) REFERENCES orders(order_uid),
    "name" VARCHAR(255) NOT NULL,
    phone VARCHAR(12) NOT NULL,
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(255) NOT NULL,
    "address" VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
    transaction VARCHAR(19) REFERENCES orders(order_uid),
    request_id VARCHAR(19) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    "provider" VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    payment_dt INT NOT NULL,
    bank VARCHAR(255) NOT NULL,
    delivery_cost INT NOT NULL,
    goods_total INT NOT NULL,
    custom_fee INT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    chrt_id BIGINT PRIMARY KEY,
    track_number VARCHAR(55) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    sale INT DEFAULT 0,
    size VARCHAR(255) DEFAULT '0',
    total_price INT NOT NULL,
    nm_id INT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    "status" SMALLINT NOT NULL
);
