package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// Payload contains the payload data of the token
type Payload struct {
	ID       uuid.UUID `json:"id"`
	UserID   int64     `json:"user_id"`
	ExpireAt time.Time `json:"expire_at"`
}

// NewPayload creates a new token payload with a specific user
func NewPayload(userID int64, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:       tokenID,
		UserID:   userID,
		ExpireAt: time.Now().Add(duration),
	}
	return payload, nil
}

// IsExpired checks if the token is expired or not
func (payload *Payload) isExpired() error {
	if time.Now().After(payload.ExpireAt) {
		return ErrExpiredToken
	}
	return nil
}
