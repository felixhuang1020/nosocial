package wx

import (
	"nosocial/internal/middleware"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req service.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetWXUserID(c)
	order, err := h.orderService.CreateOrder(userID, &req)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	page := 1
	size := 20
	orders, total, err := h.orderService.GetUserOrders(userID, page, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  orders,
		"total": total,
	})
}

func (h *OrderHandler) Pay(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的订单ID")
		return
	}

	userID := middleware.GetWXUserID(c)
	payParams, err := h.orderService.GetOrderPayParams(userID, id)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"pay_params": payParams,
	})
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的订单ID")
		return
	}

	userID := middleware.GetWXUserID(c)
	if err := h.orderService.CancelOrder(userID, id); err != nil {
		response.Error(c, 1, err.Error())
		return
	}

	response.Success(c, nil)
}
