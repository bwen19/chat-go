-- name: InsertRoom :one
INSERT INTO rooms (
    name, cover, category
  )
VALUES (
    $1, $2, $3
  )
RETURNING *;

-- name: RetrieveRoom :one
SELECT * FROM rooms
WHERE id = $1 LIMIT 1;

-- name: RetrieveUserRooms :many
WITH rooms_cte AS (
  SELECT
    id AS room_id,
    name,
    cover,
    category,
    create_at
  FROM rooms
  WHERE id IN (
    SELECT room_id
    FROM room_members AS m
    WHERE m.member_id = $1
  )
)
SELECT
  room_id,
  name,
  cover,
  category,
  create_at,
  member_id,
  rank,
  join_at,
  nickname,
  avatar
FROM rooms_cte AS r,
  LATERAL (
    SELECT
      member_id,
      rank,
      join_at,
      nickname,
      avatar
    FROM
      room_members AS y
      JOIN users AS u
        ON y.member_id = u.id
    WHERE y.room_id = r.room_id
  ) AS m;

-- name: RetrieveFriendRooms :many
SELECT
  m.room_id,
  m.rank,
  m.join_at,
  r.category,
  r.create_at,
  m.member_id,
  u.nickname,
  u.avatar
FROM
  room_members AS m
  JOIN rooms AS r
    ON r.id = m.room_id
  JOIN users AS u
    ON u.id = m.member_id
WHERE m.room_id = $1;

-- name: UpdateRoom :one
UPDATE rooms
SET
  name = COALESCE(sqlc.narg(name), name),
  cover = COALESCE(sqlc.narg(cover), cover)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM rooms
WHERE id = $1;

-- name: DeleteRooms :exec
DELETE FROM rooms
WHERE id = ANY(@room_ids::bigint[]);