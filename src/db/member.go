package db

import (
	"context"
	"gochat/src/db/sqlc"
)

func (s *dbStore) GetMembersInfo(ctx context.Context, roomID int64) ([]*MemberInfo, error) {
	memberRows, err := s.RetrieveMembers(ctx, roomID)
	if err != nil {
		return nil, err
	}

	members := make([]*MemberInfo, 0, len(memberRows))
	for _, m := range memberRows {
		member := &MemberInfo{
			ID:     m.ID,
			Name:   m.Name,
			Avatar: m.Avatar,
			Rank:   m.Rank,
			JoinAt: m.JoinAt,
		}
		members = append(members, member)
	}

	return members, nil
}

func (s *dbStore) AddMembers(ctx context.Context, roomID int64, memberIDs []int64) ([]*MemberInfo, error) {
	var ret []*MemberInfo

	err := s.execTx(ctx, func(q *sqlc.Queries) error {
		var err error
		for _, memberID := range memberIDs {
			err = q.InsertMember(ctx, &sqlc.InsertMemberParams{
				RoomID:   roomID,
				MemberID: memberID,
				Rank:     RankMember,
			})
			if err != nil {
				return err
			}
		}

		if ret, err = s.GetMembersInfo(ctx, roomID); err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
}
