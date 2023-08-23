package token

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func GenerateToken(userId string) (string, error) {
	token_lifespan, err := strconv.Atoi(os.Getenv("JWT_AUTH_LIFESPAN_HOURS"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))
}

func parseRawToken(rawToken string) (*jwt.Token, error) {
	return jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
}

func CheckTokenValidity(rawToken string) bool {
	_, err := parseRawToken(rawToken)

	return err == nil
}

func ExtractTokenFromBearerToken(bearerToken string) string {
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}

func ExtractUserIdFromToken(rawToken string) (string, error) {
	token, err := parseRawToken(rawToken)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		userId := claims["userId"].(string)

		return userId, nil
	}

	return "", nil
}

func IsTokenValid(c *gin.Context) bool {
	rawToken := ExtractToken(c)

	return CheckTokenValidity(rawToken)
}

func ExtractToken(c *gin.Context) string {
	return ExtractTokenFromBearerToken(c.Request.Header.Get("Authorization"))
}

func ExtractTokenID(c *gin.Context) (string, error) {
	rawToken := ExtractToken(c)

	return ExtractUserIdFromToken(rawToken)
}
