package service

import (
	"errors"
	"fmt"
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
	db        *gorm.DB
}

func NewReviewService(reviewDAO *dao.ReviewDAO, couponDAO *dao.CouponDAO, userDAO *dao.UserDAO, db *gorm.DB) *ReviewService {
	return &ReviewService{
		reviewDAO: reviewDAO,
		couponDAO: couponDAO,
		userDAO:   userDAO,
		db:        db,
	}
}

func (s *ReviewService) SubmitReview(userID uint64, imageURLs []string, content string) (*model.Review, error) {
	if len(imageURLs) == 0 {
		return nil, errors.New("请上传至少一张截图")
	}

	// 检查用户是否已有待审核或已通过的评价
	existing, err := s.reviewDAO.GetActiveByUser(userID)
	if err == nil && existing != nil {
		return nil, errors.New("您已提交过评价，请等待审核结果")
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

// AuditReview 审核点评：通过时原子「创建券 + 更新审核状态」
func (s *ReviewService) AuditReview(reviewID uint64, adminID uint32, pass bool, rejectReason string) error {
	review, err := s.reviewDAO.GetByID(reviewID)
	if err != nil {
		return err
	}
	// 幂等：已审核直接返回错误，避免重复发券
	if review.Status != 0 {
		return fmt.Errorf("review already audited")
	}

	if !pass {
		// 未通过：CAS 确保只有 pending 状态可被置为 rejected
		return s.auditCAS(reviewID, 2, rejectReason, 0, adminID)
	}

	// 通过：事务 = 创建券 + CAS 更新审核状态（若 CAS 失败事务回滚）
	validStart := time.Now()
	validEnd := validStart.AddDate(0, 0, config.C.Business.ReviewCouponValidDays)

	return s.db.Transaction(func(tx *gorm.DB) error {
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
		if err := tx.Create(coupon).Error; err != nil {
			return err
		}
		res := tx.Model(&model.Review{}).
			Where("id = ? AND status = 0", reviewID).
			Updates(map[string]interface{}{
				"status":         1,
				"audit_time":     gorm.Expr("NOW()"),
				"audit_admin_id": adminID,
				"coupon_id":      coupon.ID,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("review already audited")
		}
		return nil
	})
}

// auditCAS 仅在 status=0 时执行审核更新
func (s *ReviewService) auditCAS(reviewID uint64, status int8, rejectReason string, couponID uint64, adminID uint32) error {
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
	res := s.db.Model(&model.Review{}).
		Where("id = ? AND status = 0", reviewID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("review already audited")
	}
	return nil
}
