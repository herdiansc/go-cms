package services

import (
	"log"
	"net/http"

	"github.com/herdiansc/go-cms/models"
)

// ProfileServices defines login service struct
type ProfileServices struct {
	authData any
	repo     AuthFinder
}

// NewProfileServices inits ProfileServices
func NewProfileServices(ad any, af AuthFinder) ProfileServices {
	return ProfileServices{
		authData: ad,
		repo:     af,
	}
}

// GetProfile performs action of getting profile of currently logged in user
func (svc ProfileServices) GetProfile() (int, models.Response) {
	authData, ok := svc.authData.(models.VerifyData)
	if !ok {
		log.Printf("Failed to read authData\n")
		return http.StatusBadRequest, models.Response{Message: "error", Data: nil}
	}
	auth, err := svc.repo.FindByUsername(authData.Username)
	if err != nil {
		log.Printf("Failed to get data: %+v\n", err.Error())
		return http.StatusNotFound, models.Response{Message: "Internal server error", Data: err.Error()}
	}

	return http.StatusOK, models.Response{Message: "ok", Data: auth.ProfileResponse()}
}
