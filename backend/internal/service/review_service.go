package service

import (
	"errors"
	"nosocial/config"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type ReviewService struct {
	reviewDAO *dao.ReviewDAO
	couponDAO *dao.CouponDAO
	userDAO   *dao.UserDAO
}

func NewReviewService(reviewDAO *dao.ReviewDAO, couponDAO *dao.CouponDAO, userDAO *dao.UserDAO) *ReviewService {
	return &ReviewService{
		reviewDAO: reviewDAO,
		couponDAO: couponDAO,
		userDAO:   userDAO,
	}
}

func (s *ReviewService) SubmitReview(userID uint64, imageURLs []string, content string) (*model.Review, error) {
	if len(imageURLs) == 0 {
		return nil, errors.New("请上传至少一张截图")
	}
	review := &model.Review{
		UserID:        userID,
		ScreenshotURL: imageURLs[0], // 向后兼容：单图字段写入首张
		ImageURLs:     model.StringArray(imageURLs),
		ReviewContent: &content,
		Status:        0,
	}
	if err := s.reviewDAO.Create(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) GetReviewStatus(userID uint64) (*model.Review, error) {
	review, err := s.reviewDAO.GetByUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) GetReviewList(status int8, page, size int) ([]*model.Review, int64, error) {
	offset := (page - 1) * size
	return s.reviewDAO.List(status, offset, size)
}

func (s *ReviewService) AuditReview(reviewID uint64, adminID uint32, pass bool, rejectReason string) error {
	review, err := s.reviewDAO.GetByID(reviewID)
	if err != nil {
		return err
	}

	if !pass {
		return s.reviewDAO.Audit(reviewID, 2, rejectReason, 0, adminID)
	}

	// 通过：创建优惠券
	validStart := time.Now()
	validEnd := validStart.AddDate(0, 0, config.C.Business.ReviewCouponValidDays)
	coupon := &model.Coupon{
		UserID:         review.UserID,
		CouponNo:       utils.GenerateCouponNo(),
		Type:           1,
		Name:           "点评感谢券",
		Amount:         config.C.Business.ReviewCouponAmount,
		MinOrderAmount: config.C.Business.ReviewCouponMinOrder,
		Status:         0,
		ValidStart:     validStart.Format("2006-01-02"),
		ValidEnd:       validEnd.Format("2006-01-02"),
		SourceID:       &review.ID,
	}

	if err := s.couponDAO.Create(coupon); err != nil {
		return err
	}

	return s.reviewDAO.Audit(reviewID, 1, "", coupon.ID, adminID)
}
