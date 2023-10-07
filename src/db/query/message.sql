-- name: CreateMessage :exec
INSERT INTO messages (
    room_id, sender_id, content, kind
) VALUES (
    $1, $2, $3, $4
);

-- name: DeleteMessageByUser :exec
DELETE FROM messages
WHERE
    sender_id = @user_id::bigint OR
    room_id = ANY(@room_ids::bigint[]);