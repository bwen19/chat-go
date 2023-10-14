// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: friendship.sql

package db

import (
	"context"
	"time"
)

const createFriend = `-- name: CreateFriend :exec
INSERT INTO friendships (
    user_id, friend_id, room_id, status
) VALUES (
    $1, $2, $3, $4
)
`

type CreateFriendParams struct {
	UserID   int64  `json:"user_id"`
	FriendID int64  `json:"friend_id"`
	RoomID   int64  `json:"room_id"`
	Status   string `json:"status"`
}

func (q *Queries) CreateFriend(ctx context.Context, arg CreateFriendParams) error {
	_, err := q.db.Exec(ctx, createFriend,
		arg.UserID,
		arg.FriendID,
		arg.RoomID,
		arg.Status,
	)
	return err
}

const deleteFriend = `-- name: DeleteFriend :many
DELETE FROM friendships
WHERE (
    user_id = $1::bigint AND friend_id = $2::bigint
) OR (
    user_id = $2::bigint AND friend_id = $1::bigint
) RETURNING room_id
`

type DeleteFriendParams struct {
	Id1 int64 `json:"id1"`
	Id2 int64 `json:"id2"`
}

func (q *Queries) DeleteFriend(ctx context.Context, arg DeleteFriendParams) ([]int64, error) {
	rows, err := q.db.Query(ctx, deleteFriend, arg.Id1, arg.Id2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var room_id int64
		if err := rows.Scan(&room_id); err != nil {
			return nil, err
		}
		items = append(items, room_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteFriendByUser = `-- name: DeleteFriendByUser :many
DELETE FROM friendships
WHERE user_id = $1::bigint OR
    friend_id = $1::bigint
RETURNING room_id
`

func (q *Queries) DeleteFriendByUser(ctx context.Context, id int64) ([]int64, error) {
	rows, err := q.db.Query(ctx, deleteFriendByUser, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var room_id int64
		if err := rows.Scan(&room_id); err != nil {
			return nil, err
		}
		items = append(items, room_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserFriends = `-- name: GetUserFriends :many
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
WHERE f.friend_id = $1 AND status IN ('adding', 'accepted')
`

type GetUserFriendsRow struct {
	RoomID   int64     `json:"room_id"`
	Status   string    `json:"status"`
	CreateAt time.Time `json:"create_at"`
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Nickname string    `json:"nickname"`
	Avatar   string    `json:"avatar"`
	First    bool      `json:"first"`
}

func (q *Queries) GetUserFriends(ctx context.Context, userID int64) ([]GetUserFriendsRow, error) {
	rows, err := q.db.Query(ctx, getUserFriends, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetUserFriendsRow{}
	for rows.Next() {
		var i GetUserFriendsRow
		if err := rows.Scan(
			&i.RoomID,
			&i.Status,
			&i.CreateAt,
			&i.ID,
			&i.Username,
			&i.Nickname,
			&i.Avatar,
			&i.First,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
