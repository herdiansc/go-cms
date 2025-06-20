package services

import (
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/herdiansc/go-cms/models"
)

// TokenVerifyServices defines TokenVerifyServices struct
type TokenVerifyServices struct{}

// NewTokenVerifyServices inits TokenVerifyServices
func NewTokenVerifyServices() TokenVerifyServices {
	return TokenVerifyServices{}
}

// Verify verifies token
func (svc TokenVerifyServices) Verify(authHeader string) (int, models.Response) {
	log.Printf("authHeader: %+v\n", authHeader)
	authHeaders := strings.Split(authHeader, " ")
	if len(authHeaders) < 2 {
		return http.StatusBadRequest, models.Response{Message: "Invalid token", Data: nil}
	}
	token, err := jwt.Parse(authHeaders[1], func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		log.Printf("Failed to parse token: %+v\n", err.Error())
		return http.StatusBadRequest, models.Response{Message: "Invalid token", Data: nil}
	}

	claims := token.Claims.(jwt.MapClaims)
	responseData := models.VerifyData{
		ID:       int64(claims["id"].(float64)),
		Username: claims["username"].(string),
	}

	return http.StatusOK, models.Response{Message: "ok", Data: responseData}
}
