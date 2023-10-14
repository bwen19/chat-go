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

-- name: GetUserFriends :many
SELECT
    f.room_id, f.status, f.create_at, u.id, u.username,
    u.nickname, u.avatar, (f.user_id = $1) AS first
FROM friendships AS f
INNER JOIN users AS u ON u.id = f.friend_id
WHERE f.user_id = $1 AND status IN ('adding', 'accepted')
UNION
SELECT
    f.room_id, f.status, f.create_at, u.id, u.username,
    u.nickname, u.avatar, (f.user_id = $1) AS first
FROM friendships AS f
INNER JOIN users AS u ON u.id = f.user_id
WHERE f.friend_id = $1 AND status IN ('adding', 'accepted');

