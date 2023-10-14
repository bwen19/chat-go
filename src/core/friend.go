package core

import (
	"context"
	"time"
)

type FriendInfo struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Nickname string    `json:"nickname"`
	Avatar   string    `json:"avatar"`
	Status   string    `json:"status"`
	RoomID   int64     `json:"room_id"`
	First    bool      `json:"first"`
	CreateAt time.Time `json:"create_at"`
}

func (s *State) GetUserFriends(ctx context.Context, userID int64) ([]*FriendInfo, error) {
	rows, err := s.Store.GetUserFriends(ctx, userID)
	if err != nil {
		return nil, err
	}

	rsp := make([]*FriendInfo, 0, len(rows))
	for _, row := range rows {
		friend := &FriendInfo{
			ID:       row.ID,
			Username: row.Username,
			Nickname: row.Nickname,
			Avatar:   row.Avatar,
			Status:   row.Status,
			RoomID:   row.RoomID,
			First:    row.First,
			CreateAt: row.CreateAt,
		}
		rsp = append(rsp, friend)
	}

	return rsp, nil
}
