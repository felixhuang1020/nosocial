package service

import (
	"fmt"
	"nosocial/config"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/utils"
	"time"
)

type OrderService struct {
	orderDAO *dao.DrinkOrderDAO
}

func NewOrderService(orderDAO *dao.DrinkOrderDAO) *OrderService {
	return &OrderService{orderDAO: orderDAO}
}

type CreateOrderReq struct {
	TotalAmount    float64 `json:"total_amount"`
	DiscountAmount float64 `json:"discount_amount,omitempty"`
	CouponID       *uint64 `json:"coupon_id,omitempty"`
	ShareholderID  *uint64 `json:"shareholder_id,omitempty"`
}

func (s *OrderService) CreateOrder(userID uint64, req *CreateOrderReq) (*model.DrinkOrder, error) {
	// 金额校验
	if req.TotalAmount < 0 {
		return nil, fmt.Errorf("订单总金额不能为负数")
	}
	if req.TotalAmount > 999999.99 {
		return nil, fmt.Errorf("订单总金额超出限制")
	}
	if req.DiscountAmount < 0 {
		return nil, fmt.Errorf("优惠金额不能为负数")
	}
	if req.DiscountAmount > req.TotalAmount {
		return nil, fmt.Errorf("优惠金额不能超过订单总金额")
	}

	payAmount := req.TotalAmount - req.DiscountAmount
	if payAmount < 0 {
		payAmount = 0
	}

	order := &model.DrinkOrder{
		OrderNo:        utils.GenerateOrderNo("DO"),
		UserID:         userID,
		ShareholderID:  req.ShareholderID,
		TotalAmount:    req.TotalAmount,
		DiscountAmount: req.DiscountAmount,
		PayAmount:      payAmount,
		CouponID:       req.CouponID,
		Status:         0,
	}

	if err := s.orderDAO.Create(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) GetUserOrders(userID uint64, page, size int) ([]*model.DrinkOrder, int64, error) {
	offset := (page - 1) * size
	return s.orderDAO.ListByUser(userID, offset, size)
}

func (s *OrderService) GetOrderList(status int8, page, size int) ([]*model.DrinkOrder, int64, error) {
	offset := (page - 1) * size
	return s.orderDAO.ListAll(status, offset, size)
}

func (s *OrderService) UpdateOrderStatus(id uint64, status int8) error {
	return s.orderDAO.UpdateStatus(id, status)
}

// GetOrderByID 根据ID获取订单
func (s *OrderService) GetOrderByID(id uint64) (*model.DrinkOrder, error) {
	return s.orderDAO.GetByID(id)
}

// GetOrderPayParams 获取订单支付参数（MOCK）
func (s *OrderService) GetOrderPayParams(userID uint64, orderID uint64) (map[string]interface{}, error) {
	order, err := s.orderDAO.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, fmt.Errorf("order not belong to user")
	}
	if order.Status != 0 {
		return nil, fmt.Errorf("order cannot be paid")
	}

	// MOCK 支付参数，实际应调用微信支付统一下单接口
	return map[string]interface{}{
		"appId":     config.C.WX.AppID,
		"timeStamp": fmt.Sprintf("%d", time.Now().Unix()),
		"nonceStr":  utils.GenerateInviteCode(),
		"package":   fmt.Sprintf("prepay_id=mock_prepay_%s", order.OrderNo),
		"signType":  "RSA",
		"paySign":   "mock_pay_sign",
	}, nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(userID uint64, orderID uint64) error {
	order, err := s.orderDAO.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.UserID != userID {
		return fmt.Errorf("order not belong to user")
	}
	if order.Status != 0 {
		return fmt.Errorf("only unpaid order can be cancelled")
	}
	return s.orderDAO.UpdateStatus(orderID, 4)
}
