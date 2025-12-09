package model

type Product struct {
	ProductID    uint    `gorm:"primaryKey" json:"product_id"`
	ProductName  string  `gorm:"size:250; not null" json:"product_name"`
	ProductPrice float64 `gorm:"not null" json:"product_price"`
	ProductStock int     `gorm:"not null" json:"product_stock"`
	ProductSlug  string  `gorm:"size:250; unique; not null" json:"product_slug"`
	CategoryID   uint    `gorm:"not null;index" json:"category_id"`
	ImgUrlMain   string  `gorm:"column:image_url_main" json:"image_url_main"`

	// Relationships
	Category      Category       `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category"`
	ProductImages []ProductImage `gorm:"foreignKey:ProductID" json:"product_images"`
}
