CREATE TABLE customers_tokens (
    token TEXT not NULL UNIQUE,
    customers_id BIGINT NOT NULL REFERENCES customers,
    expire TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created TIMESTAMP not NULL DEFAULT CURRENT_TIMESTAMP
);


INSERT INTO customers_tokens (token, customers_id, expire, created)
VALUES (
    'token:text',
    'customers_id:bigint',
    'expire:timestamp without time zone',
    'created:timestamp without time zone'
  );

  select customer_id, expire from customers_tokens where token=$1