package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorExpiredToken = errors.New("token has expried")
	ErrInvalidToken   = errors.New("token is invalid")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(id uuid.UUID, UserID uuid.UUID, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		ID:        id,
		UserID:    UserID,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}
