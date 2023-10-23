package db

import (
	"errors"
	"fmt"
	"gochat/src/util"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
)

var (
	ErrRedisNil       = redis.Nil
	ErrRecordNotFound = pgx.ErrNoRows
	ErrInvalidStatus  = errors.New("invalid friend status")
)

const (
	RoleUser         = "user"
	RoleAdmin        = "admin"
	StatusAdding     = "adding"
	StatusAccepted   = "accepted"
	StatusDeleted    = "deleted"
	CategoryPublic   = "public"
	CategoryPrivate  = "private"
	CategoryPersonal = "personal"
	RankOwner        = "owner"
	RankManager      = "manager"
	RankMember       = "member"
	KindText         = "text"
	KindFile         = "file"
)

const (
	personalRoomName  = "My Room"
	privateRoomName   = "My Friend"
	personalRoomCover = "/cover/personal"
	publicRoomCover   = "/cover/public"
	privateRoomCover  = "/cover/personal"
	defaultAvatar     = "/avatar/default"

	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
)

func HasViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == uniqueViolation || pgErr.Code == foreignKeyViolation {
			return true
		}
	}
	return false
}

func runDatabaseMigration(config *util.Config) error {
	migration, err := migrate.New(config.MigrationUrl, config.DatabaseUrl)
	if err != nil {
		return errors.New("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.New("failed to run migrate up")
	}
	return nil
}

func userKey(userID int64) string {
	return fmt.Sprintf("user:%d", userID)
}

func sessionKey(ID uuid.UUID) string {
	return fmt.Sprintf("session:%s", ID)
}

func roomKey(roomID int64) string {
	return fmt.Sprintf("room:%d", roomID)
}
