-- name: CreateCustomer :one
INSERT INTO customers (name, email)
VALUES ($1, $2)
RETURNING id, name, email;

-- name: GetCustomerByEmail :one
SELECT id, name, email
FROM customers
WHERE email = $1;
