package db

import (
	"context"
	"gochat/src/db/rdb"
	"gochat/src/db/sqlc"
	"time"

	"github.com/google/uuid"
)

// Store defines all functions to execute queries
type Store interface {
	sqlc.Querier
	rdb.Cacher

	CreateUser(ctx context.Context, arg *CreateUserParams) (*UserInfo, error)
	GetUserByID(ctx context.Context, userID int64) (*UserInfo, error)
	GetUserByName(ctx context.Context, username string) (*sqlc.User, error)
	ListUsers(ctx context.Context, arg *ListUsersParams) (int64, []*UserInfo, error)
	ModifyUser(ctx context.Context, arg *ModifyUserParams) (*UserInfo, error)
	RemoveUser(ctx context.Context, userID int64) error

	CreateSession(ctx context.Context, arg *CreateSessionParams) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (*SessionInfo, error)
	ListSessions(ctx context.Context, arg *ListSessionsParams) (int64, []*SessionInfo, error)
	RemoveSession(ctx context.Context, sessionID uuid.UUID) error

	GetUserRooms(ctx context.Context, userID int64, endTime time.Time) ([]*RoomInfo, error)
	GetUserFriends(ctx context.Context, userID int64) ([]*FriendInfo, error)
	CacheMessage(ctx context.Context, msg *MessageInfo) error
	DumpPartialMessages(ctx context.Context) func()
	DumpAllMessages(ctx context.Context)

	AddFriend(ctx context.Context, userID int64, friendID int64) (*FriendInfo, *FriendInfo, error)
	AcceptFriend(ctx context.Context, userID int64, friendID int64) (*FriendInfo, *RoomInfo, *FriendInfo, *RoomInfo, error)
	RefuseFriend(ctx context.Context, userID int64, friendID int64) error
	RemoveFriend(ctx context.Context, userID int64, friendID int64) (int64, error)

	NewRoom(ctx context.Context, name string, ownerID int64, memberIDs []int64) (*RoomInfo, error)
	GetRoomInfo(ctx context.Context, roomID int64) (*RoomInfo, error)
	RemoveRoom(ctx context.Context, roomID int64) error

	AddMembers(ctx context.Context, roomID int64, memberIDs []int64) ([]*MemberInfo, error)
}

var _ Store = (*dbStore)(nil)
