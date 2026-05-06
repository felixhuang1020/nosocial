package wx

import (
	"nosocial/internal/middleware"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
)

type ShareholderHandler struct {
	shareholderService *service.ShareholderService
}

func NewShareholderHandler(shareholderService *service.ShareholderService) *ShareholderHandler {
	return &ShareholderHandler{shareholderService: shareholderService}
}

func (h *ShareholderHandler) Apply(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	order, err := h.shareholderService.ApplyShareholder(userID)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *ShareholderHandler) Pay(c *gin.Context) {
	// 注意：股东注册费由后端配置项 Business.ShareholderFee 决定，
	// 不接受客户端传入的金额，防止篡改支付金额。
	userID := middleware.GetWXUserID(c)
	// 先申请创建订单
	order, err := h.shareholderService.ApplyShareholder(userID)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}

	payParams, err := h.shareholderService.GetPayParams(userID, order.OrderNo)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"order_no":   order.OrderNo,
		"pay_params": payParams,
	})
}

func (h *ShareholderHandler) Profile(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	profile, err := h.shareholderService.GetProfile(userID)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, profile)
}

func (h *ShareholderHandler) Earnings(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	page := 1
	size := 20
	records, total, err := h.shareholderService.GetEarnings(userID, (page-1)*size, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  records,
		"total": total,
	})
}

func (h *ShareholderHandler) Team(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	team, err := h.shareholderService.GetTeam(userID)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  team,
		"total": len(team),
	})
}

func (h *ShareholderHandler) Withdraw(c *gin.Context) {
	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetWXUserID(c)
	if err := h.shareholderService.Withdraw(userID, req.Amount); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

// Withdrawals 提现记录列表
func (h *ShareholderHandler) Withdrawals(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	page := 1
	size := 20
	list, total, err := h.shareholderService.GetWithdrawals(userID, (page-1)*size, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  list,
		"total": total,
	})
}
