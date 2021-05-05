package authtoken

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func verifyToken(r *http.Request) (*jwt.Token, error) {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) != 2 {
		return nil, errors.New("")
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

func StoreAuth(userID uint64, t *AuthToken, c redis.Client) error {
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

func FetchAuth(r *http.Request, c redis.Client) (uint64, error) {
	ad, err := extractTokenMetadata(r)
	if err != nil {
		return 0, err
	}

	userid, err := c.Get(ad.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	userID, _ := strconv.ParseUint(userid, 10, 64)

	return userID, nil
}

func DeleteAuth(r *http.Request, c redis.Client) (int64, error) {
	ad, err := extractTokenMetadata(r)
	if err != nil {
		return 0, err
	}

	userID, err := c.Del(ad.AccessUuid).Result()
	if err != nil {
		return 0, err
	}

	return userID, nil
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
