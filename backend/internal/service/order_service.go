package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/utils"
	"nosocial/internal/pkg/wxpay"
	"time"
)

type OrderService struct {
	orderDAO  *dao.DrinkOrderDAO
	drinkDAO  *dao.DrinkDAO
	couponDAO *dao.CouponDAO
	userDAO   *dao.UserDAO
	wx        *wxpay.Client
}

func NewOrderService(orderDAO *dao.DrinkOrderDAO, drinkDAO *dao.DrinkDAO, couponDAO *dao.CouponDAO, userDAO *dao.UserDAO, wx *wxpay.Client) *OrderService {
	return &OrderService{
		orderDAO:  orderDAO,
		drinkDAO:  drinkDAO,
		couponDAO: couponDAO,
		userDAO:   userDAO,
		wx:        wx,
	}
}

// CreateOrderItem 下单明细（服务端按 drink_id 重算价格）
type CreateOrderItem struct {
	DrinkID  uint64 `json:"drink_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,min=1,max=99"`
}

// CreateOrderReq 创建订单请求
// 说明：TotalAmount / DiscountAmount 均由服务端重算，前端金额入参仅作展示校验，不会采纳
type CreateOrderReq struct {
	Items         []CreateOrderItem `json:"items" binding:"required,min=1,max=30"`
	CouponID      *uint64           `json:"coupon_id,omitempty"`
	ShareholderID *uint64           `json:"shareholder_id,omitempty"`
}

