package dao

import (
	"nosocial/internal/model"

	"gorm.io/gorm"
)

type ShareholderOrderDAO struct {
	db *gorm.DB
}

func NewShareholderOrderDAO(db *gorm.DB) *ShareholderOrderDAO {
	return &ShareholderOrderDAO{db: db}
}

func (d *ShareholderOrderDAO) Create(order *model.ShareholderOrder) error {
	return d.db.Create(order).Error
}

func (d *ShareholderOrderDAO) GetByOrderNo(orderNo string) (*model.ShareholderOrder, error) {
	var order model.ShareholderOrder
	err := d.db.Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

func (d *ShareholderOrderDAO) UpdatePayStatus(orderNo string, status int8, transactionID string) error {
	updates := map[string]interface{}{
		"pay_status":     status,
		"transaction_id": transactionID,
		"pay_time":       gorm.Expr("NOW()"),
	}
	return d.db.Model(&model.ShareholderOrder{}).Where("order_no = ?", orderNo).Updates(updates).Error
}

// GetPendingByUser 获取用户的未支付股东订单（幂等：避免重复创建）
func (d *ShareholderOrderDAO) GetPendingByUser(userID uint64) (*model.ShareholderOrder, error) {
	var order model.ShareholderOrder
	err := d.db.Where("user_id = ? AND pay_status = 0", userID).
		Order("id DESC").First(&order).Error
	return &order, err
}

// UpdatePrepayID 写回 prepay_id（仅在未支付状态）
func (d *ShareholderOrderDAO) UpdatePrepayID(orderNo, prepayID string) error {
	return d.db.Model(&model.ShareholderOrder{}).
		Where("order_no = ? AND pay_status = 0", orderNo).
		Update("prepay_id", prepayID).Error
}

// MarkPaidCAS 幂等置为已支付，返回影响行数
func (d *ShareholderOrderDAO) MarkPaidCAS(orderNo, transactionID, rawNotify string) (int64, error) {
	res := d.db.Model(&model.ShareholderOrder{}).
		Where("order_no = ? AND pay_status = 0", orderNo).
		Updates(map[string]interface{}{
			"pay_status":     1,
			"transaction_id": transactionID,
			"notify_raw":     rawNotify,
			"pay_time":       gorm.Expr("NOW()"),
		})
	return res.RowsAffected, res.Error
}

// CommissionRecordDAO 佣金记录DAO
type CommissionRecordDAO struct {
	db *gorm.DB
}

func NewCommissionRecordDAO(db *gorm.DB) *CommissionRecordDAO {
	return &CommissionRecordDAO{db: db}
}

func (d *CommissionRecordDAO) Create(record *model.CommissionRecord) error {
	return d.db.Create(record).Error
}

func (d *CommissionRecordDAO) ListByShareholder(shareholderID uint64, offset, limit int) ([]*model.CommissionRecord, int64, error) {
	var records []*model.CommissionRecord
	var total int64
	query := d.db.Where("shareholder_id = ?", shareholderID)
	query.Model(&model.CommissionRecord{}).Count(&total)
	err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&records).Error
	return records, total, err
}

func (d *CommissionRecordDAO) UpdateStatus(id uint64, status int8) error {
	updates := map[string]interface{}{
		"status":     status,
		"settled_at": gorm.Expr("NOW()"),
	}
	return d.db.Model(&model.CommissionRecord{}).Where("id = ?", id).Updates(updates).Error
}

// WithdrawalDAO 提现记录DAO
type WithdrawalDAO struct {
	db *gorm.DB
}

func NewWithdrawalDAO(db *gorm.DB) *WithdrawalDAO {
	return &WithdrawalDAO{db: db}
}

func (d *WithdrawalDAO) Create(record *model.Withdrawal) error {
	return d.db.Create(record).Error
}

func (d *WithdrawalDAO) ListByUser(userID uint64, offset, limit int) ([]*model.Withdrawal, int64, error) {
	var records []*model.Withdrawal
	var total int64
	query := d.db.Where("user_id = ?", userID)
	query.Model(&model.Withdrawal{}).Count(&total)
	err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&records).Error
	return records, total, err
}
