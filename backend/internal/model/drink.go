package model

import (
	"time"

	"gorm.io/gorm"
)

// DrinkCategory 酒水分类
type DrinkCategory struct {
	ID        uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(32);not null" json:"name"`
	Icon      *string   `gorm:"type:varchar(255)" json:"icon,omitempty"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	Status    int8      `gorm:"type:smallint;not null;default:1" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (DrinkCategory) TableName() string {
	return "drink_categories"
}

// Drink 酒水表
type Drink struct {
	ID            uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"type:varchar(64);not null" json:"name"`
	EnglishName   *string        `gorm:"type:varchar(128)" json:"english_name,omitempty"`
	CategoryID    uint32         `gorm:"not null;index:idx_category_id" json:"category_id"`
	Price         float64        `gorm:"type:decimal(8,2);not null" json:"price"`
	CostPrice     *float64       `gorm:"type:decimal(8,2)" json:"cost_price,omitempty"`
	Alcohol       *float64       `gorm:"type:decimal(4,1)" json:"alcohol,omitempty"`
	Volume        *string        `gorm:"type:varchar(16)" json:"volume,omitempty"`
	Description   *string        `gorm:"type:text" json:"description,omitempty"`
	Ingredients   *string        `gorm:"type:varchar(255)" json:"ingredients,omitempty"`
	ImageURL      *string        `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	IsRecommended int8           `gorm:"type:smallint;not null;default:0" json:"is_recommended"`
	IsFreeDrink   int8           `gorm:"type:smallint;not null;default:0" json:"is_free_drink"`
	Status        int8           `gorm:"type:smallint;not null;default:1;index:idx_status" json:"status"`
	SortOrder     int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Category      *DrinkCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Drink) TableName() string {
	return "drinks"
}

// DrinkOrder 酒水订单
type DrinkOrder struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo        string     `gorm:"type:varchar(32);not null;uniqueIndex:uk_order_no" json:"order_no"`
	UserID         uint64     `gorm:"not null;index:idx_user_id" json:"user_id"`
	ShareholderID  *uint64    `gorm:"index:idx_shareholder_id" json:"shareholder_id,omitempty"`
	TotalAmount    float64    `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	DiscountAmount float64    `gorm:"type:decimal(10,2);not null;default:0.00" json:"discount_amount"`
	PayAmount      float64    `gorm:"type:decimal(10,2);not null" json:"pay_amount"`
	CouponID       *uint64    `json:"coupon_id,omitempty"`
	Items          *string    `gorm:"type:jsonb" json:"items,omitempty"`
	PrepayID       *string    `gorm:"type:varchar(64)" json:"prepay_id,omitempty"`
	TransactionID  *string    `gorm:"type:varchar(64);uniqueIndex:uk_drink_tx_id" json:"transaction_id,omitempty"`
	NotifyRaw      *string    `gorm:"type:jsonb" json:"notify_raw,omitempty"`
	Status         int8       `gorm:"type:smallint;not null;default:0" json:"status"`
	PayTime        *time.Time `gorm:"type:timestamptz" json:"pay_time,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (DrinkOrder) TableName() string {
	return "drink_orders"
}
