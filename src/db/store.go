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
	GetUsers(ctx context.Context, arg *sqlc.ListUsersParams) (int64, []*UserInfo, error)
	ModifyUser(ctx context.Context, arg *sqlc.UpdateUserParams) (*UserInfo, error)
	RemoveUser(ctx context.Context, userID int64) error

	CreateSession(ctx context.Context, arg *CreateSessionParams) error
	GetSession(ctx context.Context, sessionID uuid.UUID) (*SessionInfo, error)
	RemoveSession(ctx context.Context, sessionID uuid.UUID) error

	GetUserFriends(ctx context.Context, userID int64) ([]*FriendInfo, error)
	AddFriend(ctx context.Context, userID int64, friendID int64) (*FriendInfo, *FriendInfo, error)
	AcceptFriend(ctx context.Context, userID int64, friendID int64) (*FriendInfo, *RoomInfo, *FriendInfo, *RoomInfo, error)
	RefuseFriend(ctx context.Context, userID int64, friendID int64) error
	RemoveFriend(ctx context.Context, userID int64, friendID int64) (int64, error)

	CacheMessage(ctx context.Context, msg *MessageInfo) error
	GetUserRooms(ctx context.Context, userID int64) ([]*RoomInfo, error)
	GetRoomMessages(ctx context.Context, rooms []*RoomInfo, endTime time.Time) error

	DumpPartialMessages()
	DumpAllMessages(ctx context.Context)
}

var _ Store = (*dbStore)(nil)
