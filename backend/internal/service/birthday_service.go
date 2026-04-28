package service

import (
	"errors"
	"fmt"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/utils"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BirthdayService struct {
	giftDAO *dao.BirthdayGiftDAO
	userDAO *dao.UserDAO
	db      *gorm.DB
}

func NewBirthdayService(giftDAO *dao.BirthdayGiftDAO, userDAO *dao.UserDAO, db *gorm.DB) *BirthdayService {
	return &BirthdayService{
		giftDAO: giftDAO,
		userDAO: userDAO,
		db:      db,
	}
}

func (s *BirthdayService) CheckGift(userID uint64) (*model.BirthdayGift, error) {
	year := time.Now().Year()
	return s.giftDAO.GetByUserAndYear(userID, year)
}

// ClaimGift 领取生日礼品（并发安全 + 幂等）
// 1. 存在记录则 CAS: status=0 -> 1
// 2. 不存在则 INSERT（依赖 (user_id, year) 唯一约束防重复）
func (s *BirthdayService) ClaimGift(userID uint64) error {
	year := time.Now().Year()
	gift, err := s.giftDAO.GetByUserAndYear(userID, year)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if gift != nil && gift.ID > 0 {
		if gift.Status != 0 {
			return fmt.Errorf("birthday gift already claimed")
		}
		res := s.db.Model(&model.BirthdayGift{}).
			Where("id = ? AND status = 0", gift.ID).
			Updates(map[string]interface{}{
				"status":     1,
				"claimed_at": gorm.Expr("NOW()"),
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return fmt.Errorf("birthday gift already claimed")
		}
		return nil
	}
	// 不存在：创建已领取记录；并发冲突由唯一索引兜底
	newGift := &model.BirthdayGift{
		UserID:    userID,
		Year:      year,
		GiftType:  1,
		GiftName:  "生日特调鸡尾酒一杯",
		GiftValue: float64Ptr(88.00),
		Status:    1,
	}
	if err := s.giftDAO.Create(newGift); err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("birthday gift already claimed")
		}
		return err
	}
	return nil
}

// GetHistory 获取生日礼品历史
func (s *BirthdayService) GetHistory(userID uint64) ([]*model.BirthdayGift, error) {
	return s.giftDAO.ListByUser(userID)
}

func (s *BirthdayService) AutoCreateGifts() error {
	today := utils.TodayMonthDay()
	users, err := s.userDAO.GetByBirthday(today)
	if err != nil {
		return err
	}

	year := time.Now().Year()
	for _, user := range users {
		exists, err := s.giftDAO.GetByUserAndYear(user.ID, year)
		if err == nil && exists != nil && exists.ID > 0 {
			continue
		}
		gift := &model.BirthdayGift{
			UserID:    user.ID,
			Year:      year,
			GiftType:  1,
			GiftName:  "生日特调鸡尾酒一杯",
			GiftValue: float64Ptr(88.00),
			Status:    0,
		}
		// 忽略唯一键冲突（并发/重复执行场景）
		_ = s.giftDAO.Create(gift)
	}
	return nil
}

func float64Ptr(f float64) *float64 {
	return &f
}

// isUniqueViolation 判断 Postgres/MySQL 唯一键冲突
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "Duplicate entry") ||
		strings.Contains(msg, "UNIQUE constraint") ||
		strings.Contains(msg, "unique constraint")
}
