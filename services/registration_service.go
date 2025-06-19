package services

import (
	"log"
	"net/http"

	"github.com/herdiansc/go-cms/models"
)

// JsonDecoder defines json decoder function
type JsonDecoder interface {
	Decode(v any) error
}

// RequestValidator defines request validator function
type RequestValidator interface {
	Struct(s interface{}) error
}

// AuthCreator defines auth creator function
type AuthCreator interface {
	Create(auth models.Auth) error
}

// PasswordHasher defines password hasher function
type PasswordHasher interface {
	HashPassword(password string) (string, error)
}

// RegistrationServices defines registration services struct
type RegistrationServices struct {
	decoder   JsonDecoder
	validator RequestValidator
	hasher    PasswordHasher
	repo      AuthCreator
}

// NewRegistrationServices inits RegistrationServices
func NewRegistrationServices(jd JsonDecoder, rv RequestValidator, ph PasswordHasher, ac AuthCreator) RegistrationServices {
	return RegistrationServices{
		decoder:   jd,
		validator: rv,
		hasher:    ph,
		repo:      ac,
	}
}

// Register performs register service
func (svc RegistrationServices) Register() (int, models.Response) {
	var data models.RegisterRequest
	err := svc.decoder.Decode(&data)
	if err != nil {
		return http.StatusBadRequest, models.Response{Message: "Bad Request", Data: err.Error()}
	}
	err = svc.validator.Struct(data)
	if err != nil {
		log.Printf("Failed to validate data: %+v\n", err.Error())
		return http.StatusBadRequest, models.Response{Message: "Bad Request", Data: err.Error()}
	}
	authData := data.Auth()
	authData.Password, err = svc.hasher.HashPassword(data.Password)
	if err != nil {
		log.Printf("Failed to hash password: %+v\n", err.Error())
		return http.StatusInternalServerError, models.Response{Message: "Failed to save", Data: err.Error()}
	}

	err = svc.repo.Create(authData)
	if err != nil {
		log.Printf("Failed to save: %+v\n", err.Error())
		return http.StatusInternalServerError, models.Response{Message: "Failed to save", Data: err.Error()}
	}
	return http.StatusCreated, models.Response{Message: "ok", Data: nil}
}
