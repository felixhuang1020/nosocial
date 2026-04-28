package model

import (
	"time"
)

// PaymentNotifyLog 支付回调原始日志，用于幂等去重与审计
type PaymentNotifyLog struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNo       string    `gorm:"type:varchar(32);not null;index:idx_pnl_order_no" json:"order_no"`
	TransactionID string    `gorm:"type:varchar(64);not null;uniqueIndex:uk_pnl_tx_id" json:"transaction_id"`
	OrderType     string    `gorm:"type:varchar(8);not null;default:''" json:"order_type"` // drink / shareholder
	RawBody       string    `gorm:"type:jsonb" json:"raw_body"`
	Processed     int8      `gorm:"type:smallint;not null;default:0" json:"processed"`
	ReceivedAt    time.Time `gorm:"autoCreateTime" json:"received_at"`
}

func (PaymentNotifyLog) TableName() string {
	return "payment_notify_logs"
}
