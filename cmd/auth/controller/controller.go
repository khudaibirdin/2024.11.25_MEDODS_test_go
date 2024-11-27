package controller

import (
	"app/cmd/auth/model"
	"app/cmd/auth/service"
	"app/internal/database"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetTokens(c *gin.Context) {
	var GetTokensRequest model.GetTokensRequestStruct
	if err := c.ShouldBindJSON(&GetTokensRequest); err != nil {
		err_message := fmt.Sprintf("Request parameters parsing error: %v", err)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}

	// формируем токены
	access_token, err := service.GenerateAccessToken(GetTokensRequest.GUID, GetTokensRequest.IP)
	if err != nil {
		err_message := fmt.Sprintf("Access token generating error: %s", err)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}
	refresh_token, refresh_token_hashed, err := service.GenerateRefreshToken(GetTokensRequest.GUID)
	if err != nil {
		err_message := fmt.Sprintf("Refresh token generating error: %s", err)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}
	// сохранить refresh токен в бд
	user := model.AuthRefreshTokens{GUID: GetTokensRequest.GUID, RefreshToken: refresh_token_hashed}
	result := database.DB.Create(&user)
	if result.Error != nil {
		err_message := fmt.Sprintf("Database value creation error: %s", result.Error)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}

	// вернуть токены в веб
	c.JSON(http.StatusOK, gin.H{"access_token": access_token, "refresh_token": refresh_token})
}

func RefreshTokens(c *gin.Context) {
	var RefreshTokenRequest model.RefreshTokenRequestStruct
	if err := c.ShouldBindJSON(&RefreshTokenRequest); err != nil {
		err_message := fmt.Sprintf("Request parameters parsing error: %v", err)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}

	// из AccessToken вытаскиваем GUID и IP
	var RequestAccessClaims model.AccessTokenClaims
	err := service.GetAccessTokenClaims(RefreshTokenRequest.AccessToken, &RequestAccessClaims)
	if err != nil {
		err_message := fmt.Sprintf("Token claim error: %s", err)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}
	// если IP другой, формируем сообщение на почту
	if RequestAccessClaims.IP != RefreshTokenRequest.IP {
		err_message := "Ip adresses is not the same in token and request. Email messaged to user."
		service.SendEmail("user@test.com", "Your IP is wrong")
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}

	// ищем в БД RefreshToken по полученному GUID
	var UserAuthRefreshToken model.AuthRefreshTokens
	result := database.DB.First(&UserAuthRefreshToken, "GUID = ?", RequestAccessClaims.GUID)
	if result.Error != nil {
		err_message := fmt.Sprintf("Database value getting error: %s", result.Error)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}

	// сравниваем хэш Refresh из БД и из запроса
	err = bcrypt.CompareHashAndPassword([]byte(UserAuthRefreshToken.RefreshToken), []byte(RefreshTokenRequest.RefreshToken))
	if err != nil {
		err_message := fmt.Sprintf("Hash comparing error: %s", result.Error)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}

	// если совпадает, удаляем старый токен
	result = database.DB.Delete(&UserAuthRefreshToken, "GUID = ?", RequestAccessClaims.GUID)
	if result.Error != nil {
		err_message := fmt.Sprintf("Database value getting error: %s", result.Error)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}

	// выдаем новую пару
	access_token, err := service.GenerateAccessToken(UserAuthRefreshToken.GUID, RefreshTokenRequest.IP)
	if err != nil {
		err_message := fmt.Sprintf("Access token generating error: %s", err)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}
	refresh_token, refresh_token_hashed, err := service.GenerateRefreshToken(UserAuthRefreshToken.GUID)
	if err != nil {
		err_message := fmt.Sprintf("Refresh token generating error: %s", err)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}

	// сохранить refresh токен в бд
	user := model.AuthRefreshTokens{GUID: UserAuthRefreshToken.GUID, RefreshToken: refresh_token_hashed}
	result = database.DB.Create(&user)
	if result.Error != nil {
		err_message := fmt.Sprintf("Database value creation error: %s", result.Error)
		log.Print(err_message)
		c.JSON(http.StatusBadRequest, gin.H{"error": err_message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"AccessToken": access_token, "RefreshToken": refresh_token})
}
