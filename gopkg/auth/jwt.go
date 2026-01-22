package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(uid string) (string, error) {
	secret := getSecretFromViper()
	duration := viper.GetInt("auth.jwt.expire_hour")
	expire := time.Now().Add(time.Hour * time.Duration(duration))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaims{
		UserID: uid,
		StandardClaims: jwt.StandardClaims{
			Audience:  viper.GetString("auth.jwt.audience"),
			Issuer:    viper.GetString("auth.jwt.issuer"),
			ExpiresAt: expire.Unix(),
		},
	})
	return token.SignedString(secret)
}

func ParseToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, GetJWTSecret)
	if err != nil {
		return nil, err

	}

	if claims, ok := token.Claims.(*JWTClaims); !ok || !token.Valid {
		return nil, errors.New("parse token error")
	} else {
		return claims, nil
	}
}

func getSecretFromViper() []byte {
	return []byte(viper.GetString("auth.jwt.secret"))
}

func GetJWTSecret(token *jwt.Token) (interface{}, error) {
	return getSecretFromViper(), nil
}
