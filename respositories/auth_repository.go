package respositories

import (
	"github.com/herdiansc/go-cms/models"
	"gorm.io/gorm"
)

// AuthRepository struct
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository inits AuthRepository
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return AuthRepository{db: db}
}

// Create saves an auth data
func (repo AuthRepository) Create(auth models.Auth) error {
	result := repo.db.Create(&auth)
	return result.Error
}

// FindByUsername finds an auth by username
func (repo AuthRepository) FindByUsername(username string) (models.Auth, error) {
	var auth models.Auth
	result := repo.db.Where("username = ?", username).First(&auth)
	return auth, result.Error
}
