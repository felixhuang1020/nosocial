package service

import (
	"fmt"
	"nosocial/config"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type ShareholderService struct {
	userDAO       *dao.UserDAO
	orderDAO      *dao.ShareholderOrderDAO
	commissionDAO *dao.CommissionRecordDAO
	withdrawalDAO *dao.WithdrawalDAO
	db            *gorm.DB
}

func NewShareholderService(userDAO *dao.UserDAO, orderDAO *dao.ShareholderOrderDAO, commissionDAO *dao.CommissionRecordDAO, withdrawalDAO *dao.WithdrawalDAO, db *gorm.DB) *ShareholderService {
	return &ShareholderService{
		userDAO:       userDAO,
		orderDAO:      orderDAO,
		commissionDAO: commissionDAO,
		withdrawalDAO: withdrawalDAO,
		db:            db,
	}
}

// ApplyShareholder 申请成为股东（创建支付订单）
func (s *ShareholderService) ApplyShareholder(userID uint64) (*model.ShareholderOrder, error) {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user.IsShareholder == 1 {
		return nil, fmt.Errorf("already a shareholder")
	}

	order := &model.ShareholderOrder{
		UserID:  userID,
		OrderNo: utils.GenerateOrderNo("SH"),
		Amount:  config.C.Business.ShareholderFee,
	}
	if err := s.orderDAO.Create(order); err != nil {
		return nil, err
	}
	return order, nil
}

// GetPayParams 获取Mock支付参数
func (s *ShareholderService) GetPayParams(userID uint64, orderNo string) (map[string]interface{}, error) {
	order, err := s.orderDAO.GetByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, fmt.Errorf("order not belong to user")
	}

	// MOCK 支付参数
	return map[string]interface{}{
		"appId":     config.C.WX.AppID,
		"timeStamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonceStr":  utils.GenerateInviteCode(),
		"package":   fmt.Sprintf("prepay_id=mock_prepay_%s", orderNo),
		"signType":  "RSA",
		"paySign":   "mock_pay_sign",
	}, nil
}

// PayCallback 支付回调（MOCK）
func (s *ShareholderService) PayCallback(orderNo, transactionID string) error {
	order, err := s.orderDAO.GetByOrderNo(orderNo)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		if err := tx.Model(&model.ShareholderOrder{}).Where("order_no = ?", orderNo).Updates(map[string]interface{}{
			"pay_status":     1,
			"transaction_id": transactionID,
			"pay_time":       gorm.Expr("NOW()"),
		}).Error; err != nil {
			return err
		}

		// 更新用户为股东（有效期一年）
		inviteCode := utils.GenerateInviteCode()
		expireAt := time.Now().AddDate(1, 0, 0)
		if err := tx.Model(&model.User{}).Where("id = ?", order.UserID).Updates(map[string]interface{}{
			"is_shareholder":        1,
			"shareholder_level":     1,
			"register_fee_paid":     1,
			"shareholder_expire_at": expireAt,
			"invite_code":           inviteCode,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *ShareholderService) GetShareholderList(offset, limit int) ([]*model.User, int64, error) {
	return s.userDAO.ListShareholders(offset, limit)
}

func (s *ShareholderService) GetProfile(userID uint64) (*model.User, error) {
	return s.userDAO.GetByID(userID)
}

func (s *ShareholderService) GetEarnings(userID uint64, offset, limit int) ([]*model.CommissionRecord, int64, error) {
	return s.commissionDAO.ListByShareholder(userID, offset, limit)
}

func (s *ShareholderService) GetTeam(userID uint64) ([]*model.User, error) {
	return s.userDAO.GetChildren(userID)
}

func (s *ShareholderService) GetTeamCount(userID uint64) (int64, error) {
	return s.userDAO.GetChildrenCount(userID)
}

// CalculateCommission 计算佣金
// Withdraw 申请提现
func (s *ShareholderService) Withdraw(userID uint64, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("invalid amount")
	}

	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return err
	}
	if user.IsShareholder != 1 {
		return fmt.Errorf("not a shareholder")
	}
	if user.Balance < amount {
		return fmt.Errorf("insufficient balance")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 扣除余额
		if err := tx.Model(&model.User{}).Where("id = ?", userID).
			Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}

		// 创建提现记录
		record := &model.Withdrawal{
			UserID: userID,
			Amount: amount,
			Status: 0,
		}
		if err := tx.Create(record).Error; err != nil {
			return err
		}

		return nil
	})
}

// GetWithdrawals 获取提现记录
func (s *ShareholderService) GetWithdrawals(userID uint64, offset, limit int) ([]*model.Withdrawal, int64, error) {
	return s.withdrawalDAO.ListByUser(userID, offset, limit)
}

func (s *ShareholderService) CalculateCommission(order *model.DrinkOrder) error {
	if order.ShareholderID == nil || *order.ShareholderID == 0 {
		return nil
	}

	shareholder, err := s.userDAO.GetByID(*order.ShareholderID)
	if err != nil || shareholder.IsShareholder != 1 {
		return nil
	}

	rate := config.C.Business.CommissionRate
	commission := order.PayAmount * rate

	record := &model.CommissionRecord{
		ShareholderID:    *order.ShareholderID,
		ConsumerID:       order.UserID,
		OrderID:          order.ID,
		OrderType:        1,
		OrderAmount:      order.PayAmount,
		CommissionRate:   rate,
		CommissionAmount: commission,
		Status:           0,
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.User{}).Where("id = ?", *order.ShareholderID).
			Update("total_earning", gorm.Expr("total_earning + ?", commission)).Error; err != nil {
			return err
		}
		return nil
	})
}
