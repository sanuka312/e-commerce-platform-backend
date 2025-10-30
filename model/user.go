package model

type User struct {
	UserId    uint   `gorm:"size:100;not null" json:"keycloak_user_id"`
	FullName  string `gorm:"size:250;not null" json:"full_name"`
	UserEmail string `gorm:"size:200;not null" json:"user_email"`
}
