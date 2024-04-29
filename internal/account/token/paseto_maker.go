package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func (maker *PasetoMaker) CreateToken(id  uuid.UUID , username string, duration time.Duration) (string, *Payload ,  error) {
	payload, err := NewPayload(id ,username, duration)
	if err != nil {
		return "", payload , nil
	}

	token , err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	return token , payload , err

}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, err
	}
	err = payload.Valid()
	if err != nil {
		return nil, err
	}
	
	return payload, nil
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	// chacha20poly1305.KeySize = 32
	if len(symmetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be as least %d characters", chacha20poly1305.KeySize)
	}
	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}
