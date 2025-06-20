package services

import (
	"log"
	"net/http"
	"net/url"

	"github.com/herdiansc/go-cms/models"
)

// ArticleHistoryLister defines article history lister function
type ArticleHistoryLister interface {
	List(params map[string]interface{}) ([]models.ArticleHistory, error)
}

// ListArticleHistoryServices defines list article history service struct
type ListArticleHistoryServices struct {
	authData    any
	articleRepo ArticleDetailer
	repo        ArticleHistoryLister
}

// NewListArticleHistoryServices inits ListArticleHistoryServices
func NewListArticleHistoryServices(ad any, ar ArticleDetailer, al ArticleHistoryLister) ListArticleHistoryServices {
	return ListArticleHistoryServices{
		authData:    ad,
		articleRepo: ar,
		repo:        al,
	}
}

// List performs action of listing article histories
func (svc ListArticleHistoryServices) List(articleUuid string, q url.Values) (int, models.Response) {
	_, ok := svc.authData.(models.VerifyData)
	if !ok {
		log.Printf("Failed to read authData\n")
		return http.StatusBadRequest, models.Response{Message: "error", Data: nil}
	}

	article, err := svc.articleRepo.FindByParam("uuid", articleUuid)
	if err != nil {
		log.Printf("Failed to get data: %+v\n", err.Error())
		return http.StatusNotFound, models.Response{Message: "not found", Data: err.Error()}
	}

	params := make(map[string]interface{})
	for k, v := range q {
		params[k] = v[0]
	}
	params["article_id"] = article.ID

	data, _ := svc.repo.List(params)
	if len(data) == 0 {
		log.Printf("Failed to get data")
		return http.StatusNotFound, models.Response{Message: "Not found", Data: nil}
	}

	return http.StatusOK, models.Response{Message: "ok", Data: data}
}

// ArticleHistoryDetailer defines article history detailer function
type ArticleHistoryDetailer interface {
	FindByParam(param string, value any) (models.ArticleHistory, error)
}

// DetailArticleHistoryServices defines detail article history service struct
type DetailArticleHistoryServices struct {
	authData any
	repo     ArticleHistoryDetailer
}

// NewDetailArticleHistoryServices inits DetailArticleHistoryServices
func NewDetailArticleHistoryServices(ad any, al ArticleHistoryDetailer) DetailArticleHistoryServices {
	return DetailArticleHistoryServices{
		authData: ad,
		repo:     al,
	}
}

// GetDetailByUUID gets detail of an article history by uuid
func (svc DetailArticleHistoryServices) GetDetailByUUID(uuid string) (int, models.Response) {
	_, ok := svc.authData.(models.VerifyData)
	if !ok {
		log.Printf("Failed to read authData\n")
		return http.StatusBadRequest, models.Response{Message: "error", Data: nil}
	}

	data, err := svc.repo.FindByParam("uuid", uuid)
	if err != nil {
		log.Printf("Failed to get data: %+v\n", err.Error())
		return http.StatusNotFound, models.Response{Message: "not found", Data: err.Error()}
	}

	return http.StatusOK, models.Response{Message: "ok", Data: data}
}
