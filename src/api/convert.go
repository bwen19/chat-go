package api

import (
	db "gochat/src/db/sqlc"
)

func convertUser(user *db.User) *User {
	return &User{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Role:     user.Role,
		Deleted:  user.Deleted,
		RoomID:   user.RoomID,
		CreateAt: user.CreateAt,
	}
}

func convertListUsers(users []db.ListUsersRow) *ListUsersResponse {
	if len(users) == 0 {
		return &ListUsersResponse{}
	}

	res := make([]*User, 0, 5)
	for _, user := range users {
		res = append(res, &User{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Role:     user.Role,
			Deleted:  user.Deleted,
			CreateAt: user.CreateAt,
		})
	}

	return &ListUsersResponse{
		Total: users[0].Total,
		Users: res,
	}
}

func convertListSessions(sessions []db.ListSessionsRow) *ListSessionsResponse {
	if len(sessions) == 0 {
		return &ListSessionsResponse{}
	}

	sess := make([]*Session, 0, 5)
	for _, session := range sessions {
		sess = append(sess, &Session{
			ID:        session.ID,
			ClientIp:  session.ClientIp,
			UserAgent: session.UserAgent,
			ExpireAt:  session.ExpireAt,
			CreateAt:  session.CreateAt,
		})
	}

	return &ListSessionsResponse{
		Total:    sessions[0].Total,
		Sessions: sess,
	}
}
