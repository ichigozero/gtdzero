package authtoken

import (
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ichigozero/gtdzero/services/redis"
	"github.com/twinj/uuid"
)

type AuthToken struct {
	AccessToken       string
	AccessUuid        string
	AccessExpiration  int64
	RefreshToken      string
	RefreshUuid       string
	RefreshExpiration int64
}

func Create(userID uint64) (*AuthToken, error) {
	t := &AuthToken{}
	t.AccessUuid = uuid.NewV4().String()
	t.AccessExpiration = time.Now().Add(time.Minute * 15).Unix()

	t.RefreshUuid = uuid.NewV4().String()
	t.RefreshExpiration = time.Now().Add(time.Hour * 24 * 7).Unix()

	var err error
	// TODO use config file instead
	os.Setenv("ACCESS_SECRET", "access-secret")

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["uuid"] = t.AccessUuid
	atClaims["user_id"] = userID
	atClaims["exp"] = t.AccessExpiration

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	t.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// TODO use config file instead
	os.Setenv("REFRESH_SECRET", "refresh-secret")

	rtClaims := jwt.MapClaims{}
	rtClaims["uuid"] = t.RefreshUuid
	rtClaims["user_id"] = userID
	rtClaims["exp"] = t.RefreshExpiration

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	t.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func Store(userID uint64, t *AuthToken, c redis.Client) error {
	at := time.Unix(t.AccessExpiration, 0)
	rt := time.Unix(t.RefreshExpiration, 0)
	now := time.Now()

	var err error

	err = c.Set(t.AccessUuid, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if err != nil {
		return err
	}

	err = c.Set(t.RefreshUuid, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}
