-- name: GetUsers :many
SELECT * FROM users
ORDER BY id;

-- name: GetUserWithPasswordHash :one
SELECT
    u.id,
    u.name,
    u.email,
    c.password_hash
FROM users u
JOIN user_credentials c ON c.user_id = u.id
WHERE u.email = $1;

-- name: CreateUserWithPasswordHash :one
WITH new_user AS (
        INSERT INTO users (
            name,
            email
        ) VALUES (
            $1,
            $2
        )
        RETURNING *
    ),
    new_credentials AS (
        INSERT INTO user_credentials (
            user_id,
            password_hash
        )
        SELECT id, $3
        FROM new_user
    )
SELECT *
FROM new_user;
