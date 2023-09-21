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
