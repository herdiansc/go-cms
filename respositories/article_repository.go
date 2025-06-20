package respositories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/herdiansc/go-cms/models"
	"gorm.io/gorm"
)

// ArticleRepository struct
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository inits ArticleRepository
func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return ArticleRepository{db: db}
}

// List finds list of all articles by filter
func (repo ArticleRepository) List(params map[string]interface{}) ([]models.Article, error) {
	fmt.Println(params)
	var data []models.Article
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

// Create saves an article data
func (repo ArticleRepository) Create(writerID int64, data models.CreateArticleRequest) (models.Article, error) {
	tx := repo.db.Begin()

	article := data.Article()
	article.WriterID = writerID
	if err := tx.Create(&article).Error; err != nil {
		tx.Rollback()
		return models.Article{}, err
	}

	for _, reqTag := range data.Tags {
		var tag models.Tag
		result := tx.Where("lower(title) = ?", reqTag).First(&tag)
		tagID := tag.ID
		if result.Error != nil {
			newTag := models.Tag{
				Title: strings.ToLower(reqTag),
			}
			_ = tx.Create(&newTag)
			tagID = newTag.ID
		}
		_ = tx.Create(&models.ArticleTag{
			ArticleID: article.ID,
			TagID:     tagID,
		})
	}

	tx.Commit()

	return article, nil
}

// FindByParam finds an article by a specific param
func (repo ArticleRepository) FindByParam(param string, value any) (models.Article, error) {
	var data models.Article
	result := repo.db.Where(fmt.Sprintf("%s = ?", param), value).First(&data)
	return data, result.Error
}

// FindByParam finds an article by a specific param
func (repo ArticleRepository) DeleteByParam(param string, value any) error {
	var data models.Article
	result := repo.db.Where(fmt.Sprintf("%s = ?", param), value).First(&data)
	if result.Error != nil {
		return result.Error
	}

	tx := repo.db.Begin()
	result = tx.Delete(&data)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	result = tx.Where("article_id = ?", data.ID).Delete(&models.ArticleTag{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	result = tx.Where("article_id = ?", data.ID).Delete(&models.ArticleHistory{})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()

	return nil
}

// PatchByParam patches an article by a specific param
func (repo ArticleRepository) PatchByParam(id int64, param string, value any) (models.Article, error) {
	var data models.Article
	result := repo.db.Where("id = ?", id).First(&data)
	if result.Error != nil {
		return models.Article{}, result.Error
	}

	switch param {
	case "status":
		data.Status = fmt.Sprintf("%s", value)
	default:
		return models.Article{}, errors.New("empty param")
	}

	result = repo.db.Save(&data)

	return data, result.Error
}
