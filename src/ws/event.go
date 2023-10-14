package ws

import (
	"gochat/src/core"
	"time"
)

type WebsocketEvent struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type InitializeRequest struct {
	Timestamp time.Time `json:"timestamp"`
}
type InitializeResponse struct {
	Rooms   []*core.RoomInfo   `json:"rooms"`
	Friends []*core.FriendInfo `json:"friends"`
}

type SendMessageRequest struct {
	RoomID  int64  `json:"room_id"`
	Content string `json:"content"`
	Kind    string `json:"kind"`
}
type SendMessageResponse struct {
	Message *core.MessageInfo `json:"message"`
}

type AddFriendRequest struct {
	FriendID int64 `json:"friend_id"`
}
type AddFriendResponse struct {
	Friend *core.FriendInfo `json:"friend"`
}

type AcceptFriendRequest struct {
	FriendID int64 `json:"friend_id"`
}
type AcceptFriendResponse struct {
	Friend *core.FriendInfo `json:"friend"`
	Room   *core.RoomInfo   `json:"room"`
}

type RefuseFriendRequest struct {
	FriendID int64 `json:"friend_id"`
}
type RefuseFriendResponse struct {
	FriendID int64 `json:"friend_id"`
}

type DeleteFriendRequest struct {
	FriendID int64 `json:"friend_id"`
}
type DeleteFriendResponse struct {
	FriendID int64 `json:"friend_id"`
	RoomID   int64 `json:"room_id"`
}

type NewRoomRequest struct {
	Name      string  `json:"name"`
	MemberIDs []int64 `json:"member_ids"`
}
type NewRoomResponse struct {
	Room *core.RoomInfo `json:"room"`
}

type UpdateRoomResquest struct {
	RoomID int64  `json:"room_id"`
	Name   string `json:"name"`
}
type UpdateRoomResponse struct {
	RoomID int64  `json:"room_id"`
	Name   string `json:"name"`
}

type DeleteRoomRequest struct {
	RoomID int64 `json:"room_id"`
}
type DeleteRoomResponse struct {
	RoomID int64 `json:"room_id"`
}

type LeaveRoomRequest struct {
	RoomID int64 `json:"room_id"`
}

type AddMembersRequest struct {
	RoomID    int64   `json:"room_id"`
	MemberIDs []int64 `json:"member_ids"`
}
type AddMembersResponse struct {
	RoomID  int64              `json:"room_id"`
	Members []*core.MemberInfo `json:"members"`
}

type DeleteMembersRequest struct {
	RoomID    int64   `json:"room_id"`
	MemberIDs []int64 `json:"member_ids"`
}
type DeleteMembersResponse struct {
	RoomID    int64   `json:"room_id"`
	MemberIDs []int64 `json:"member_ids"`
}
