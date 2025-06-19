package respositories

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/herdiansc/go-cms/models"
	"gorm.io/gorm"
)

// ArticleHistoryRepository struct
type ArticleHistoryRepository struct {
	db *gorm.DB
}

// NewArticleHistoryRepository inits ArticleHistoryRepository
func NewArticleHistoryRepository(db *gorm.DB) ArticleHistoryRepository {
	return ArticleHistoryRepository{db: db}
}

// Create saves an article data
func (repo ArticleHistoryRepository) Create(action string, data models.Article) error {
	return repo.db.Transaction(func(tx *gorm.DB) error {
		lastArticleHistory := models.ArticleHistory{}
		result := repo.db.Where("article_id = ?", data.ID).Order("id desc").First(&lastArticleHistory)
		var version int64 = 1
		if result.Error == nil {
			version = lastArticleHistory.Version + 1
		}
		articleJson, _ := json.Marshal(data)
		newArticleHistory := models.ArticleHistory{
			Article:   string(articleJson),
			Version:   version,
			Status:    data.Status,
			ArticleID: data.ID,
			Action:    action,
		}

		if err := tx.Create(&newArticleHistory).Error; err != nil {
			return err
		}

		return nil
	})
}

// List lists of all article histories by filter
func (repo ArticleHistoryRepository) List(params map[string]interface{}) ([]models.ArticleHistory, error) {
	var data []models.ArticleHistory
	limit := 10
	if _, ok := params["limit"]; ok {
		limitStr, _ := params["limit"].(string)
		limit, _ = strconv.Atoi(limitStr)
		delete(params, "limit")
	}
	page := 1
	if _, ok := params["page"]; ok {
		pageStr, _ := params["page"].(string)
		page, _ = strconv.Atoi(pageStr)
		delete(params, "page")
	}
	orderField := "id"
	if _, ok := params["orderField"]; ok {
		orderField, _ = params["orderField"].(string)
		delete(params, "orderField")
	}
	orderDir := "desc"
	if _, ok := params["orderDir"]; ok {
		orderDir, _ = params["orderDir"].(string)
		delete(params, "orderDir")
	}

	result := repo.db.Debug().Where(params).
		Order(fmt.Sprintf("%s %s", orderField, orderDir)).
		Limit(limit).
		Offset(limit * (page - 1)).
		Find(&data)
	return data, result.Error
}

// FindByParam finds an article by a specific param
func (repo ArticleHistoryRepository) FindByParam(param string, value any) (models.ArticleHistory, error) {
	var data models.ArticleHistory
	result := repo.db.Where(fmt.Sprintf("%s = ?", param), value).First(&data)
	return data, result.Error
}
