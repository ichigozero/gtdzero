package auth

import (
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
)

type AuthClient interface {
	Store(userID uint64, td *TokenDetails) error
	Fetch(tokenUUID string) (uint64, error)
	Delete(tokenUUID string) error
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
		td.AccessUUID,
		strconv.Itoa(int(userID)),
		at.Sub(now),
	).Err()
	if err != nil {
		return err
	}

	err = a.client.Set(
		td.RefreshUUID,
		strconv.Itoa(int(userID)),
		rt.Sub(now),
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (a *authClient) Fetch(tokenUUID string) (uint64, error) {
	userID, err := a.client.Get(tokenUUID).Result()
	if err != nil {
		return 0, err
	}

	parsedID, _ := strconv.ParseUint(userID, 10, 64)

	return parsedID, nil
}

func (a *authClient) Delete(tokenUUID string) error {
	_, err := a.client.Del(tokenUUID).Result()
	if err != nil {
		return err
	}

	return nil
}
