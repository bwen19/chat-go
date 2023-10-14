package core

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type MessageInfo struct {
	ID       int64     `json:"id"`
	RoomID   int64     `json:"room_id"`
	SenderID int64     `json:"sender_id"`
	Name     string    `json:"name"`
	Avatar   string    `json:"avatar"`
	Content  string    `json:"content"`
	Kind     string    `json:"kind"`
	Divide   bool      `json:"divide"`
	SendAt   time.Time `json:"send_at"`
}

func (m *MessageInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *MessageInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func roomKey(roomID int64) string {
	return fmt.Sprintf("room:%d", roomID)
}

func (s *State) CacheMessage(ctx context.Context, roomID int64, msg *MessageInfo) error {
	return s.Cache.RPush(ctx, roomKey(roomID), msg)
}

func (s *State) GetRoomMessages(ctx context.Context, rooms []*RoomInfo) error {
	for _, room := range rooms {
		msg := make([]*MessageInfo, 0)
		if err := s.Cache.LGetAll(ctx, roomKey(room.ID), &msg); err != nil {
			return err
		}
		room.Messages = msg
	}
	return nil
}
