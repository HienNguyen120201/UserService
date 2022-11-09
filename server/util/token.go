package util

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

func CreateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"username":   username,
		"exp":        time.Now().Add(time.Hour * 12).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("SECRET_JWT")))
}

func IsValidToken(token string) (bool, error) {
	logger, err := NewLogger(zap.DebugLevel, true)
	if err != nil {
		log.Println("Fail to initial log:", err)
	}
	defer logger.Core().Sync()

	_, err = jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_JWT")), nil
	})
	if err != nil {
		logger.Sugar().Debug("Fail to parse with claims:", err)
		return false, err
	}
	return true, err
}
