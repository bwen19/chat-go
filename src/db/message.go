package db

import (
	"context"
	"gochat/src/db/sqlc"
	"sort"
	"time"
)

func (s *dbStore) CacheMessage(ctx context.Context, msg *MessageInfo) error {
	return s.PushCacheList(ctx, roomKey(msg.RoomID), msg)
}

func (s *dbStore) GetRoomMessages(ctx context.Context, rooms []*RoomInfo, endTime time.Time) error {
	endTimeMilli := endTime.UnixMilli()
	for _, room := range rooms {
		messages := make([]*MessageInfo, 0)
		if err := s.GetCacheList(ctx, roomKey(room.ID), &messages); err != nil {
			return err
		}

		var offset int64 = 0
		for _, message := range messages {
			curOffset := (endTimeMilli - message.SendAt.UnixMilli()) / 86400000
			if curOffset != offset {
				message.Divide = true
				offset = curOffset
			}
		}
		room.Messages = messages
	}
	sort.Sort(RoomSlice(rooms))
	return nil
}

func (s *dbStore) DumpPartialMessages() {
	ctx := context.Background()
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
