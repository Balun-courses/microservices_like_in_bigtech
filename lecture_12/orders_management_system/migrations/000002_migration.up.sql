CREATE TABLE IF NOT EXISTS orders (
    id uuid PRIMARY KEY,
    user_id int8 NOT NULL,
    items json,
    delivery_variant_id int8,
    delivery_date TIMESTAMP
);
