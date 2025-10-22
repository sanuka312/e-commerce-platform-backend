package model

type ProductImage struct {
	ImageID   uint   `gorm:"primaryKey" json:"image_id"`
	ProductID uint   `gorm:"not null" json:"product_id"`
	ImageURL  string `gorm:"size:500; not null" json:"image_url"`
}
