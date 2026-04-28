package model

import (
	"time"
)

// BirthdayGift 生日礼品记录
type BirthdayGift struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint64     `gorm:"not null;uniqueIndex:uk_user_year,priority:1" json:"user_id"`
	Year      int        `gorm:"not null;uniqueIndex:uk_user_year,priority:2" json:"year"`
	GiftType  int8       `gorm:"type:smallint;not null;default:1" json:"gift_type"`
	GiftName  string     `gorm:"type:varchar(64);not null" json:"gift_name"`
	GiftValue *float64   `gorm:"type:decimal(8,2)" json:"gift_value,omitempty"`
	Status    int8       `gorm:"type:smallint;not null;default:0" json:"status"`
	ClaimedAt *time.Time `gorm:"type:timestamptz" json:"claimed_at,omitempty"`
	UsedAt    *time.Time `gorm:"type:timestamptz" json:"used_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (BirthdayGift) TableName() string {
	return "birthday_gifts"
}

// Review 离店点评
type Review struct {
	ID            uint64      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64      `gorm:"not null;index:idx_user_id" json:"user_id"`
	DianpingURL   *string     `gorm:"type:varchar(255)" json:"dianping_url,omitempty"`
	ScreenshotURL string      `gorm:"type:varchar(255);not null" json:"screenshot_url"` // 兼容字段：存首张图片 URL
	ImageURLs     StringArray `gorm:"type:jsonb" json:"image_urls,omitempty"`           // 多图数组（JSONB）
	ReviewContent *string     `gorm:"type:text" json:"review_content,omitempty"`
	Status        int8        `gorm:"type:smallint;not null;default:0;index:idx_status" json:"status"`
	RejectReason  *string     `gorm:"type:varchar(255)" json:"reject_reason,omitempty"`
	CouponID      *uint64     `json:"coupon_id,omitempty"`
	AuditTime     *time.Time  `gorm:"type:timestamptz" json:"audit_time,omitempty"`
	AuditAdminID  *uint32     `json:"audit_admin_id,omitempty"`
	CreatedAt     time.Time   `gorm:"autoCreateTime" json:"created_at"`
}

func (Review) TableName() string {
	return "reviews"
}

// Coupon 优惠券
type Coupon struct {
	ID             uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64     `gorm:"not null;index:idx_user_id" json:"user_id"`
	CouponNo       string     `gorm:"type:varchar(32);not null;uniqueIndex:uk_coupon_no" json:"coupon_no"`
	Type           int8       `gorm:"type:smallint;not null" json:"type"`
	Name           string     `gorm:"type:varchar(64);not null" json:"name"`
	Amount         float64    `gorm:"type:decimal(8,2);not null" json:"amount"`
	MinOrderAmount float64    `gorm:"type:decimal(8,2);not null;default:0.00" json:"min_order_amount"`
	Status         int8       `gorm:"type:smallint;not null;default:0;index:idx_status" json:"status"`
	ValidStart     string     `gorm:"type:date;not null" json:"valid_start"`
	ValidEnd       string     `gorm:"type:date;not null" json:"valid_end"`
	UsedAt         *time.Time `gorm:"type:timestamptz" json:"used_at,omitempty"`
	UsedOrderID    *uint64    `json:"used_order_id,omitempty"`
	SourceID       *uint64    `json:"source_id,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (Coupon) TableName() string {
	return "coupons"
}
