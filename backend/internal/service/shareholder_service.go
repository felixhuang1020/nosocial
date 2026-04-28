package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"nosocial/config"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/utils"
	"nosocial/internal/pkg/wxpay"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShareholderService struct {
	userDAO       *dao.UserDAO
	orderDAO      *dao.ShareholderOrderDAO
	commissionDAO *dao.CommissionRecordDAO
	withdrawalDAO *dao.WithdrawalDAO
	wx            *wxpay.Client
	db            *gorm.DB
}

func NewShareholderService(userDAO *dao.UserDAO, orderDAO *dao.ShareholderOrderDAO, commissionDAO *dao.CommissionRecordDAO, withdrawalDAO *dao.WithdrawalDAO, wx *wxpay.Client, db *gorm.DB) *ShareholderService {
	return &ShareholderService{
		userDAO:       userDAO,
		orderDAO:      orderDAO,
		commissionDAO: commissionDAO,
		withdrawalDAO: withdrawalDAO,
		wx:            wx,
		db:            db,
	}
}

// ApplyShareholder 申请成为股东（幂等：复用未支付订单）
func (s *ShareholderService) ApplyShareholder(userID uint64) (*model.ShareholderOrder, error) {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user.IsShareholder == 1 {
		return nil, fmt.Errorf("already a shareholder")
	}

	// 幂等：已有未支付订单则直接复用
	if existing, err := s.orderDAO.GetPendingByUser(userID); err == nil && existing.ID > 0 {
		// 金额漂移保护：若配置改动，复用订单金额保持一致
		return existing, nil
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

// GetPayParams 获取股东订单支付参数（微信 JSAPI 真实下单）
func (s *ShareholderService) GetPayParams(userID uint64, orderNo string) (*wxpay.PrepayParams, error) {
	order, err := s.orderDAO.GetByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, errors.New("order not belong to user")
	}
	if order.PayStatus != 0 {
		return nil, errors.New("order already paid")
	}
	if !s.wx.IsConfigured() {
		return nil, wxpay.ErrNotConfigured
	}

	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user.Openid == "" {
		return nil, errors.New("user openid missing")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 已有 prepay_id 复用
	if order.PrepayID != nil && *order.PrepayID != "" {
		if p, err := s.wx.BuildMiniPaySign(ctx, *order.PrepayID); err == nil {
			return p, nil
		}
	}

	amountFen := int64(math.Round(order.Amount * 100))
	if amountFen <= 0 {
		return nil, errors.New("invalid amount")
	}

	params, err := s.wx.PrepayJSAPI(ctx, wxpay.PrepayJSAPIReq{
		OutTradeNo:  order.OrderNo,
		Description: "NoSocial 股东注册",
		AmountFen:   amountFen,
		OpenID:      user.Openid,
		Attach:      "shareholder",
	})
	if err != nil {
		return nil, err
	}

	_ = s.orderDAO.UpdatePrepayID(order.OrderNo, params.PrepayID)
	return params, nil
}

// PayCallback 股东注册订单支付回调处理（幂等 + 事务）
// 返回值 changed 表示本次调用是否真正将订单置为已支付（用于上层判断是否发通知等）
func (s *ShareholderService) PayCallback(orderNo, transactionID, rawNotify string, paidFen int64) (changed bool, err error) {
	order, err := s.orderDAO.GetByOrderNo(orderNo)
	if err != nil {
		return false, err
	}
	// 金额一致性校验
	expected := int64(math.Round(order.Amount * 100))
	if paidFen > 0 && expected > 0 && paidFen != expected {
		return false, fmt.Errorf("amount mismatch: expect %d, got %d", expected, paidFen)
	}
	// 已支付直接视为幂等成功
	if order.PayStatus == 1 {
		return false, nil
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		// CAS 置为已支付
		res := tx.Model(&model.ShareholderOrder{}).
			Where("order_no = ? AND pay_status = 0", orderNo).
			Updates(map[string]interface{}{
				"pay_status":     1,
				"transaction_id": transactionID,
				"notify_raw":     rawNotify,
				"pay_time":       gorm.Expr("NOW()"),
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			// 已被其它回调处理过
			changed = false
			return nil
		}
		changed = true

		// 升级用户为股东（有效期一年）
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
	return changed, err
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

// Withdraw 申请提现（原子扣余额 + 创建提现记录 + CAS 防超卖）
func (s *ShareholderService) Withdraw(userID uint64, amount float64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	amount = round2(amount)

	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return err
	}
	if user.IsShareholder != 1 {
		return errors.New("not a shareholder")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// CAS 扣余额：WHERE balance >= amount，避免并发提现超额
		res := tx.Model(&model.User{}).
			Where("id = ? AND balance >= ?", userID, amount).
			Update("balance", gorm.Expr("balance - ?", amount))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("insufficient balance")
		}

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

// CalculateCommission 在订单支付成功后调用：原子创建佣金记录 + 增加股东累计收益
func (s *ShareholderService) CalculateCommission(order *model.DrinkOrder) error {
	if order == nil || order.ShareholderID == nil || *order.ShareholderID == 0 {
		return nil
	}
	shareholder, err := s.userDAO.GetByID(*order.ShareholderID)
	if err != nil || shareholder.IsShareholder != 1 {
		return nil
	}

	rate := config.C.Business.CommissionRate
	commission := round2(order.PayAmount * rate)

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
		// 幂等插入：命中唯一索引 uk_commission_order 则跳过，不翻倍结算
		res := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(record)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			// 已存在相同 (order_id, order_type) 的佣金记录，结算核将由原始事务完成
			return nil
		}
		if err := tx.Model(&model.User{}).Where("id = ?", *order.ShareholderID).
			Update("total_earning", gorm.Expr("total_earning + ?", commission)).Error; err != nil {
			return err
		}
		return nil
	})
}
