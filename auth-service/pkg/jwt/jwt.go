package jwt

import (
	"fmt"

	"github.com/ennemli/todo/todo/configs"
	gojwt "github.com/golang-jwt/jwt/v5"
)

func Create(data interface{}, exp int64) (string, error) {
	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256,
		gojwt.MapClaims{
			"data": data,
			"exp":  exp,
		})
	tokenStr, err := token.SignedString([]byte(configs.GetConfig().Service.JWT_SK))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func Validate(tokenString string) (gojwt.MapClaims, error) {
	token, err := gojwt.Parse(tokenString, func(token *gojwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(configs.GetConfig().Service.JWT_SK), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("token invalid")
	}
	if claims, ok := token.Claims.(gojwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("token invalid")
}
