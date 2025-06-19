package services

import (
	"log"
	"net/http"
	"net/url"

	"github.com/herdiansc/go-cms/models"
)

// ArticleProcessor defines article creator function
type ArticleProcessor interface {
	Create(writerID int64, data models.CreateArticleRequest) (models.Article, error)
	FindByParam(param string, value any) (models.Article, error)
}

// ArticleHistoryCreator defines article history creator function
type ArticleHistoryCreator interface {
	Create(action string, data models.Article) error
}

// CreateArticleServices defines article service struct
type CreateArticleServices struct {
	authData    any
	decoder     JsonDecoder
	validator   RequestValidator
	repo        ArticleProcessor
	historyRepo ArticleHistoryCreator
}

// NewCreateArticleServices inits CreateArticleServices
func NewCreateArticleServices(ad any, jd JsonDecoder, rv RequestValidator, ac ArticleProcessor, hr ArticleHistoryCreator) CreateArticleServices {
	return CreateArticleServices{
		authData:    ad,
		decoder:     jd,
		validator:   rv,
		repo:        ac,
		historyRepo: hr,
	}
}

// Create performs action of creating article
func (svc CreateArticleServices) Create() (int, models.Response) {
	authData, ok := svc.authData.(models.VerifyData)
	if !ok {
		log.Printf("Failed to read authData\n")
		return http.StatusBadRequest, models.Response{Message: "error", Data: nil}
	}

	var data models.CreateArticleRequest
	err := svc.decoder.Decode(&data)
	if err != nil {
		log.Printf("Failed to decode json data: %+v\n", err.Error())
		return http.StatusBadRequest, models.Response{Message: "Bad Request", Data: err.Error()}
	}

	err = svc.validator.Struct(data)
	if err != nil {
		log.Printf("Failed to validate data: %+v\n", err.Error())
		return http.StatusBadRequest, models.Response{Message: "Bad Request", Data: err.Error()}
	}

	article, err := svc.repo.Create(authData.ID, data)
	if err != nil {
		log.Printf("Failed to save data: %+v\n", err.Error())
		return http.StatusInternalServerError, models.Response{Message: "Failed to save data", Data: err.Error()}
	}

	go svc.historyRepo.Create("create", article)

	return http.StatusOK, models.Response{Message: "ok", Data: nil}
}

// ArticleLister defines article lister function
type ArticleLister interface {
	List(params map[string]interface{}) ([]models.Article, error)
}

// ListArticleServices defines list article service struct
type ListArticleServices struct {
	authData any
	repo     ArticleLister
}

// NewListArticleServices inits ListArticleServices
func NewListArticleServices(ad any, al ArticleLister) ListArticleServices {
	return ListArticleServices{
		authData: ad,
		repo:     al,
	}
}

// List performs action of listing articles
func (svc ListArticleServices) List(q url.Values) (int, models.Response) {
	_, ok := svc.authData.(models.VerifyData)
	if !ok {
		log.Printf("Failed to read authData\n")
		return http.StatusBadRequest, models.Response{Message: "error", Data: nil}
	}

	params := make(map[string]interface{})
	for k, v := range q {
		params[k] = v[0]
	}
	data, _ := svc.repo.List(params)
	if len(data) == 0 {
		log.Printf("Failed to get data")
		return http.StatusNotFound, models.Response{Message: "Not found", Data: nil}
	}

	return http.StatusOK, models.Response{Message: "ok", Data: data}
}

// ArticleDetailer defines article detailer function
type ArticleDetailer interface {
	FindByParam(param string, value any) (models.Article, error)
}

// DetailArticleServices defines detail article service struct
type DetailArticleServices struct {
	authData any
	repo     ArticleDetailer
}

// NewDetailArticleServices inits DetailArticleServices
func NewDetailArticleServices(ad any, al ArticleDetailer) DetailArticleServices {
	return DetailArticleServices{
		authData: ad,
		repo:     al,
	}
}

// DetailArticleServices gets detail of an article by uuid
func (svc DetailArticleServices) GetDetailByUUID(uuid string) (int, models.Response) {
	data, err := svc.repo.FindByParam("uuid", uuid)
	if err != nil {
		log.Printf("Failed to get data: %+v\n", err.Error())
		return http.StatusNotFound, models.Response{Message: "not found", Data: err.Error()}
	}

	return http.StatusOK, models.Response{Message: "ok", Data: data}
}

// ArticleDeleter defines article remover function
type ArticleDeleter interface {
	DeleteByParam(param string, value any) error
}

// DeleteArticleServices defines delete article service struct
type DeleteArticleServices struct {
	authData any
	repo     ArticleDeleter
}

// NewDeleteArticleServices inits DeleteArticleServices
func NewDeleteArticleServices(ad any, ade ArticleDeleter) DeleteArticleServices {
	return DeleteArticleServices{
		authData: ad,
		repo:     ade,
	}
}

// Delete gets detail of an article by uuid
func (svc DeleteArticleServices) Delete(uuid string) (int, models.Response) {
	err := svc.repo.DeleteByParam("uuid", uuid)
	if err != nil {
		log.Printf("Failed to delete data: %+v\n", err.Error())
		return http.StatusNotFound, models.Response{Message: "Failed to delete article", Data: err.Error()}
	}

	return http.StatusOK, models.Response{Message: "ok"}
}

// ArticlePatcher defines article patcher function
type ArticlePatcher interface {
	PatchByParam(uuid string, param string, value any) (models.Article, error)
}

// PatchArticleServices defines patch article service struct
type PatchArticleServices struct {
	authData    any
	decoder     JsonDecoder
	validator   RequestValidator
	repo        ArticlePatcher
	historyRepo ArticleHistoryCreator
}

// NewPatchArticleServices inits PatchArticleServices
func NewPatchArticleServices(ad any, jd JsonDecoder, rv RequestValidator, ac ArticlePatcher, hr ArticleHistoryCreator) PatchArticleServices {
	return PatchArticleServices{
		authData:    ad,
		decoder:     jd,
		validator:   rv,
		repo:        ac,
		historyRepo: hr,
	}
}

// Patch performs action of patching an article
func (svc PatchArticleServices) Patch(uuid string) (int, models.Response) {
	_, ok := svc.authData.(models.VerifyData)
	if !ok {
		log.Printf("Failed to read authData\n")
		return http.StatusBadRequest, models.Response{Message: "error", Data: nil}
	}

	var data models.PatchArticleRequest
	err := svc.decoder.Decode(&data)
	if err != nil {
		log.Printf("Failed to decode json data: %+v\n", err.Error())
		return http.StatusBadRequest, models.Response{Message: "Bad Request", Data: err.Error()}
	}

	err = svc.validator.Struct(data)
	if err != nil {
		log.Printf("Failed to validate data: %+v\n", err.Error())
		return http.StatusBadRequest, models.Response{Message: "Bad Request", Data: err.Error()}
	}

	article, err := svc.repo.PatchByParam(uuid, "status", data.Status)
	if err != nil {
		log.Printf("Failed to save data: %+v\n", err.Error())
		return http.StatusInternalServerError, models.Response{Message: "Failed to save data", Data: err.Error()}
	}

	go svc.historyRepo.Create("patch", article)

	return http.StatusOK, models.Response{Message: "ok", Data: nil}
}
