package service

import (
	"app/cmd/auth/model"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// функция генерирует jwt access токен
func GenerateAccessToken(guid string, IP string) (string, error) {
	claims := model.AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{},
		GUID:             guid,
		IP:               IP,
	}
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	secret_key := os.Getenv("JWT_SECRET_KEY")
	access_token_str, err := access_token.SignedString([]byte(secret_key))
	if err != nil {
		return "", err
	}
	return access_token_str, nil
}

// функция генерирует refresh токен
func GenerateRefreshToken(guid string) (string, string, error) {
	token := base64.StdEncoding.EncodeToString([]byte(uuid.New().String()))
	hashed_token, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return token, string(hashed_token), nil
}

func GetAccessTokenClaims(RefreshToken string, RequestAccessClaims *model.AccessTokenClaims) error {
	token, err := jwt.Parse(RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("sign method error")
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return err
	}

	// Проверяем, является ли токен валидным и содержит claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		RequestAccessClaims.GUID = claims["guid"].(string)
		RequestAccessClaims.IP = claims["ip"].(string)
		fmt.Println(RequestAccessClaims.GUID)
		fmt.Println(RequestAccessClaims.IP)
		return nil
	}
	return fmt.Errorf("claim error")
}

func SendEmail(address string, message string) {
	log.Printf("Message %s sent to user %s", message, address)
}
