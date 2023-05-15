CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(19) PRIMARY KEY,
    track_number VARCHAR(14) UNIQUE NOT NULL,
    "entry" VARCHAR(4) NOT NULL,
    delivery JSONB NOT NULL,
    payment JSONB NOT NULL,
    locale VARCHAR(2) NOT NULL,
    internal_signature VARCHAR(55) DEFAULT '',
    customer_id VARCHAR(55) NOT NULL,
    delivery_service VARCHAR(55) NOT NULL,
    shardkey VARCHAR(55) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    off_shard VARCHAR(55) NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    order_id VARCHAR(19) REFERENCES orders (order_uid),
    chrt_id BIGINT NOT NULL,
    track_number VARCHAR(14) NOT NULL,
    price INT NOT NULL,
    rid VARCHAR(25) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    sale INT NOT NULL,
    size VARCHAR(55) DEFAULT '0',
    total_price INT NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(255) NOT NULL,
    "status" SMALLINT NOT NULL
);
