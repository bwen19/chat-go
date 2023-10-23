-- name: InsertMember :exec
INSERT INTO room_members (
    room_id, member_id, rank
  )
VALUES (
    $1, $2, $3
  );

-- name: RetrieveMember :one
SELECT * FROM room_members
WHERE
  room_id = $1
  AND member_id = $2
LIMIT 1;

-- name: RetrieveMembers :many
SELECT
  u.id,
  u.nickname AS name,
  u.avatar,
  m.rank,
  m.join_at
FROM
  room_members AS m
  JOIN users AS u
    ON m.member_id = u.id
WHERE m.room_id = @room_id::bigint;

-- name: RetrieveOwnerRoomIDs :many
SELECT room_id
FROM room_members
WHERE
  member_id = @user_id::bigint
  AND rank = 'owner';

-- name: DeleteMember :exec
DELETE FROM room_members
WHERE
  room_id = $1
  AND member_id = $2;

-- name: DeleteMembers :many
DELETE FROM room_members
WHERE
  room_id = @room_id::bigint
  AND member_id = ANY(@member_ids::bigint[])
  AND rank <> 'owner'
RETURNING member_id;

-- name: DeleteMembersByRoom :exec
DELETE FROM room_members
WHERE room_id = $1;

-- name: DeleteMembersByUser :exec
DELETE FROM room_members
WHERE
    member_id = @user_id::bigint OR
    room_id = ANY(@room_ids::bigint[]);