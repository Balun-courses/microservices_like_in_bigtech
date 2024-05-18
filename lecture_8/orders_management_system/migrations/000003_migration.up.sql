CREATE TABLE IF NOT EXISTS orders_outbox_messages (
    id serial PRIMARY KEY,
    order_id uuid NOT NULL
);
