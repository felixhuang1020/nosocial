package service

import (
	"fmt"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/utils"
	"time"
)

type BirthdayService struct {
	giftDAO *dao.BirthdayGiftDAO
	userDAO *dao.UserDAO
}

func NewBirthdayService(giftDAO *dao.BirthdayGiftDAO, userDAO *dao.UserDAO) *BirthdayService {
	return &BirthdayService{
		giftDAO: giftDAO,
		userDAO: userDAO,
	}
}

func (s *BirthdayService) CheckGift(userID uint64) (*model.BirthdayGift, error) {
	year := time.Now().Year()
	return s.giftDAO.GetByUserAndYear(userID, year)
}

func (s *BirthdayService) ClaimGift(userID uint64) error {
	year := time.Now().Year()
	gift, err := s.giftDAO.GetByUserAndYear(userID, year)
	if err != nil {
		// 自动创建礼品并领取
		gift = &model.BirthdayGift{
			UserID:    userID,
			Year:      year,
			GiftType:  1,
			GiftName:  "生日特调鸡尾酒一杯",
			GiftValue: float64Ptr(88.00),
			Status:    1,
		}
		return s.giftDAO.Create(gift)
	}

	if gift.Status != 0 {
		return fmt.Errorf("生日礼品已领取")
	}
	return s.giftDAO.UpdateStatus(gift.ID, 1)
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
		exists, _ := s.giftDAO.GetByUserAndYear(user.ID, year)
		if exists != nil {
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
		s.giftDAO.Create(gift)
	}
	return nil
}

func float64Ptr(f float64) *float64 {
	return &f
}
