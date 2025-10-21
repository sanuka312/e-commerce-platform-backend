package model

type User struct {
	KeycloakUserID uint   `gorm:"primaryKey" json:"keycloak_user_id"`
	FullName       string `gorm:"size:250; not null" json:"full_name"`
	UserEmail      string `gorm:"size:200;not null" json:"user_email"`
	UserPassword   string `gorm:"size:10;not null json:user_password"`
}
