package model

import "time"

type User struct {
	UserId       uint       `gorm:"size:100;not null" json:"user_id"`
	RefreshToken string     `gorm:"type:text"`
	TokenExpiry  *time.Time `gorm:"index"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime"`
}
