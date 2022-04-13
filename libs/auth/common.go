package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type TokenDetails struct {
	AccessToken       string
	AccessUuid        string
	AccessExpiration  int64
	RefreshToken      string
	RefreshUuid       string
	RefreshExpiration int64
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) != 2 {
		return nil, errors.New("invalid authorization request header")
	}

	token, err := jwt.Parse(strArr[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}
