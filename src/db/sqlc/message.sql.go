// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: message.sql

package sqlc

import (
	"context"
	"time"
)

const deleteMessagesByRoom = `-- name: DeleteMessagesByRoom :exec
DELETE FROM messages
WHERE room_id = $1
`

func (q *Queries) DeleteMessagesByRoom(ctx context.Context, roomID int64) error {
	_, err := q.db.Exec(ctx, deleteMessagesByRoom, roomID)
	return err
}

const deleteMessagesByUser = `-- name: DeleteMessagesByUser :exec
DELETE FROM messages
WHERE
  sender_id = $1::bigint OR
  room_id = ANY($2::bigint[])
`

type DeleteMessagesByUserParams struct {
	UserID  int64   `json:"user_id"`
	RoomIds []int64 `json:"room_ids"`
}

func (q *Queries) DeleteMessagesByUser(ctx context.Context, arg *DeleteMessagesByUserParams) error {
	_, err := q.db.Exec(ctx, deleteMessagesByUser, arg.UserID, arg.RoomIds)
	return err
}

const insertMessage = `-- name: InsertMessage :exec
INSERT INTO messages (
    room_id, sender_id, content, kind, send_at
  )
VALUES (
    $1, $2, $3, $4, $5
  )
`

type InsertMessageParams struct {
	RoomID   int64     `json:"room_id"`
	SenderID int64     `json:"sender_id"`
	Content  string    `json:"content"`
	Kind     string    `json:"kind"`
	SendAt   time.Time `json:"send_at"`
}

func (q *Queries) InsertMessage(ctx context.Context, arg *InsertMessageParams) error {
	_, err := q.db.Exec(ctx, insertMessage,
		arg.RoomID,
		arg.SenderID,
		arg.Content,
		arg.Kind,
		arg.SendAt,
	)
	return err
}

const retrieveMessages = `-- name: RetrieveMessages :many
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
LIMIT $2 OFFSET $3
`

type RetrieveMessagesParams struct {
	RoomID int64 `json:"room_id"`
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type RetrieveMessagesRow struct {
	ID       int64     `json:"id"`
	RoomID   int64     `json:"room_id"`
	SenderID int64     `json:"sender_id"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar"`
	Content  string    `json:"content"`
	Kind     string    `json:"kind"`
	SendAt   time.Time `json:"send_at"`
}

func (q *Queries) RetrieveMessages(ctx context.Context, arg *RetrieveMessagesParams) ([]*RetrieveMessagesRow, error) {
	rows, err := q.db.Query(ctx, retrieveMessages, arg.RoomID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*RetrieveMessagesRow{}
	for rows.Next() {
		var i RetrieveMessagesRow
		if err := rows.Scan(
			&i.ID,
			&i.RoomID,
			&i.SenderID,
			&i.Name,
			&i.Avatar,
			&i.Content,
			&i.Kind,
			&i.SendAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
