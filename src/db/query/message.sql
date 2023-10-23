-- name: InsertMessage :exec
INSERT INTO messages (
    room_id, sender_id, content, kind, send_at
  )
VALUES (
    $1, $2, $3, $4, $5
  );

-- name: RetrieveMessages :many
SELECT
  m.id,
  m.room_id,
  m.sender_id,
  u.nickname AS name,
  u.avatar,
  m.content,
  m.kind,
  m.send_at
FROM
  messages AS m
  JOIN users AS u
    ON u.id = m.sender_id
  JOIN room_members AS r
    ON r.room_id = m.room_id
    AND r.member_id = m.sender_id
WHERE
  m.room_id = $1
  AND m.send_at > r.join_at
ORDER BY m.send_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteMessagesByRoom :exec
DELETE FROM messages
WHERE room_id = $1;

-- name: DeleteMessagesByUser :exec
DELETE FROM messages
WHERE
  sender_id = @user_id::bigint OR
  room_id = ANY(@room_ids::bigint[]);