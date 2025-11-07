package model

import "time"

type User struct {
	UserId         uint       `gorm:"primaryKey;autoIncrement" json:"user_id"`           // Internal numeric ID
	KeycloakUserId string     `gorm:"type:uuid;unique;not null" json:"keycloak_user_id"` // Keycloak UUID
	RefreshToken   string     `gorm:"type:text" json:"refresh_token"`
	TokenExpiry    *time.Time `gorm:"index" json:"token_expiry"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
