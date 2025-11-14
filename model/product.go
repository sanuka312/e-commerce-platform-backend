package model

type Product struct {
	ProductID    uint    `gorm:"primaryKey" json:"product_id"`
	ProductName  string  `gorm:"size:250; not null" json:"product_name"`
	ProductPrice float64 `gorm:"not null" json:"price"`
	ProductStock int     `gorm:"not null" json:"product_stock"`
	ProductSlug  string  `gorm:"size:250; unique; not null" json:"product_slug"`
	CategoryID   uint    `gorm:"not null;index" json:"category_id"`
	CategoryName string  `gorm:"size:100; not null" json:"category_name"`
	ImgUrlMain   string  `json:"image_url_main"`

	// Relationships
	Category      Category       `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category"`
	ProductImages []ProductImage `gorm:"foreignKey:ProductID" json:"product_images"`
}
