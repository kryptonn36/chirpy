-- name: CreateRefreshToken :one
INSERT INTO refreshTokens(token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NULL,
    $3
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refreshTokens
WHERE token=$1;

-- name: RevoketimeUpdate :exec
UPDATE refreshTokens
SET revoked_at=NOW(), updated_at=NOW()
WHERE token=$1;