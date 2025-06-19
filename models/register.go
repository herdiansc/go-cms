package models

// RegisterRequest struct
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role"`
}

// Auth creates auth struct from RegisterRequest
func (m RegisterRequest) Auth() Auth {
	roleName := "WRITER"
	if m.Role != "" {
		roleName = m.Role
	}
	return Auth{
		Username: m.Username,
		Password: m.Password,
		RoleName: roleName,
	}
}
