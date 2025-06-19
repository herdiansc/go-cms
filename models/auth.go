package models

// Auth struct
type Auth struct {
	Base
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	RoleName string `gorm:"not null"`
}

// ProfileResponse struct
type ProfileResponse struct {
	PublicBase
	Username string `gorm:"not null;unique"`
	RoleName string `gorm:"not null"`
}

func (a Auth) ProfileResponse() ProfileResponse {
	return ProfileResponse{
		PublicBase: PublicBase{
			UUID:      a.UUID,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
			DeletedAt: a.DeletedAt,
		},
		Username: a.Username,
		RoleName: a.RoleName,
	}
}
