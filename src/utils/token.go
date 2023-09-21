package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

// -------------------------------------------------------------------
// Payload of paseto token
type Payload struct {
	ID       uuid.UUID `json:"id"`
	UserID   int64     `json:"user_id"`
	Role     string    `json:"role"`
	ExpireAt time.Time `json:"expire_at"`
}

func NewPayload(userID int64, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:       tokenID,
		UserID:   userID,
		Role:     role,
		ExpireAt: time.Now().Add(duration),
	}
	return payload, nil
}

func (payload *Payload) IsExpired() error {
	if time.Now().After(payload.ExpireAt) {
		return ErrExpiredToken
	}
	return nil
}

func (payload *Payload) IsAdmin() bool {
	return payload.Role == "admin"
}

// -------------------------------------------------------------------
// Define of maker interface
type TokenMaker interface {
	CreateToken(userID int64, role string, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

// Maker of paseto token
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (TokenMaker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

func (maker *PasetoMaker) CreateToken(userID int64, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, role, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.IsExpired()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
