package db

import (
	"context"
	"errors"
	"gochat/src/db/sqlc"
)

func (s *dbStore) GetUserFriends(ctx context.Context, userID int64) ([]*FriendInfo, error) {
	rows, err := s.ListUserFriends(ctx, userID)
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

func (s *dbStore) AddFriend(ctx context.Context, userID int64, friendID int64) (*FriendInfo, *FriendInfo, error) {
	friend, err := s.RetrieveFriend(ctx, &sqlc.RetrieveFriendParams{
		UserID:   userID,
		FriendID: friendID,
	})
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			friend, err = s.createFriend(ctx, userID, friendID)
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	if friend.Status == StatusDeleted {
		friend, err = s.UpdateAddFriend(ctx, &sqlc.UpdateAddFriendParams{
			UserID:   userID,
			FriendID: friendID,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	if friend.Status != StatusAdding {
		return nil, nil, ErrInvalidStatus
	}

	return s.getFriendDetail(ctx, friend)
}

func (s *dbStore) AcceptFriend(ctx context.Context, userID int64, friendID int64) (*FriendInfo, *RoomInfo, *FriendInfo, *RoomInfo, error) {
	friend, err := s.RetrieveFriend(ctx, &sqlc.RetrieveFriendParams{
		UserID:   userID,
		FriendID: friendID,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if friend.FriendID != userID || friend.Status != StatusAdding {
		return nil, nil, nil, nil, ErrInvalidStatus
	}

	err = s.execTx(ctx, func(q *sqlc.Queries) error {
		err := s.UpdateFriend(ctx, &sqlc.UpdateFriendParams{
			Status:   StatusAccepted,
			UserID:   friend.UserID,
			FriendID: friend.FriendID,
		})
		if err != nil {
			return err
		}

		err = s.InsertRoomMember(ctx, &sqlc.InsertRoomMemberParams{
			RoomID:   friend.RoomID,
			MemberID: friend.UserID,
			Rank:     RankMember,
		})
		if err != nil {
			return err
		}

		err = s.InsertRoomMember(ctx, &sqlc.InsertRoomMemberParams{
			RoomID:   friend.RoomID,
			MemberID: friend.FriendID,
			Rank:     RankMember,
		})
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	fFriend, uFriend, err := s.getFriendDetail(ctx, friend)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	fRoom, uRoom, err := s.getFriendRooms(ctx, friend)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return uFriend, uRoom, fFriend, fRoom, nil
}

func (s *dbStore) RefuseFriend(ctx context.Context, userID int64, friendID int64) error {
	friend, err := s.RetrieveFriend(ctx, &sqlc.RetrieveFriendParams{
		UserID:   userID,
		FriendID: friendID,
	})
	if err != nil {
		return err
	}

	if friend.Status != StatusAdding {
		return ErrInvalidStatus
	}

	err = s.UpdateFriend(ctx, &sqlc.UpdateFriendParams{
		Status:   StatusDeleted,
		UserID:   friend.UserID,
		FriendID: friend.FriendID,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *dbStore) RemoveFriend(ctx context.Context, userID int64, friendID int64) (int64, error) {
	friend, err := s.RetrieveFriend(ctx, &sqlc.RetrieveFriendParams{
		UserID:   userID,
		FriendID: friendID,
	})
	if err != nil {
		return 0, err
	}

	if friend.Status != StatusAccepted {
		return 0, ErrInvalidStatus
	}

	err = s.execTx(ctx, func(q *sqlc.Queries) error {
		err := s.UpdateFriend(ctx, &sqlc.UpdateFriendParams{
			Status:   StatusDeleted,
			UserID:   friend.UserID,
			FriendID: friend.FriendID,
		})
		if err != nil {
			return err
		}

		if err = s.DeleteMemberByRoom(ctx, friend.RoomID); err != nil {
			return err
		}

		return err
	})

	return friend.RoomID, err
}

func (s *dbStore) getFriendDetail(ctx context.Context, friend *sqlc.Friendship) (*FriendInfo, *FriendInfo, error) {
	v, err := s.RetrieveFriendDetail(ctx, &sqlc.RetrieveFriendDetailParams{
		UserID:   friend.UserID,
		FriendID: friend.FriendID,
	})
	if err != nil {
		return nil, nil, err
	}

	fFriend := &FriendInfo{
		ID:       v.UserID,
		Username: v.UUsername,
		Nickname: v.UNickname,
		Avatar:   v.UAvatar,
		Status:   v.Status,
		RoomID:   v.RoomID,
		First:    false,
		CreateAt: v.CreateAt,
	}
	uFriend := &FriendInfo{
		ID:       v.FriendID,
		Username: v.FUsername,
		Nickname: v.FNickname,
		Avatar:   v.FAvatar,
		Status:   v.Status,
		RoomID:   v.RoomID,
		First:    true,
		CreateAt: v.CreateAt,
	}
	return uFriend, fFriend, nil
}

func (s *dbStore) getFriendRooms(ctx context.Context, friend *sqlc.Friendship) (*RoomInfo, *RoomInfo, error) {
	vals, err := s.RetrieveFriendRooms(ctx, friend.RoomID)
	if err != nil {
		return nil, nil, err
	}

	if len(vals) != 2 {
		return nil, nil, errors.New("error num of people in friend room")
	}

	var uRoom, fRoom *RoomInfo
	members := make([]*MemberInfo, 2)
	for _, v := range vals {
		m := &MemberInfo{
			ID:     v.MemberID,
			Name:   v.Nickname,
			Avatar: v.Avatar,
			Rank:   v.Rank,
			JoinAt: v.JoinAt,
		}
		members = append(members, m)

		room := &RoomInfo{
			ID:       v.RoomID,
			Name:     v.Nickname,
			Cover:    v.Avatar,
			Category: v.Category,
			CreateAt: v.CreateAt,
			Unreads:  0,
			Members:  nil,
			Messages: make([]*MessageInfo, 0),
		}
		if v.MemberID == friend.UserID {
			fRoom = room
		} else if v.MemberID == friend.FriendID {
			uRoom = room
		} else {
			return nil, nil, errors.New("mismatched members")
		}
	}
	uRoom.Members = members
	fRoom.Members = members

	return uRoom, fRoom, nil
}

func (s *dbStore) createFriend(ctx context.Context, userID int64, friendID int64) (*sqlc.Friendship, error) {
	var ret *sqlc.Friendship

	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		room, err := q.InsertRoom(ctx, &sqlc.InsertRoomParams{
			Name:     privateRoomName,
			Cover:    privateCover,
			Category: CategoryPrivate,
		})
		if err != nil {
			return err
		}

		ret, err = q.InsertFriend(ctx, &sqlc.InsertFriendParams{
			UserID:   userID,
			FriendID: friendID,
			RoomID:   room.ID,
			Status:   StatusAdding,
		})
		if err != nil {
			return err
		}

		return err
	})

	return ret, err
}
