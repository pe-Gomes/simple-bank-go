package token

import (
	"fmt"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
	"golang.org/x/crypto/chacha20poly1305"
)

type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
	implicit     []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	return &PasetoMaker{
		symmetricKey: paseto.NewV4SymmetricKey(),
		implicit:     []byte(symmetricKey),
	}, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	token := paseto.NewToken()

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	token.Set("id", tokenID.String())
	token.Set("username", username)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(duration))
	return token.V4Encrypt(maker.symmetricKey, maker.implicit), nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, token, maker.implicit)

	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, ErrExpiredToken
		}
		return nil, err
	}

	payload, err := getPayloadFromToken(parsedToken)

	if err != nil {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func getPayloadFromToken(token *paseto.Token) (*Payload, error) {
	id, err := token.GetString("id")
	if err != nil {
		return nil, ErrInvalidToken
	}

	username, err := token.GetString("username")
	if err != nil {
		return nil, ErrInvalidToken
	}

	issuedAt, err := token.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}

	expiredAt, err := token.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        uuid.MustParse(id),
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiresAt: expiredAt,
	}, nil
}

var _ Maker = (*PasetoMaker)(nil)
