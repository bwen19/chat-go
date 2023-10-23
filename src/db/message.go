package db

import (
	"context"
	"gochat/src/db/sqlc"
	"sort"
	"time"
)

func (s *dbStore) GetUserRooms(ctx context.Context, userID int64, endTime time.Time) ([]*RoomInfo, error) {
	rows, err := s.RetrieveUserRooms(ctx, userID)
	if err != nil {
		return nil, err
	}

	endTimeMilli := endTime.UnixMilli()
	roomMap := make(map[int64]*RoomInfo)
	for _, row := range rows {
		member := &MemberInfo{
			ID:     row.MemberID,
			Name:   row.Nickname,
			Avatar: row.Avatar,
			JoinAt: row.JoinAt,
		}

		isPrivateRoom := row.Category == "private" && row.MemberID != userID
		if room, ok := roomMap[row.RoomID]; ok {
			if isPrivateRoom {
				room.Name = row.Nickname
				room.Cover = row.Avatar
			}
			room.Members = append(room.Members, member)
			continue
		}

		messages, err := s.getRoomMessages(ctx, row.RoomID, endTimeMilli)
		if err != nil {
			return nil, err
		}

		room := &RoomInfo{
			ID:       row.RoomID,
			Name:     row.Name,
			Cover:    row.Cover,
			Category: row.Category,
			CreateAt: row.CreateAt,
			Members:  []*MemberInfo{member},
			Messages: messages,
		}
		if isPrivateRoom {
			room.Name = row.Nickname
			room.Cover = row.Avatar
		}

		roomMap[row.RoomID] = room
	}

	rsp := make([]*RoomInfo, 0, len(roomMap))
	for _, room := range roomMap {
		rsp = append(rsp, room)
	}
	sort.Sort(RoomSlice(rsp))

	return rsp, nil
}

func (s *dbStore) getRoomMessages(ctx context.Context, roomID int64, endTime int64) ([]*MessageInfo, error) {
	messages := make([]*MessageInfo, 0)

	if err := s.GetCacheList(ctx, roomKey(roomID), &messages); err != nil {
		return nil, err
	}

	var offset int64
	for _, message := range messages {
		curOffset := (endTime - message.SendAt.UnixMilli()) / 86400000
		if curOffset != offset {
			message.Divide = true
			offset = curOffset
		}
	}

	return messages, nil
}

func (s *dbStore) GetUserFriends(ctx context.Context, userID int64) ([]*FriendInfo, error) {
	rows, err := s.RetrieveUserFriends(ctx, userID)
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

func (s *dbStore) CacheMessage(ctx context.Context, msg *MessageInfo) error {
	return s.PushCacheList(ctx, roomKey(msg.RoomID), msg)
}

func (s *dbStore) DumpPartialMessages(ctx context.Context) func() {
	return func() {
		keys, err := s.GetCacheKeys(ctx, "room:*")
		if err != nil {
			return
		}

		for _, key := range keys {
			messages := make([]*MessageInfo, 0)
			if err = s.PopCacheList(ctx, key, 15, &messages); err != nil {
				return
			}

			for _, msg := range messages {
				err := s.InsertMessage(ctx, &sqlc.InsertMessageParams{
					RoomID:   msg.RoomID,
					SenderID: msg.SenderID,
					Content:  msg.Content,
					Kind:     msg.Kind,
					SendAt:   msg.SendAt,
				})
				if err != nil {
					return
				}
			}
		}
	}
}

func (s *dbStore) DumpAllMessages(ctx context.Context) {
	keys, err := s.GetCacheKeys(ctx, "room:*")
	if err != nil {
		return
	}

	for _, key := range keys {
		messages := make([]*MessageInfo, 0)
		if err = s.PopCacheList(ctx, key, 0, &messages); err != nil {
			return
		}

		for _, msg := range messages {
			err := s.InsertMessage(ctx, &sqlc.InsertMessageParams{
				RoomID:   msg.RoomID,
				SenderID: msg.SenderID,
				Content:  msg.Content,
				Kind:     msg.Kind,
				SendAt:   msg.SendAt,
			})
			if err != nil {
				return
			}
		}
	}
}
