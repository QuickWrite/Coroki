-- name: GetUsers :many
SELECT * FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
    name, email
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetUserWithPasswordHash :one
SELECT
    u.id,
    u.name,
    u.email,
    c.password_hash
FROM users u
JOIN user_credentials c ON c.user_id = u.id
WHERE u.email = $1;
