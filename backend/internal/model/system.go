package model

import (
	"time"

	"gorm.io/gorm"
)

// Banner 轮播图/展示图
type Banner struct {
	ID        uint32         `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     *string        `gorm:"type:varchar(64)" json:"title,omitempty"`
	ImageURL  string         `gorm:"type:varchar(255);not null" json:"image_url"`
	LinkType  int8           `gorm:"type:smallint;not null;default:0" json:"link_type"`
	LinkValue *string        `gorm:"type:varchar(255)" json:"link_value,omitempty"`
	Position  int8           `gorm:"type:smallint;not null;index:idx_position" json:"position"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	Status    int8           `gorm:"type:smallint;not null;default:1" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Banner) TableName() string {
	return "banners"
}

// AdminUser 商家管理员
type AdminUser struct {
	ID          uint32     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string     `gorm:"type:varchar(32);not null;uniqueIndex:uk_username" json:"username"`
	Password    string     `gorm:"type:varchar(255);not null" json:"-"`
	Nickname    *string    `gorm:"type:varchar(32)" json:"nickname,omitempty"`
	Role        int8       `gorm:"type:smallint;not null;default:1" json:"role"`
	Status      int8       `gorm:"type:smallint;not null;default:1" json:"status"`
	LastLoginAt *time.Time `gorm:"type:timestamptz" json:"last_login_at,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (AdminUser) TableName() string {
	return "admin_users"
}

// Setting 系统配置表
type Setting struct {
	ID        uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string    `gorm:"type:varchar(64);not null;uniqueIndex:uk_key" json:"key"`
	Value     string    `gorm:"type:text;not null" json:"value"`
	Desc      *string   `gorm:"type:varchar(255)" json:"desc,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Setting) TableName() string {
	return "settings"
}
