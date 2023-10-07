-- name: CreateRoom :one
INSERT INTO rooms (
    name, cover, category
) VALUES (
    $1, $2, $3
) RETURNING
    id, name, cover, category, create_at;

-- name: DeleteRooms :exec
DELETE FROM rooms
WHERE id = ANY(@room_ids::bigint[]);

-- name: GetUserRooms :many
WITH rooms_cte AS (
    SELECT id AS room_id, name, cover, category, create_at
    FROM rooms WHERE id IN (
        SELECT room_id FROM room_members AS m
        WHERE m.member_id = $1)
)
SELECT room_id, name, cover, category, create_at,
    member_id, rank, join_at, nickname, avatar
FROM rooms_cte AS r,
    LATERAL (
        SELECT member_id, rank, join_at, nickname, avatar
        FROM room_members AS y
        INNER JOIN users AS u ON y.member_id = u.id
        WHERE y.room_id = r.room_id
    ) AS m;