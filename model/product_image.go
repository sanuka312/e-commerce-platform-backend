package model

type ProductImage struct {
	ImageID   uint   `gorm:"primaryKey" json:"image_id"`
	ProductID uint   `gorm:"not null;index" json:"product_id"`
	ImageURL  string `gorm:"size:500; not null" json:"image_url"`

	// Relationship
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
