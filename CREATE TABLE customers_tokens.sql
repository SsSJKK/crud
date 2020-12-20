CREATE TABLE managers_tokens (
    token TEXT not NULL UNIQUE,
    manager_id BIGINT NOT NULL REFERENCES managers,
    expire TIMESTAMP not null DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created TIMESTAMP not NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO customers_tokens (token, customers_id, expire, created)
VALUES (
        'token:text',
        'customer_id:bigint',
        'expire:timestamp without time zone',
        'created:timestamp without time zone'
    );
select customer_id,
    expire
from customers_tokens
where token = $1 drop TABLE customers_tokens CREATE TABLE products (
        id BIGSERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        price INTEGER NOT NULL CHECK(price > 0),
        qty INTEGER NOT NULL DEFAULT 0 CHECK (qty >= 0),
        active BOOLEAN NOT NULL DEFAULT TRUE,
        created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );
DROP table managers;
CREATE TABLE managers (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT not NULL UNIQUE,
    password text not NULL,
    roles text [] NOT NULL DEFAULT '{}',
    active BOOLEAN not NULL DEFAULT TRUE,
    creatred TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

UPDATE products
SET name = 'ipp', price = 200, qty = 3 where id = 4

select  * from products


CREATE TABLE sales 
(
    id BIGSERIAL PRIMARY KEY,
    manager_id BIGINT NOT NULL REFERENCES managers,
    customer_id BIGINT NOT NULL DEFAULT 0,
    crated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sale_positions 
(
    id BIGSERIAL PRIMARY KEY,
    sale_id BIGINT NOT NULL REFERENCES sales,
    product_id BIGINT NOT NULL REFERENCES products,
    price INTEGER NOT NULL CHECK (price >= 0),
    qty INTEGER NOT NULL CHECK (qty > 0),
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



TRUNCATE TABLE sales CASCADE