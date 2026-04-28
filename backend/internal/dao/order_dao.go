package dao

import (
	"nosocial/internal/model"

	"gorm.io/gorm"
)

type DrinkOrderDAO struct {
	db *gorm.DB
}

func NewDrinkOrderDAO(db *gorm.DB) *DrinkOrderDAO {
	return &DrinkOrderDAO{db: db}
}

func (d *DrinkOrderDAO) Create(order *model.DrinkOrder) error {
	return d.db.Create(order).Error
}

func (d *DrinkOrderDAO) GetByID(id uint64) (*model.DrinkOrder, error) {
	var order model.DrinkOrder
	err := d.db.First(&order, id).Error
	return &order, err
}

func (d *DrinkOrderDAO) GetByOrderNo(orderNo string) (*model.DrinkOrder, error) {
	var order model.DrinkOrder
	err := d.db.Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

func (d *DrinkOrderDAO) ListByUser(userID uint64, offset, limit int) ([]*model.DrinkOrder, int64, error) {
	var orders []*model.DrinkOrder
	var total int64
	query := d.db.Where("user_id = ?", userID)
	query.Model(&model.DrinkOrder{}).Count(&total)
	err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&orders).Error
	return orders, total, err
}

func (d *DrinkOrderDAO) ListAll(status int8, offset, limit int) ([]*model.DrinkOrder, int64, error) {
	var orders []*model.DrinkOrder
	var total int64
	query := d.db.Model(&model.DrinkOrder{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&orders).Error
	return orders, total, err
}

func (d *DrinkOrderDAO) UpdateStatus(id uint64, status int8) error {
	updates := map[string]interface{}{"status": status}
	if status == 1 { // 已支付
		updates["pay_time"] = gorm.Expr("NOW()")
	}
	return d.db.Model(&model.DrinkOrder{}).Where("id = ?", id).Updates(updates).Error
}

// UpdatePrepayID 写回微信 prepay_id（仅在未支付状态）
func (d *DrinkOrderDAO) UpdatePrepayID(id uint64, prepayID string) error {
	return d.db.Model(&model.DrinkOrder{}).
		Where("id = ? AND status = 0", id).
		Update("prepay_id", prepayID).Error
}

// MarkPaidCAS 幂等置为已支付：仅当 status=0 时执行更新，返回影响行数（=1表示本次真正改状态）
func (d *DrinkOrderDAO) MarkPaidCAS(orderNo string, transactionID string, rawNotify string) (int64, error) {
	res := d.db.Model(&model.DrinkOrder{}).
		Where("order_no = ? AND status = 0", orderNo).
		Updates(map[string]interface{}{
			"status":         1,
			"transaction_id": transactionID,
			"notify_raw":     rawNotify,
			"pay_time":       gorm.Expr("NOW()"),
		})
	return res.RowsAffected, res.Error
}

func (d *DrinkOrderDAO) GetTodayStats() (float64, int64, error) {
	var totalAmount float64
	var count int64
	d.db.Model(&model.DrinkOrder{}).
		Where("DATE(created_at) = CURRENT_DATE AND status >= 1").
		Select("COALESCE(SUM(pay_amount), 0)").
		Scan(&totalAmount)
	d.db.Model(&model.DrinkOrder{}).
		Where("DATE(created_at) = CURRENT_DATE AND status >= 1").
		Count(&count)
	return totalAmount, count, nil
}

// RevenueTrend 近N天营业额趋势
type RevenueTrend struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
}

func (d *DrinkOrderDAO) GetRevenueTrend(days int) ([]RevenueTrend, error) {
	var results []RevenueTrend
	err := d.db.Model(&model.DrinkOrder{}).
		Select("TO_CHAR(created_at, 'MM-DD') as date, COALESCE(SUM(pay_amount), 0) as amount").
		Where("created_at >= CURRENT_DATE - ?::integer AND status >= 1", days).
		Group("TO_CHAR(created_at, 'YYYY-MM-DD'), TO_CHAR(created_at, 'MM-DD')").
		Order("TO_CHAR(created_at, 'YYYY-MM-DD') ASC").
		Scan(&results).Error
	return results, err
}
