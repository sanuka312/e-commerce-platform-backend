package model

type Category struct {
	CategoryID   uint   `gorm:"primaryKey" json:"category_id"`
	CategoryName string `gorm:"size:100; not null" json:"category_name"`
	CategorySlug string `gorm:"size:100; not null" json:"category_slug"`

	Products []Product `gorm:"foreignKey:CategoryID" json:"products"`
}
