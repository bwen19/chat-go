package db

import (
	"context"
)

func (s *dbStore) GetUserRooms(ctx context.Context, userID int64) ([]*RoomInfo, error) {
	rows, err := s.RetrieveUserRooms(ctx, userID)
	if err != nil {
		return nil, err
	}

	roomMap := make(map[int64]*RoomInfo)
	for _, row := range rows {
		isPrivateRoom := row.Category == "private" && row.MemberID != userID

		member := &MemberInfo{
			ID:     row.MemberID,
			Name:   row.Nickname,
			Avatar: row.Avatar,
			JoinAt: row.JoinAt,
		}

		if room, ok := roomMap[row.RoomID]; ok {
			if isPrivateRoom {
				room.Name = row.Nickname
				room.Cover = row.Avatar
			}
			room.Members = append(room.Members, member)
		} else {
			room := &RoomInfo{
				ID:       row.RoomID,
				Name:     row.Name,
				Cover:    row.Cover,
				Category: row.Category,
				CreateAt: row.CreateAt,
				Members:  []*MemberInfo{member},
				Messages: make([]*MessageInfo, 0),
			}
			if isPrivateRoom {
				room.Name = row.Nickname
				room.Cover = row.Avatar
			}

			roomMap[row.RoomID] = room
		}
	}
	rsp := make([]*RoomInfo, 0, len(roomMap))
	for _, room := range roomMap {
		rsp = append(rsp, room)
	}
	return rsp, nil
}

func (s *dbStore) NewRoom(ctx context.Context, name string, userID int64, memberIDs []int64) (*RoomInfo, error) {
	return nil, nil
}
