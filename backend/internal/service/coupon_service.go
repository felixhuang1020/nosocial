package service

import (
	"errors"
	"fmt"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"time"

	"gorm.io/gorm"
)

type CouponService struct {
	couponDAO *dao.CouponDAO
	db        *gorm.DB
}

func NewCouponService(couponDAO *dao.CouponDAO, db *gorm.DB) *CouponService {
	return &CouponService{couponDAO: couponDAO, db: db}
}

func (s *CouponService) GetMyCoupons(userID uint64, status int8, page, size int) ([]*model.Coupon, int64, error) {
	offset := (page - 1) * size
	return s.couponDAO.ListByUser(userID, status, offset, size)
}

// UseCoupon 使用优惠券：CAS 保证并发安全（防止同一张券被两个订单同时核销）
func (s *CouponService) UseCoupon(userID uint64, couponID, orderID uint64) error {
	coupon, err := s.couponDAO.GetByID(couponID)
	if err != nil {
		return errors.New("优惠券不存在")
	}
	if coupon.UserID != userID {
		return errors.New("优惠券不属于当前用户")
	}
	if coupon.Status != 0 {
		return errors.New("优惠券已被使用或过期")
	}
	today := time.Now().Format("2006-01-02")
	if coupon.ValidStart > today || coupon.ValidEnd < today {
		return errors.New("优惠券不在有效期内")
	}

	res := s.db.Model(&model.Coupon{}).
		Where("id = ? AND user_id = ? AND status = 0", couponID, userID).
		Updates(map[string]interface{}{
			"status":        1,
			"used_at":       gorm.Expr("NOW()"),
			"used_order_id": orderID,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("优惠券已被使用")
	}
	return nil
}
