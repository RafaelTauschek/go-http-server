-- name: CreateRefreshToken :one
INSERT INTO refresh_token(token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES(
    $1,
    NOW(),
    NOW(),
    $2,
    (NOW() + interval '60 day'),
    NULL
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT * FROM refresh_token WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_token
SET updated_at = Now(), revoked_at = Now()
WHERE token = $1;