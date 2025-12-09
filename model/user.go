package model

import "time"

type User struct {
	KeycloakUserID string     `gorm:"primaryKey;index" json:"keycloak_user_id"`
	RefreshToken   string     `gorm:"type:text" json:"refresh_token"`
	TokenExpiry    *time.Time `gorm:"index" json:"token_expiry"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
