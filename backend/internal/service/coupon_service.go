package service

import (
	"errors"
	"nosocial/internal/dao"
	"nosocial/internal/model"
)

type CouponService struct {
	couponDAO *dao.CouponDAO
}

func NewCouponService(couponDAO *dao.CouponDAO) *CouponService {
	return &CouponService{couponDAO: couponDAO}
}

func (s *CouponService) GetMyCoupons(userID uint64, status int8, page, size int) ([]*model.Coupon, int64, error) {
	offset := (page - 1) * size
	return s.couponDAO.ListByUser(userID, status, offset, size)
}

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
	return s.couponDAO.UseCoupon(couponID, orderID)
}
