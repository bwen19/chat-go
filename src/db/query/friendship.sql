-- name: CreateFriend :exec
INSERT INTO friendships (
    user_id, friend_id, room_id, status
) VALUES (
    $1, $2, $3, $4
);

-- name: DeleteFriend :many
DELETE FROM friendships
WHERE (
    user_id = @id1::bigint AND friend_id = @id2::bigint
) OR (
    user_id = @id2::bigint AND friend_id = @id1::bigint
) RETURNING room_id;

-- name: DeleteFriendByUser :many
DELETE FROM friendships
WHERE user_id = @id::bigint OR
    friend_id = @id::bigint
RETURNING room_id;