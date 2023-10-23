-- name: InsertSession :one
INSERT INTO sessions (
    id, user_id, refresh_token,
    client_ip, user_agent, expire_at
  )
VALUES (
    $1, $2, $3, $4, $5, $6
  )
RETURNING *;

-- name: RetrieveSession :one
SELECT * FROM sessions
WHERE id = @id::uuid LIMIT 1;

-- name: RetrieveSessions :many
SELECT
  id,
  client_ip,
  user_agent,
  expire_at,
  create_at,
  count(*) OVER() AS total
FROM sessions
WHERE user_id = @user_id::bigint
ORDER BY create_at
LIMIT $1 OFFSET $2;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = @id::uuid;

-- name: DeleteSessionsByUser :exec
DELETE FROM sessions
WHERE user_id = @user_id::bigint;