-- name: InsertUser :one
INSERT INTO users (
    username, hashed_password, nickname,
    avatar, role, room_id
	)
VALUES (
    $1, $2, $3, $4, $5, $6
	)
RETURNING *;

-- name: RetrieveUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: RetrieveUserByName :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: RetrieveUsers :many
SELECT
	id,
	username,
	avatar,
	nickname,
	role,
	deleted,
	create_at,
	count(*) OVER() AS total
FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
	username = COALESCE(sqlc.narg(username), username),
	hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password),
	avatar = COALESCE(sqlc.narg(avatar), avatar),
	nickname = COALESCE(sqlc.narg(nickname), nickname),
	role = COALESCE(sqlc.narg(role), role),
	deleted = COALESCE(sqlc.narg(deleted), deleted)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING room_id;
