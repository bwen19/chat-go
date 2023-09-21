-- name: CreateRoomMember :exec
INSERT INTO room_members (
    room_id, member_id, rank
) VALUES (
    $1, $2, $3
);

-- name: DeleteMember :exec
DELETE FROM room_members
WHERE room_id = $1 AND member_id = $2;

-- name: DeleteMemberByRoom :exec
DELETE FROM room_members
WHERE room_id = $1;

-- name: DeleteMemberByUser :exec
DELETE FROM room_members
WHERE
    member_id = @user_id::bigint OR
    room_id = ANY(@room_ids::bigint[]);