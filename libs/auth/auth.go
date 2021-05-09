package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v7"
)

type AuthClient interface {
	Store(userID uint64, td *TokenDetails) error
	Fetch(r *http.Request) (uint64, error)
	Delete(r *http.Request) (uint64, error)
}

type authClient struct {
	client *redis.Client
}

func NewAuthClient(client *redis.Client) AuthClient {
	return &authClient{client: client}
}

func (a *authClient) Store(userID uint64, td *TokenDetails) error {
	at := time.Unix(td.AccessExpiration, 0)
	rt := time.Unix(td.RefreshExpiration, 0)
	now := time.Now()

	var err error

	err = a.client.Set(
		td.AccessUuid,
		strconv.Itoa(int(userID)),
		at.Sub(now),
	).Err()
	if err != nil {
		return err
	}

	err = a.client.Set(
		td.RefreshUuid,
		strconv.Itoa(int(userID)),
		rt.Sub(now),
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (a *authClient) Fetch(r *http.Request) (uint64, error) {
	ad, err := extractTokenMetadata(r)
	if err != nil {
		return 0, err
	}

	userID, err := a.client.Get(ad.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	parsedID, _ := strconv.ParseUint(userID, 10, 64)

	return parsedID, nil
}

func (a *authClient) Delete(r *http.Request) (uint64, error) {
	ad, err := extractTokenMetadata(r)
	if err != nil {
		return 0, err
	}

	userID, err := a.client.Del(ad.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	return uint64(userID), nil
}

type accessDetails struct {
	AccessUuid string
	UserId     uint64
}

func extractTokenMetadata(r *http.Request) (*accessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(
			fmt.Sprintf("%.f", claims["user_id"]),
			10,
			64,
		)
		if err != nil {
			return nil, err
		}
		return &accessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}

	return nil, err
}
