package db

import (
	"context"
	"gochat/src/db/sqlc"
)

func (s *dbStore) NewRoom(ctx context.Context, name string, ownerID int64, memberIDs []int64) (*RoomInfo, error) {
	var ret *RoomInfo

	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		room, err := q.InsertRoom(ctx, &sqlc.InsertRoomParams{
			Name:     name,
			Cover:    publicRoomCover,
			Category: CategoryPublic,
		})
		if err != nil {
			return err
		}

		err = q.InsertMember(ctx, &sqlc.InsertMemberParams{
			RoomID:   room.ID,
			MemberID: ownerID,
			Rank:     RankOwner,
		})
		if err != nil {
			return err
		}

		for _, memberID := range memberIDs {
			err = q.InsertMember(ctx, &sqlc.InsertMemberParams{
				RoomID:   room.ID,
				MemberID: memberID,
				Rank:     RankMember,
			})
			if err != nil {
				return err
			}
		}

		members, err := s.GetMembersInfo(ctx, room.ID)
		if err != nil {
			return err
		}

		ret = &RoomInfo{
			ID:       room.ID,
			Name:     room.Name,
			Cover:    room.Cover,
			Category: room.Category,
			CreateAt: room.CreateAt,
			Members:  members,
			Messages: make([]*MessageInfo, 0),
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (s *dbStore) GetRoomInfo(ctx context.Context, roomID int64) (*RoomInfo, error) {
	room, err := s.RetrieveRoom(ctx, roomID)
	if err != nil {
		return nil, err
	}

	members, err := s.GetMembersInfo(ctx, room.ID)
	if err != nil {
		return nil, err
	}

	ret := &RoomInfo{
		ID:       room.ID,
		Name:     room.Name,
		Cover:    room.Cover,
		Category: room.Category,
		CreateAt: room.CreateAt,
		Members:  members,
		Messages: make([]*MessageInfo, 0),
	}

	return ret, nil
}

func (s *dbStore) RemoveRoom(ctx context.Context, roomID int64) error {
	return s.execTx(ctx, func(q *sqlc.Queries) error {
		if err := q.DeleteMessagesByRoom(ctx, roomID); err != nil {
			return err
		}

		if err := q.DeleteMembersByRoom(ctx, roomID); err != nil {
			return err
		}

		if err := q.DeleteRoom(ctx, roomID); err != nil {
			return err
		}

		return nil
	})
}
