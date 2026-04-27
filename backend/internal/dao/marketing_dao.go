package dao

import (
	"nosocial/internal/model"

	"gorm.io/gorm"
)

// BirthdayGiftDAO
type BirthdayGiftDAO struct {
	db *gorm.DB
}

func NewBirthdayGiftDAO(db *gorm.DB) *BirthdayGiftDAO {
	return &BirthdayGiftDAO{db: db}
}

func (d *BirthdayGiftDAO) Create(gift *model.BirthdayGift) error {
	return d.db.Create(gift).Error
}

func (d *BirthdayGiftDAO) GetByUserAndYear(userID uint64, year int) (*model.BirthdayGift, error) {
	var gift model.BirthdayGift
	err := d.db.Where("user_id = ? AND year = ?", userID, year).First(&gift).Error
	return &gift, err
}

func (d *BirthdayGiftDAO) ListByUser(userID uint64) ([]*model.BirthdayGift, error) {
	var gifts []*model.BirthdayGift
	err := d.db.Where("user_id = ?", userID).Order("year DESC").Find(&gifts).Error
	return gifts, err
}

func (d *BirthdayGiftDAO) UpdateStatus(id uint64, status int8) error {
	updates := map[string]interface{}{
		"status":     status,
		"claimed_at": gorm.Expr("NOW()"),
	}
	return d.db.Model(&model.BirthdayGift{}).Where("id = ?", id).Updates(updates).Error
}

// ReviewDAO
type ReviewDAO struct {
	db *gorm.DB
}

func NewReviewDAO(db *gorm.DB) *ReviewDAO {
	return &ReviewDAO{db: db}
}

func (d *ReviewDAO) Create(review *model.Review) error {
	return d.db.Create(review).Error
}

func (d *ReviewDAO) GetByID(id uint64) (*model.Review, error) {
	var review model.Review
	err := d.db.First(&review, id).Error
	return &review, err
}

func (d *ReviewDAO) GetByUser(userID uint64) (*model.Review, error) {
	var review model.Review
	err := d.db.Where("user_id = ?", userID).Order("id DESC").First(&review).Error
	return &review, err
}

func (d *ReviewDAO) List(status int8, offset, limit int) ([]*model.Review, int64, error) {
	var reviews []*model.Review
	var total int64
	query := d.db.Model(&model.Review{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&reviews).Error
	return reviews, total, err
}

func (d *ReviewDAO) Audit(id uint64, status int8, rejectReason string, couponID uint64, adminID uint32) error {
	updates := map[string]interface{}{
		"status":         status,
		"audit_time":     gorm.Expr("NOW()"),
		"audit_admin_id": adminID,
	}
	if couponID > 0 {
		updates["coupon_id"] = couponID
	}
	if rejectReason != "" {
		updates["reject_reason"] = rejectReason
	}
	return d.db.Model(&model.Review{}).Where("id = ?", id).Updates(updates).Error
}

// CouponDAO
type CouponDAO struct {
	db *gorm.DB
}

func NewCouponDAO(db *gorm.DB) *CouponDAO {
	return &CouponDAO{db: db}
}

func (d *CouponDAO) Create(coupon *model.Coupon) error {
	return d.db.Create(coupon).Error
}

func (d *CouponDAO) GetByID(id uint64) (*model.Coupon, error) {
	var coupon model.Coupon
	err := d.db.First(&coupon, id).Error
	return &coupon, err
}

func (d *CouponDAO) ListByUser(userID uint64, status int8, offset, limit int) ([]*model.Coupon, int64, error) {
	var coupons []*model.Coupon
	var total int64
	query := d.db.Where("user_id = ?", userID)
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	query.Model(&model.Coupon{}).Count(&total)
	err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&coupons).Error
	return coupons, total, err
}

func (d *CouponDAO) UseCoupon(id uint64, orderID uint64) error {
	return d.db.Model(&model.Coupon{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        1,
		"used_at":       gorm.Expr("NOW()"),
		"used_order_id": orderID,
	}).Error
}
