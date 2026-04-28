package dao

import (
	"nosocial/internal/model"

	"gorm.io/gorm"
)

// PaymentNotifyLogDAO 微信支付回调流水 DAO
type PaymentNotifyLogDAO struct {
	db *gorm.DB
}

func NewPaymentNotifyLogDAO(db *gorm.DB) *PaymentNotifyLogDAO {
	return &PaymentNotifyLogDAO{db: db}
}

// Insert 插入回调记录；若 transaction_id 已存在则返回 ErrDuplicated
func (d *PaymentNotifyLogDAO) Insert(log *model.PaymentNotifyLog) error {
	return d.db.Create(log).Error
}

// ExistsTxID 查询同 transaction_id 是否已入库（回调去重）
func (d *PaymentNotifyLogDAO) ExistsTxID(txID string) (bool, error) {
	var cnt int64
	if err := d.db.Model(&model.PaymentNotifyLog{}).
		Where("transaction_id = ?", txID).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// MarkProcessed 标记回调已处理
func (d *PaymentNotifyLogDAO) MarkProcessed(txID string) error {
	return d.db.Model(&model.PaymentNotifyLog{}).
		Where("transaction_id = ?", txID).
		Update("processed", 1).Error
}