// itemSnapshot 保存到 drink_orders.items 的明细快照
type itemSnapshot struct {
	DrinkID  uint64  `json:"drink_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Subtotal float64 `json:"subtotal"`
}

// CreateOrder 创建酒水订单（金额服务端重算 + 券校验）
func (s *OrderService) CreateOrder(userID uint64, req *CreateOrderReq) (*model.DrinkOrder, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("订单明细不能为空")
	}

	// 聚合 drink_id（合并同款）
	idQtyMap := make(map[uint64]int, len(req.Items))
	ids := make([]uint64, 0, len(req.Items))
	for _, it := range req.Items {
		if it.DrinkID == 0 || it.Quantity <= 0 {
			return nil, errors.New("下单明细参数无效")
		}
		if _, ok := idQtyMap[it.DrinkID]; !ok {
			ids = append(ids, it.DrinkID)
		}
		idQtyMap[it.DrinkID] += it.Quantity
	}

	drinks, err := s.drinkDAO.ListByIDs(ids)
	if err != nil {
		return nil, fmt.Errorf("load drinks: %w", err)
	}
	if len(drinks) != len(ids) {
		return nil, errors.New("部分酒水不存在或已下架")
	}

	// 服务端重算金额
	var total float64
	snapshots := make([]itemSnapshot, 0, len(drinks))
	for _, d := range drinks {
		qty := idQtyMap[d.ID]
		subtotal := round2(d.Price * float64(qty))
		total = round2(total + subtotal)
		snapshots = append(snapshots, itemSnapshot{
			DrinkID:  d.ID,
			Name:     d.Name,
			Price:    d.Price,
			Quantity: qty,
			Subtotal: subtotal,
		})
	}
	if total <= 0 {
		return nil, errors.New("订单金额无效")
	}
	if total > 999999.99 {
		return nil, errors.New("订单总金额超出限制")
	}

	// 券校验与折扣重算
	var discount float64
	if req.CouponID != nil && *req.CouponID > 0 {
		c, err := s.couponDAO.GetByID(*req.CouponID)
		if err != nil {
			return nil, errors.New("优惠券不存在")
		}
		if c.UserID != userID {
			return nil, errors.New("优惠券不属于当前用户")
		}
		if c.Status != 0 {
			return nil, errors.New("优惠券已使用或已作废")
		}
		// 有效期校验（date 字符串 YYYY-MM-DD）
		today := time.Now().Format("2006-01-02")
		if c.ValidStart > today || c.ValidEnd < today {
			return nil, errors.New("优惠券不在有效期内")
		}
		if total < c.MinOrderAmount {
			return nil, fmt.Errorf("订单满 %.2f 元可用", c.MinOrderAmount)
		}
		discount = round2(c.Amount)
		if discount > total {
			discount = total
		}
	}

	payAmount := round2(total - discount)
	if payAmount < 0 {
		payAmount = 0
	}

	// 推荐人校验（可选）
	if req.ShareholderID != nil && *req.ShareholderID > 0 {
		sh, err := s.userDAO.GetByID(*req.ShareholderID)
		if err != nil || sh.IsShareholder != 1 || *req.ShareholderID == userID {
			// 无效推荐人静默忽略，避免影响下单
			req.ShareholderID = nil
		}
	}

	itemsJSON, err := json.Marshal(snapshots)
	if err != nil {
		return nil, fmt.Errorf("marshal items: %w", err)
	}
	itemsStr := string(itemsJSON)

	order := &model.DrinkOrder{
		OrderNo:        utils.GenerateOrderNo("DO"),
		UserID:         userID,
		ShareholderID:  req.ShareholderID,
		TotalAmount:    total,
		DiscountAmount: discount,
		PayAmount:      payAmount,
		CouponID:       req.CouponID,
		Items:          &itemsStr,
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

// GetOrderPayParams 获取订单支付参数（调用微信 JSAPI 下单）
func (s *OrderService) GetOrderPayParams(userID uint64, orderID uint64) (*wxpay.PrepayParams, error) {
	order, err := s.orderDAO.GetByID(orderID)
	if err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, errors.New("order not belong to user")
	}
	if order.Status != 0 {
		return nil, errors.New("order cannot be paid")
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

	// 已有 prepay_id 复用（避免重复下单）
	if order.PrepayID != nil && *order.PrepayID != "" {
		if p, err := s.wx.BuildMiniPaySign(ctx, *order.PrepayID); err == nil {
			return p, nil
		}
		// 复用失败则走重新下单
	}

	amountFen := int64(math.Round(order.PayAmount * 100))
	if amountFen <= 0 {
		return nil, errors.New("invalid pay amount")
	}

	params, err := s.wx.PrepayJSAPI(ctx, wxpay.PrepayJSAPIReq{
		OutTradeNo:  order.OrderNo,
		Description: "NoSocial 酒水订单",
		AmountFen:   amountFen,
		OpenID:      user.Openid,
		Attach:      "drink",
	})
	if err != nil {
		return nil, err
	}

	// 写回 prepay_id（CAS 保护，只有未支付才更新）
	_ = s.orderDAO.UpdatePrepayID(order.ID, params.PrepayID)
	return params, nil
}

// CancelOrder 取消订单
func (s *OrderService) CancelOrder(userID uint64, orderID uint64) error {
	order, err := s.orderDAO.GetByID(orderID)
	if err != nil {
		return err
	}
	if order.UserID != userID {
		return errors.New("order not belong to user")
	}
	if order.Status != 0 {
		return errors.New("only unpaid order can be cancelled")
	}
	return s.orderDAO.UpdateStatus(orderID, 4)
}

// PayCallback 酒水订单支付回调处理（CAS）
// 返回 order 以便上层在真正成交时回调佣金计算
func (s *OrderService) PayCallback(orderNo, transactionID, rawNotify string, paidFen int64) (order *model.DrinkOrder, changed bool, err error) {
	order, err = s.orderDAO.GetByOrderNo(orderNo)
	if err != nil {
		return nil, false, err
	}
	// 金额一致性校验：防止失配 / 被篡改
	expected := int64(math.Round(order.PayAmount * 100))
	if paidFen > 0 && expected > 0 && paidFen != expected {
		return order, false, fmt.Errorf("amount mismatch: expect %d, got %d", expected, paidFen)
	}
	if order.Status != 0 {
		return order, false, nil
	}
	affected, err := s.orderDAO.MarkPaidCAS(orderNo, transactionID, rawNotify)
	if err != nil {
		return order, false, err
	}
	if affected == 0 {
		return order, false, nil
	}
	order.Status = 1
	if order.TransactionID == nil {
		txCopy := transactionID
		order.TransactionID = &txCopy
	}
	return order, true, nil
}

// round2 四舍五入保留 2 位小数
func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
