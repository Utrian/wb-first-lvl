CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(19) PRIMARY KEY,
    track_number VARCHAR(14) UNIQUE NOT NULL,
    "entry" VARCHAR(4) NOT NULL,
    delivery JSONB NOT NULL,
    payment JSONB NOT NULL,
    items JSONB NOT NULL,
    locale VARCHAR(2) NOT NULL,
    internal_signature VARCHAR(55) DEFAULT '',
    customer_id VARCHAR(55) NOT NULL,
    delivery_service VARCHAR(55) NOT NULL,
    shardkey VARCHAR(55) NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    off_shard VARCHAR(55) NOT NULL
);
