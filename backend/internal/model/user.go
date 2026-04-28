package model

import (
	"time"
)

// User 用户表
type User struct {
	ID                  uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Openid              string     `gorm:"type:varchar(64);not null;uniqueIndex:uk_openid" json:"openid"`
	Unionid             *string    `gorm:"type:varchar(64)" json:"unionid,omitempty"`
	Nickname            *string    `gorm:"type:varchar(64)" json:"nickname,omitempty"`
	Avatar              *string    `gorm:"type:varchar(255)" json:"avatar,omitempty"`
	Phone               *string    `gorm:"type:varchar(20)" json:"phone,omitempty"`
	Birthday            *string    `gorm:"type:date" json:"birthday,omitempty"`
	IsShareholder       int8       `gorm:"type:smallint;not null;default:0;index:idx_is_shareholder" json:"is_shareholder"`
	ShareholderLevel    int8       `gorm:"type:smallint;not null;default:0" json:"shareholder_level"`
	ShareholderExpireAt *time.Time `gorm:"type:timestamptz" json:"shareholder_expire_at,omitempty"`
	RegisterFeePaid     int8       `gorm:"type:smallint;not null;default:0" json:"register_fee_paid"`
	FreeDrinkUsed       int8       `gorm:"type:smallint;not null;default:0" json:"free_drink_used"`
	InviteCode          *string    `gorm:"type:varchar(16);uniqueIndex:uk_invite_code" json:"invite_code,omitempty"`
	ParentID            *uint64    `gorm:"index:idx_parent_id" json:"parent_id,omitempty"`
	ParentPath          *string    `gorm:"type:varchar(500)" json:"parent_path,omitempty"`
	Balance             float64    `gorm:"type:decimal(10,2);not null;default:0.00" json:"balance"`
	TotalEarning        float64    `gorm:"type:decimal(10,2);not null;default:0.00" json:"total_earning"`
	Status              int8       `gorm:"type:smallint;not null;default:1" json:"status"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// ShareholderOrder 股东注册订单
type ShareholderOrder struct {
	ID            uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64     `gorm:"not null;index:idx_user_id" json:"user_id"`
	OrderNo       string     `gorm:"type:varchar(32);not null;uniqueIndex:uk_order_no" json:"order_no"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	PayStatus     int8       `gorm:"type:smallint;not null;default:0" json:"pay_status"`
	PayTime       *time.Time `gorm:"type:timestamptz" json:"pay_time,omitempty"`
	TransactionID *string    `gorm:"type:varchar(64)" json:"transaction_id,omitempty"`
	PrepayID      *string    `gorm:"type:varchar(64)" json:"prepay_id,omitempty"`
	NotifyRaw     *string    `gorm:"type:jsonb" json:"notify_raw,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (ShareholderOrder) TableName() string {
	return "shareholder_orders"
}

// CommissionRecord 佣金记录
type CommissionRecord struct {
	ID               uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ShareholderID    uint64     `gorm:"not null;index:idx_shareholder_id" json:"shareholder_id"`
	ConsumerID       uint64     `gorm:"not null;index:idx_consumer_id" json:"consumer_id"`
	OrderID          uint64     `gorm:"not null;uniqueIndex:uk_commission_order,priority:1" json:"order_id"`
	OrderType        int8       `gorm:"type:smallint;not null;uniqueIndex:uk_commission_order,priority:2" json:"order_type"`
	OrderAmount      float64    `gorm:"type:decimal(10,2);not null" json:"order_amount"`
	CommissionRate   float64    `gorm:"type:decimal(4,2);not null" json:"commission_rate"`
	CommissionAmount float64    `gorm:"type:decimal(10,2);not null" json:"commission_amount"`
	Status           int8       `gorm:"type:smallint;not null;default:0;index:idx_status" json:"status"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	SettledAt        *time.Time `gorm:"type:timestamptz" json:"settled_at,omitempty"`
}

func (CommissionRecord) TableName() string {
	return "commission_records"
}

// Withdrawal 提现记录
type Withdrawal struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64     `gorm:"not null;index:idx_user_id" json:"user_id"`
	Amount       float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status       int8       `gorm:"type:smallint;not null;default:0" json:"status"` // 0:待处理 1:已处理 2:已拒绝
	RejectReason *string    `gorm:"type:varchar(255)" json:"reject_reason,omitempty"`
	ProcessedAt  *time.Time `gorm:"type:timestamptz" json:"processed_at,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

func (Withdrawal) TableName() string {
	return "withdrawals"
}
