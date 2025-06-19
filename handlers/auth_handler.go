package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/herdiansc/go-cms/models"
	"github.com/herdiansc/go-cms/respositories"
	"github.com/herdiansc/go-cms/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler struct
type AuthHandler struct {
	db *gorm.DB
}

// NewAuthHandler inits AuthHandler
func NewAuthHandler(db *gorm.DB) AuthHandler {
	return AuthHandler{
		db: db,
	}
}

// Register saves an auth
//
//	@Summary		Add a new auth to database
//	@Description	Add a new auth to database
//	@Tags			auth
//	@x-order		1
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.RegisterRequest	true	"Request of Creating Order Object"
//	@Success		200		{object}	models.Response			"ok"
//	@Failure		400		{object}	models.Response			"bad request"
//	@Failure		500		{object}	models.Response			"internal server error"
//	@Router			/auth/register [post]
func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	jd := json.NewDecoder(r.Body)
	rv := validator.New(validator.WithRequiredStructEnabled())
	ph := services.NewHashingService(bcrypt.GenerateFromPassword)
	ac := respositories.NewAuthRepository(h.db)

	svc := services.NewRegistrationServices(jd, rv, ph, ac)
	code, res := svc.Register()
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// Login login
//
//	@Summary		Add a new auth to database
//	@Description	Add a new auth to database
//	@x-order		2
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.LoginRequest	true	"Request of Creating Order Object"
//	@Success		200		{object}	models.Response		"ok"
//	@Failure		400		{object}	models.Response		"bad request"
//	@Failure		500		{object}	models.Response		"internal server error"
//	@Router			/auth/login [post]
func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	jd := json.NewDecoder(r.Body)
	rv := validator.New(validator.WithRequiredStructEnabled())
	ph := services.NewHashingCompareService(bcrypt.CompareHashAndPassword)
	js := jwt.NewWithClaims
	ac := respositories.NewAuthRepository(h.db)

	svc := services.NewLoginServices(jd, rv, ph, js, ac)
	code, res := svc.Login()
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// GetProfile   gets profile
//
//	@Summary		Get profile of currently logged in user
//	@Description	Get profile of currently logged in user
//	@x-order		3
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"With the bearer started"
//	@Success		200				{object}	models.Response	"ok"
//	@Failure		400				{object}	models.Response	"bad request"
//	@Failure		500				{object}	models.Response	"internal server error"
//	@Router			/auth/profile [get]
func (h AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	af := respositories.NewAuthRepository(h.db)
	svc := services.NewProfileServices(ad, af)
	code, res := svc.GetProfile()
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// // Verify verifies jwt
// // @Summary		Add a new auth to database
// // @Description	Add a new auth to database
// // @Accept		json
// // @Produce		json
// // @Param Authorization header string true "With the bearer started"
// // @Success		200		{string}	string			"ok"
// // @Router		/auth/verify [post]
// func (h AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
// 	svc := services.NewTokenVerifyServices()
// 	code, res := svc.Verify(r.Header.Get("Authorization"))
// 	w.WriteHeader(code)
// 	json.NewEncoder(w).Encode(res)
// }
