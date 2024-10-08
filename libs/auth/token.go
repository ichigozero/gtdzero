package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type TokenDetails struct {
	AccessToken       string
	AccessUUID        string
	AccessExpiration  int64
	RefreshToken      string
	RefreshUUID       string
	RefreshExpiration int64
}

type Tokenizer interface {
	Create(userID uint64) (*TokenDetails, error)
}

type tokenizer struct{}

func NewTokenizer() Tokenizer {
	return &tokenizer{}
}

func (t *tokenizer) Create(userID uint64) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AccessUUID = uuid.NewV4().String()
	td.AccessExpiration = time.Now().Add(time.Minute * 15).Unix()

	td.RefreshUUID = GenerateRefreshUUID(td.AccessUUID)
	td.RefreshExpiration = time.Now().Add(time.Hour * 24 * 7).Unix()

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["uuid"] = td.AccessUUID
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AccessExpiration

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["access_uuid"] = td.AccessUUID
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RefreshExpiration

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func GenerateRefreshUUID(accessUUID string) string {
	return uuid.NewV5(uuid.NameSpaceURL, accessUUID).String()
}

func GetTokenClaims(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, err
}

type AccessTokenDetails struct {
	UUID   string
	UserID uint64
}

func ExtractAccessToken(claims jwt.MapClaims) (*AccessTokenDetails, error) {
	uuid, ok := claims["uuid"].(string)
	if !ok {
		return nil, errors.New("invalid uuid conversion")
	}

	userID, err := strconv.ParseUint(
		fmt.Sprintf("%.f", claims["user_id"]),
		10,
		64,
	)
	if err != nil {
		return nil, err
	}

	return &AccessTokenDetails{
		UUID:   uuid,
		UserID: userID,
	}, nil
}

type RefreshTokenDetails struct {
	AccessUUID  string
	RefreshUUID string
	UserID      uint64
}

func ExtractRefreshToken(claims jwt.MapClaims) (*RefreshTokenDetails, error) {
	accessUUID, ok := claims["access_uuid"].(string)
	if !ok {
		return nil, errors.New("invalid uuid conversion")
	}

	refreshUUID, ok := claims["refresh_uuid"].(string)
	if !ok {
		return nil, errors.New("invalid uuid conversion")
	}

	userID, err := strconv.ParseUint(
		fmt.Sprintf("%.f", claims["user_id"]),
		10,
		64,
	)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenDetails{
		AccessUUID:  accessUUID,
		RefreshUUID: refreshUUID,
		UserID:      userID,
	}, nil
}
