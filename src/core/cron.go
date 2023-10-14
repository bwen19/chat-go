package core

import (
	"context"
	db "gochat/src/db/sqlc"
	"log"

	"github.com/robfig/cron/v3"
)

func (s *State) runCron() {
	c := cron.New()
	c.AddFunc("0 3 * * *", s.saveMessages)

	c.Start()
	log.Print("start Cron job")
}

func (s *State) saveMessages() {
	ctx := context.Background()
	keys, err := s.Cache.Keys(ctx, "room:*")
	if err != nil {
		return
	}

	for _, key := range keys {
		messages := make([]*MessageInfo, 0)
		if err = s.Cache.LPopCount(ctx, key, 15, &messages); err != nil {
			return
		}

		for _, msg := range messages {
			err := s.Store.CreateMessage(ctx, db.CreateMessageParams{
				RoomID:   msg.RoomID,
				SenderID: msg.SenderID,
				Content:  msg.Content,
				Kind:     msg.Kind,
				SendAt:   msg.SendAt,
			})
			if err != nil {
				return
			}
		}
	}
}

func (s *State) saveAllMessages() {
	ctx := context.Background()
	keys, err := s.Cache.Keys(ctx, "room:*")
	if err != nil {
		return
	}

	for _, key := range keys {
		messages := make([]*MessageInfo, 0)
		if err = s.Cache.LPopCount(ctx, key, 0, &messages); err != nil {
			return
		}

		for _, msg := range messages {
			err := s.Store.CreateMessage(ctx, db.CreateMessageParams{
				RoomID:   msg.RoomID,
				SenderID: msg.SenderID,
				Content:  msg.Content,
				Kind:     msg.Kind,
				SendAt:   msg.SendAt,
			})
			if err != nil {
				return
			}
		}
	}
}
