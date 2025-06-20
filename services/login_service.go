package services

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/herdiansc/go-cms/models"
)

var secretKey = []byte("jwt-secret-key")

// AuthFinder defines auth finder function
type AuthFinder interface {
	FindByUsername(msisdn string) (models.Auth, error)
}

// PasswordComparer defines password comparer function
type PasswordComparer interface {
	VerifyPassword(password, hash string) bool
}

// Signer defines signer func
type Signer func(method jwt.SigningMethod, claim jwt.Claims, opts ...jwt.TokenOption) *jwt.Token

// LoginServices defines login service struct
type LoginServices struct {
	decoder      JsonDecoder
	validator    RequestValidator
	hashComparer PasswordComparer
	signer       Signer
	repo         AuthFinder
}

// NewLoginServices inits LoginServices
func NewLoginServices(jd JsonDecoder, rv RequestValidator, ph PasswordComparer, s Signer, ac AuthFinder) LoginServices {
	return LoginServices{
		decoder:      jd,
		validator:    rv,
		hashComparer: ph,
		signer:       s,
		repo:         ac,
	}
}

// Login performs login service
func (svc LoginServices) Login() (int, models.Response) {
	var data models.LoginRequest
	err := svc.decoder.Decode(&data)
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "Bad Request", Data: err.Error()}
	}
	err = svc.validator.Struct(data)
	if err != nil {
		log.Printf("Failed to validate data: %+v\n", err.Error())
		return http.StatusBadRequest, models.Response{Message: "Bad Request", Data: err.Error()}
	}

	auth, err := svc.repo.FindByUsername(data.Username)
	if err != nil {
		log.Printf("Failed to get data: %+v\n", err.Error())
		return http.StatusNotFound, models.Response{Message: "Login failed", Data: nil}
	}

	isSimilar := svc.hashComparer.VerifyPassword(data.Password, auth.Password)
	if !isSimilar {
		log.Printf("Failed to compare password: %+v\n", isSimilar)
		return http.StatusUnauthorized, models.Response{Message: "Login failed", Data: nil}
	}

	tokenString, err := svc.signer(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       auth.ID,
		"username": auth.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
	}).SignedString(secretKey)
	if err != nil {
		log.Printf("Failed to create jwt: %+v\n", isSimilar)
		return http.StatusInternalServerError, models.Response{Message: "Login failed", Data: err.Error()}
	}

	return http.StatusOK, models.Response{Message: "ok", Data: map[string]string{"token": tokenString}}
}
