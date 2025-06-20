package services

import (
	"log"
	"net/http"
	"net/url"

	"github.com/herdiansc/go-cms/models"
)

// TagCreator defines tag creator function
type TagCreator interface {
	Create(data models.Tag) (models.Tag, error)
}

// CreateTagServices defines tag create service struct
type CreateTagServices struct {
	authData  any
	decoder   JsonDecoder
	validator RequestValidator
	repo      TagCreator
}

// NewCreateTagServices inits CreateTagServices
func NewCreateTagServices(ad any, jd JsonDecoder, rv RequestValidator, ac TagCreator) CreateTagServices {
	return CreateTagServices{
		authData:  ad,
		decoder:   jd,
		validator: rv,
		repo:      ac,
	}
}

// Create performs action of creating article
func (svc CreateTagServices) Create() (int, models.Response) {
	_, ok := svc.authData.(models.VerifyData)
	if !ok {
		log.Printf("Failed to read authData\n")
		return http.StatusBadRequest, models.Response{Message: "error", Data: nil}
	}

	var data models.CreateTagRequest
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

	tag, err := svc.repo.Create(data.Tag())
	if err != nil {
		log.Printf("Failed to save data: %+v\n", err.Error())
		return http.StatusInternalServerError, models.Response{Message: "Failed to save data", Data: err.Error()}
	}

	return http.StatusOK, models.Response{Message: "ok", Data: tag}
}

// TagLister defines tag lister function
type TagLister interface {
	List(params map[string]interface{}) ([]models.TagListItem, error)
}

// ListTagServices defines list tag service struct
type ListTagServices struct {
	authData any
	repo     TagLister
}

// NewListTagServices inits ListTagServices
func NewListTagServices(ad any, al TagLister) ListTagServices {
	return ListTagServices{
		authData: ad,
		repo:     al,
	}
}

// List performs action of listing tags
func (svc ListTagServices) List(q url.Values) (int, models.Response) {
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

// TagDetailer defines tag detailer function
type TagDetailer interface {
	FindByParam(param string, value any) (models.TagDetail, error)
}

// DetailTagServices defines detail tag service struct
type DetailTagServices struct {
	authData any
	repo     TagDetailer
}

// NewDetailTagServices inits DetailTagServices
func NewDetailTagServices(ad any, al TagDetailer) DetailTagServices {
	return DetailTagServices{
		authData: ad,
		repo:     al,
	}
}

// GetDetailByUUID gets detail of a tag by uuid
func (svc DetailTagServices) GetDetailByUUID(id int64) (int, models.Response) {
	_, ok := svc.authData.(models.VerifyData)
	if !ok {
		log.Printf("Failed to read authData\n")
		return http.StatusBadRequest, models.Response{Message: "error", Data: nil}
	}
	data, err := svc.repo.FindByParam("id", id)
	if err != nil {
		log.Printf("Failed to get data: %+v\n", err.Error())
		return http.StatusNotFound, models.Response{Message: "not found", Data: err.Error()}
	}

	return http.StatusOK, models.Response{Message: "ok", Data: data}
}
