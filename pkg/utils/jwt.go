package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	// TokenType is the type of the token.
	TokenType = "Bearer"
	// AccessExpire is the duration of the access token.
	AccessExpire = 1 * time.Hour // 1 hour
	// RefreshExpire is the duration of the refresh token.
	RefreshExpire = 1 * 24 * time.Hour // 1 day
)

var (
	// secretKey token secret key
	secretKey = []byte(viper.GetString("app.secret"))
)

// EncryptToken encrypt token
func EncryptToken(claims *jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Err(err).Msg("Error encrypt token")
		return "", err
	}

	return tokenString, nil
}

// DecryptToken decrypt token
func DecryptToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (
		interface{}, error,
	) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secretKey, nil
	})
	if err != nil {
		log.Err(err).Msg("Error decrypt token")
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("unexpected claims type: %T", token.Claims)
}
