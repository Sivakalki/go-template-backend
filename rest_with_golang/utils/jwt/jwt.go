package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ApxJwt struct {
	SecretKey     string
}



func (j *ApxJwt) Decode(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidId
	}

	return claims, nil
}



func (j *ApxJwt) GenerateJwtToken(userId string, duration time.Duration)(string, error){
	claims := jwt.MapClaims{
		"user_id": userId,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func FetchClaim(key string, claims jwt.MapClaims) string {
	if value, exists := claims[key].(string); exists && value != "" {
		return value
	}
	return ""
}
