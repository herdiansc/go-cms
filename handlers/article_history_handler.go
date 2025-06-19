package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/herdiansc/go-cms/models"
	"github.com/herdiansc/go-cms/respositories"
	"github.com/herdiansc/go-cms/services"
	"gorm.io/gorm"
)

// ArticleHistoryHandler struct
type ArticleHistoryHandler struct {
	db *gorm.DB
}

// NewArticleHistoryHandler inits ArticleHistoryHandler
func NewArticleHistoryHandler(db *gorm.DB) ArticleHistoryHandler {
	return ArticleHistoryHandler{
		db: db,
	}
}

// Detail details an article history
//
//	@Summary		details an article history
//	@Description	details an article history from the database
//	@Tags			article history
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"With the bearer started"
//	@Param			uuid			path		string			true	"UUID"
//	@Success		200				{object}	models.Response	"ok"
//	@Failure		400				{object}	models.Response	"bad request"
//	@Failure		500				{object}	models.Response	"internal server error"
//	@Router			/article-histories/{uuid} [get]
func (h ArticleHistoryHandler) Detail(w http.ResponseWriter, r *http.Request) {
	ad := r.Context().Value(models.AuthVerifyCtxKey)
	ac := respositories.NewArticleHistoryRepository(h.db)

	svc := services.NewDetailArticleHistoryServices(ad, ac)
	code, res := svc.GetDetailByUUID(r.PathValue("uuid"))
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(res)
}
