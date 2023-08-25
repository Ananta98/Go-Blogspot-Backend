package token

import (
	"blogspot-project/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var API_SECRET = utils.GetEnv("API_SECRET", "rahasiasekali")

func GetTokenString(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	header := c.Request.Header.Get("Authorization")
	splitted := strings.Split(header, " ")
	if len(splitted) == 2 {
		return splitted[1]
	}
	return ""
}

func GenerateToken(user_id uint) (string, error) {
	token_lifespan, err := strconv.Atoi(utils.GetEnv("TOKEN_HOUR_LIFESPAN", "1"))
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Duration(token_lifespan) * time.Hour).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(API_SECRET))
}

func TokenValid(c *gin.Context) error {
	tokenString := GetTokenString(c)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(API_SECRET), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractTokenID(c *gin.Context) (uint, error) {
	tokenString := GetTokenString(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(API_SECRET), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if token.Valid && ok {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)
		if err != nil {
			return 0, err
		}
		return uint(uid), nil
	}
	return 0, nil
}
