-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password_hash)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND refresh_tokens.expires_at > NOW()
AND refresh_tokens.revoked_at IS NULL;

-- name: UpdateUser :one
UPDATE users
SET email = $2, password_hash = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, email, created_at, updated_at, is_chirpy_red;

-- name: UserAddChirpyRed :one
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1
RETURNING id, is_chirpy_red;

-- name: UserRemoveChirpyRed :one
UPDATE users
SET is_chirpy_red = FALSE
WHERE id = $1
RETURNING id, is_chirpy_red;
