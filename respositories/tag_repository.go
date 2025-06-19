package respositories

import (
	"fmt"
	"strconv"

	"github.com/herdiansc/go-cms/models"
	"gorm.io/gorm"
)

// TagRepository struct
type TagRepository struct {
	db *gorm.DB
}

// NewTagRepository inits TagRepository
func NewTagRepository(db *gorm.DB) TagRepository {
	return TagRepository{db: db}
}

// Create saves a tag data
func (repo TagRepository) Create(data models.Tag) (models.Tag, error) {
	result := repo.db.Create(&data)
	return data, result.Error
}

// List finds list of all tags by filter
func (repo TagRepository) List(params map[string]interface{}) ([]models.TagListItem, error) {
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

	var data []models.TagListItem
	result := repo.db.Debug().
		Model(&models.Tag{}).
		Select("tags.id, tags.uuid, tags.title, count(at.*) as usage_count").
		Joins("left join article_tags at on at.tag_id = tags.id").
		Where(params).
		Group("tags.id, tags.uuid, tags.title").
		Order(fmt.Sprintf("%s %s", orderField, orderDir)).
		Limit(limit).
		Offset(limit * (page - 1)).
		Find(&data)

	if result.Error != nil {
		return []models.TagListItem{}, result.Error
	}

	return data, result.Error
}

// FindByParam finds a tag by a specific param
func (repo TagRepository) FindByParam(param string, value any) (models.TagDetail, error) {
	var data models.Tag
	result := repo.db.Where(fmt.Sprintf("%s = ?", param), value).First(&data)
	if result.Error != nil {
		return models.TagDetail{}, result.Error
	}

	var usageCount int64
	result = repo.db.Debug().
		Model(&models.ArticleTag{}).
		Select("count(*) as usage_count").
		Where("tag_id = ?", data.ID).
		Find(&usageCount)

	return models.TagDetail{
		Tag:        data,
		UsageCount: usageCount,
	}, result.Error
}
