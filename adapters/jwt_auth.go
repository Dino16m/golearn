package adapters

import (
	"fmt"
	"time"

	"github.com/dino16m/golearn/config"
	"github.com/dino16m/golearn/errors"
	"github.com/dino16m/golearn/types"
	"github.com/golang-jwt/jwt"
)

type JWTAuthAdapter struct {
	options config.JwtOptions
}

type JWTClaims = map[string]interface{}

func (a JWTAuthAdapter) GetRefreshToken(baseClaims JWTClaims) string {
	baseClaims["iat"] = time.Now().Unix()
	baseClaims["nbf"] = time.Now().Add(a.options.MaxRefresh).Unix()
	baseClaims["exp"] = time.Now().Add(a.options.MaxRefresh).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(baseClaims))

	tokenString, _ := token.SignedString([]byte(a.options.Key))
	return tokenString
}

func (a JWTAuthAdapter) GetToken(claim JWTClaims) string {
	claim["exp"] = time.Now().Add(a.options.Timeout).Unix()
	claim["iat"] = time.Now().Unix()
	claim["nbf"] = time.Now().Add(a.options.Timeout).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claim))

	tokenString, _ := token.SignedString([]byte(a.options.Key))
	return tokenString
}

func (a JWTAuthAdapter) GetTokenPair(claim JWTClaims) (refreshToken string, authToken string) {
	refreshToken = a.GetRefreshToken(jwt.MapClaims{"uid": claim[types.UserIdClaim]})
	authToken = a.GetToken(claim)
	return refreshToken, authToken
}

func (a JWTAuthAdapter) GetClaim(tokenStr string) (map[string]interface{}, errors.ApplicationError) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.InternalServerError(fmt.Sprintf("Unexpected signing method: %v", t.Header["alg"]))
		}
		return []byte(a.options.Key), nil
	})
	if err != nil {
		appError, ok := err.(errors.ApplicationError)
		if ok {
			return nil, appError
		} else {
			return nil, errors.InternalServerError("")
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.UnauthorizedError("Invalid token")
	}

}
