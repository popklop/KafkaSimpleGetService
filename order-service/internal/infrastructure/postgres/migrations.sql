CREATE TABLE IF NOT EXISTS orders (
    id TEXT PRIMARY KEY,
    track_number TEXT NOT NULL,
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INTEGER,
    date_created TIMESTAMP NOT NULL DEFAULT NOW(),
    oof_shard TEXT
);


CREATE TABLE IF NOT EXISTS delivery (
    order_id TEXT PRIMARY KEY REFERENCES orders(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    phone TEXT,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT
);


CREATE TABLE IF NOT EXISTS payments (
    order_id TEXT PRIMARY KEY REFERENCES orders(id) ON DELETE CASCADE,
    transaction TEXT NOT NULL,
    request_id TEXT,
    currency TEXT,
    provider TEXT,
    amount INTEGER NOT NULL,
    payment_dt TIMESTAMP,
    bank TEXT,
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);


CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_id TEXT REFERENCES orders(id) ON DELETE CASCADE,
    chrt_id INTEGER,
    track_number TEXT,
    price INTEGER,
    rid TEXT,
    name TEXT,
    sale INTEGER,
    size TEXT,
    nm_id INTEGER,
    brand TEXT,
    status INTEGER
);


CREATE INDEX IF NOT EXISTS idx_items_order ON items(order_id);
CREATE INDEX IF NOT EXISTS idx_orders_customer ON orders(customer_id);