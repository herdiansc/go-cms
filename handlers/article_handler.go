package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/herdiansc/go-cms/models"
	"github.com/herdiansc/go-cms/respositories"
	"github.com/herdiansc/go-cms/services"
	"gorm.io/gorm"
)

// ArticleHandler struct
type ArticleHandler struct {
	db *gorm.DB
}

// NewArticleHandler inits ArticleHandler
func NewArticleHandler(db *gorm.DB) ArticleHandler {
	return ArticleHandler{
		db: db,
	}
}

// Create creates new article
//
//	@Summary		creates new article
//	@Description	creates new article and saves it to the database
//	@Tags			article
//	@Accept			json
//	@Produce		json
//	@Param			request			body		models.CreateArticleRequest	true	"Request of Creating Article Object"
//	@Param			Authorization	header		string						true	"Basic [token]. Token obtained from log in endpoint"
//	@Success		200				{object}	models.Response				"ok"
//	@Failure		400				{object}	models.Response				"bad request"
//	@Failure		500				{object}	models.Response				"internal server error"
//	@Router			/articles [post]
func (h ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	jd := json.NewDecoder(r.Body)
	rv := validator.New(validator.WithRequiredStructEnabled())
	ac := respositories.NewArticleRepository(h.db)
	hr := respositories.NewArticleHistoryRepository(h.db)

	svc := services.NewCreateArticleServices(ad, jd, rv, ac, hr)
	code, res := svc.Create()
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// List lists articles
//
//	@Summary		lists articles
//	@Description	lists articles from the database
//	@Tags			article
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Basic [token]. Token obtained from log in endpoint"
//	@Param			page			query		int				false	"page number"		default(1)
//	@Param			limit			query		int				false	"limit per page"	default(10)
//	@Param			orderField		query		string			false	"order field"		default(id)
//	@Param			orderDir		query		string			false	"order dir"			default(desc)
//	@Success		200				{object}	models.Response	"ok"
//	@Failure		400				{object}	models.Response	"bad request"
//	@Failure		500				{object}	models.Response	"internal server error"
//	@Router			/articles [get]
func (h ArticleHandler) List(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	ac := respositories.NewArticleRepository(h.db)

	svc := services.NewListArticleServices(ad, ac)
	code, res := svc.List(r.URL.Query())
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// ListHistories lists articles histories for an article
//
//	@Summary		lists articles histories for an article
//	@Description	lists articles histories for an article from the database
//	@Tags			article
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Basic [token]. Token obtained from log in endpoint"
//	@Param			id				path		integer			true	"ID of article"
//	@Param			page			query		int				false	"page number"		default(1)
//	@Param			limit			query		int				false	"limit per page"	default(10)
//	@Param			orderField		query		string			false	"order field"		default(id)
//	@Param			orderDir		query		string			false	"order dir"			default(desc)
//	@Success		200				{object}	models.Response	"ok"
//	@Failure		400				{object}	models.Response	"bad request"
//	@Failure		500				{object}	models.Response	"internal server error"
//	@Router			/articles/{id}/histories [get]
func (h ArticleHandler) ListHistories(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	ar := respositories.NewArticleRepository(h.db)
	ac := respositories.NewArticleHistoryRepository(h.db)

	svc := services.NewListArticleHistoryServices(ad, ar, ac)
	id, _ := strconv.Atoi(r.PathValue("id"))
	code, res := svc.List(int64(id), r.URL.Query())
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// Detail details an article
//
//	@Summary		details an article
//	@Description	details an article from the database
//	@Tags			article
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Basic [token]. Token obtained from log in endpoint"
//	@Param			id				path		integer			true	"ID of article"
//	@Success		200				{object}	models.Response	"ok"
//	@Failure		400				{object}	models.Response	"bad request"
//	@Failure		500				{object}	models.Response	"internal server error"
//	@Router			/articles/{id} [get]
func (h ArticleHandler) Detail(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	ac := respositories.NewArticleRepository(h.db)

	svc := services.NewDetailArticleServices(ad, ac)
	id, _ := strconv.Atoi(r.PathValue("id"))
	code, res := svc.GetDetailByUUID(int64(id))
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// Delete deletes an article
//
//	@Summary		deletes an article
//	@Description	deletes an article from the database
//	@Tags			article
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Basic [token]. Token obtained from log in endpoint"
//	@Param			id				path		integer			true	"ID of article"
//	@Success		200				{object}	models.Response	"ok"
//	@Failure		400				{object}	models.Response	"bad request"
//	@Failure		500				{object}	models.Response	"internal server error"
//	@Router			/articles/{id} [delete]
func (h ArticleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	ade := respositories.NewArticleRepository(h.db)

	svc := services.NewDeleteArticleServices(ad, ade)
	id, _ := strconv.Atoi(r.PathValue("id"))
	code, res := svc.Delete(int64(id))
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}

// Patch patches an article
//
//	@Summary		patches an article
//	@Description	patches an article from the database, example to update article status
//	@Tags			article
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Basic [token]. Token obtained from log in endpoint"
//	@Param			request			body		models.PatchArticleRequest	true	"Request of Creating Article Object"
//	@Param			id				path		integer						true	"ID of article"
//	@Success		200				{object}	models.Response				"ok"
//	@Failure		400				{object}	models.Response				"bad request"
//	@Failure		500				{object}	models.Response				"internal server error"
//	@Router			/articles/{id} [patch]
func (h ArticleHandler) Patch(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	jd := json.NewDecoder(r.Body)
	rv := validator.New(validator.WithRequiredStructEnabled())
	ade := respositories.NewArticleRepository(h.db)
	hr := respositories.NewArticleHistoryRepository(h.db)

	svc := services.NewPatchArticleServices(ad, jd, rv, ade, hr)
	id, _ := strconv.Atoi(r.PathValue("id"))
	code, res := svc.Patch(int64(id))
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}
