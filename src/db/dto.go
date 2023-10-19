package db

import (
	"encoding/json"
	"gochat/src/db/sqlc"
	"time"

	"github.com/google/uuid"
)

type UserInfo struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Nickname string    `json:"nickname"`
	Role     string    `json:"role"`
	RoomID   int64     `json:"room_id"`
	Deleted  bool      `json:"deleted"`
	CreateAt time.Time `json:"create_at"`
}

func (u *UserInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (u *UserInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func NewUserInfo(v *sqlc.User) *UserInfo {
	return &UserInfo{
		ID:       v.ID,
		Username: v.Username,
		Nickname: v.Nickname,
		Avatar:   v.Avatar,
		Role:     v.Role,
		Deleted:  v.Deleted,
		RoomID:   v.RoomID,
		CreateAt: v.CreateAt,
	}
}

type SessionInfo struct {
	ID           uuid.UUID `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ClientIp     string    `json:"client_ip"`
	UserAgent    string    `json:"user_agent"`
}

func (s *SessionInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s *SessionInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func NewSessionInfo(v *sqlc.Session) *SessionInfo {
	return &SessionInfo{
		ID:           v.ID,
		UserID:       v.UserID,
		RefreshToken: v.RefreshToken,
		ClientIp:     v.ClientIp,
		UserAgent:    v.UserAgent,
	}
}

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

type MemberInfo struct {
	ID     int64     `json:"id"`
	Name   string    `json:"name"`
	Avatar string    `json:"avatar"`
	Rank   string    `json:"rank"`
	JoinAt time.Time `json:"join_at"`
}

type RoomInfo struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Cover    string         `json:"cover"`
	Category string         `json:"category"`
	Unreads  int64          `json:"unreads"`
	CreateAt time.Time      `json:"create_at"`
	Members  []*MemberInfo  `json:"members"`
	Messages []*MessageInfo `json:"messages"`
}

type RoomSlice []*RoomInfo

func (r RoomSlice) Len() int {
	return len(r)
}

func (r RoomSlice) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RoomSlice) Less(i, j int) bool {
	aSlice := r[i].Messages
	bSlice := r[j].Messages
	if alen := len(aSlice); alen != 0 {
		if blen := len(bSlice); blen != 0 {
			return aSlice[alen-1].SendAt.After(bSlice[blen-1].SendAt)
		} else {
			return true
		}
	} else {
		return false
	}
}
