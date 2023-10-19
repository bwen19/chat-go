-- name: InsertFriend :one
INSERT INTO friendships (
    user_id, friend_id, room_id, status
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: DeleteFriendByUser :many
DELETE FROM friendships
WHERE user_id = @id::bigint OR
    friend_id = @id::bigint
RETURNING room_id;

-- name: UpdateFriend :exec
UPDATE friendships
SET status = @status::varchar
WHERE user_id = @user_id::bigint AND friend_id = @friend_id::bigint;

-- name: UpdateAddFriend :one
UPDATE friendships
SET
    user_id = @user_id::bigint,
    friend_id = @friend_id::bigint,
    status = 'adding'
WHERE (
    user_id = @user_id::bigint AND friend_id = @friend_id::bigint
) OR (
    user_id = @friend_id::bigint AND friend_id = @user_id::bigint
) RETURNING *;

-- name: RetrieveFriend :one
SELECT * FROM friendships
WHERE (
    user_id = @user_id::bigint AND friend_id = @friend_id::bigint
) OR (
    user_id = @friend_id::bigint AND friend_id = @user_id::bigint
);

-- name: RetrieveFriendDetail :one
SELECT
    f.room_id, f.status, f.create_at, f.user_id, u1.username AS u_username,
    u1.nickname AS u_nickname, u1.avatar AS u_avatar, f.friend_id,
    u2.username AS f_username, u2.nickname AS f_nickname, u2.avatar AS f_avatar
FROM friendships AS f
INNER JOIN users AS u1 ON u1.id = f.user_id
INNER JOIN users AS u2 ON u2.id = f.friend_id
WHERE f.user_id = $1 AND f.friend_id = $2;

-- name: ListUserFriends :many
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

