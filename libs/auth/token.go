package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type Tokenizer interface {
	Create(userID uint64) (*TokenDetails, error)
}

type tokenizer struct{}

var _ Tokenizer = (*tokenizer)(nil)

func NewTokenizer() Tokenizer {
	return &tokenizer{}
}

func (t *tokenizer) Create(userID uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AccessUuid = uuid.NewV4().String()
	td.AccessExpiration = time.Now().Add(time.Minute * 15).Unix()

	td.RefreshUuid = uuid.NewV4().String()
	td.RefreshExpiration = time.Now().Add(time.Hour * 24 * 7).Unix()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["uuid"] = td.AccessUuid
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AccessExpiration

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RefreshExpiration

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func ValidateToken(r *http.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}
