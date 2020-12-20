INSERT INTO customers_tokens (token, customer_id, expire, created)
VALUES (
        'token:text',
        'customer_id:bigint',
        'expire:timestamp without time zone',
        'created:timestamp without time zone'
    ) ;

INSERT INTO managers (
    name,
    phone,
    password,
    roles
  )
VALUES (
    'vasya',
    '+992000000001',
    '$2a$10$blk3/GSurOOfeisjL6R0WeO3M1GTHWVea51Wc4lffetNrOZT8xiDK',
    '{MANAGER, ADMIN}'
  );

TRUNCATE TABLE managers CASCADE;


INSERT INTO sales (manager_id)
VALUES (
    14
  );



SELECT s.manager_id, sum(sp.qty * sp.price) FROM sales s
JOIN sale_positions as sp on sp.sale_id = s.id
where s.manager_id = 14
GROUP by s.manager_id;

select * from managers WHERE roles[2] = 'ADMIN' and id = 14

UPDATE products
SET
qty = 0
WHERE id=1