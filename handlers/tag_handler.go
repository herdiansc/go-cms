package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/herdiansc/go-cms/models"
	"github.com/herdiansc/go-cms/respositories"
	"github.com/herdiansc/go-cms/services"
	"gorm.io/gorm"
)

// TagHandler struct
type TagHandler struct {
	db *gorm.DB
}

// NewTagHandler inits TagHandler
func NewTagHandler(db *gorm.DB) TagHandler {
	return TagHandler{
		db: db,
	}
}

// Create creates new tag
//
//	@Summary		creates new tag
//	@Description	creates new tag and saves it to the database
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			request			body		models.CreateTagRequest	true	"Request of Creating Tag Object"
//	@Param			Authorization	header		string					true	"With the bearer started"
//	@Success		200				{object}	models.Response			"ok"
//	@Failure		400				{object}	models.Response			"bad request"
//	@Failure		500				{object}	models.Response			"internal server error"
//	@Router			/tags [post]
func (h TagHandler) Create(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	jd := json.NewDecoder(r.Body)
	rv := validator.New(validator.WithRequiredStructEnabled())
	ac := respositories.NewTagRepository(h.db)

	svc := services.NewCreateTagServices(ad, jd, rv, ac)
	code, res := svc.Create()
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// List lists tags
//
//	@Summary		lists tags
//	@Description	lists tags from the database
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"With the bearer started"
//	@Param			page			query		int				false	"page number"		default(1)
//	@Param			limit			query		int				false	"limit per page"	default(10)
//	@Param			orderField		query		string			false	"order field"		default(id)
//	@Param			orderDir		query		string			false	"order dir"			default(desc)
//	@Success		200				{object}	models.Response	"ok"
//	@Failure		400				{object}	models.Response	"bad request"
//	@Failure		500				{object}	models.Response	"internal server error"
//	@Router			/tags [get]
func (h TagHandler) List(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	ac := respositories.NewTagRepository(h.db)

	svc := services.NewListTagServices(ad, ac)
	code, res := svc.List(r.URL.Query())
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// Detail details a tag
//
//	@Summary		details a tag
//	@Description	details a tag from the database
//	@Tags			tag
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"With the bearer started"
//	@Param			uuid			path		string			true	"UUID"
//	@Success		200				{object}	models.Response	"ok"
//	@Failure		400				{object}	models.Response	"bad request"
//	@Failure		500				{object}	models.Response	"internal server error"
//	@Router			/tags/{uuid} [get]
func (h TagHandler) Detail(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	ac := respositories.NewTagRepository(h.db)

	svc := services.NewDetailTagServices(ad, ac)
	code, res := svc.GetDetailByUUID(r.PathValue("uuid"))
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}
