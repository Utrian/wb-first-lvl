CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR PRIMARY KEY,
    track_number VARCHAR UNIQUE NOT NULL,
    entry VARCHAR NOT NULL,
    locale VARCHAR(2) NOT NULL,
    internal_signature VARCHAR,
    customer_id VARCHAR NOT NULL,
    delivery_service VARCHAR NOT NULL,
    shardkey VARCHAR NOT NULL,
    sm_id INT NOT NULL,
    date_created TIMESTAMP NOT NULL, -- Возможно стоит поставить дефолтное NOW()
    off_shard VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS deliveries (
    order_uid VARCHAR REFERENCES orders(order_uid),
    name VARCHAR NOT NULL,
    phone VARCHAR NOT NULL,
    zip VARCHAR(10) NOT NULL,
    city VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS payments (
    transaction VARCHAR REFERENCES orders(order_uid),
    request_id VARCHAR,
    currency VARCHAR(3) NOT NULL,
    provider VARCHAR NOT NULL,
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
    price INT, -- думаю цена может быть NULL, например, если товар закончился или еще не поступил в продажу;
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(255) DEFAULT
    total_price
    nm_id
    brand
    status
);