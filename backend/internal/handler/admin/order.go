package admin

import (
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderAdminHandler struct {
	orderService *service.OrderService
}

func NewOrderAdminHandler(orderService *service.OrderService) *OrderAdminHandler {
	return &OrderAdminHandler{orderService: orderService}
}

func (h *OrderAdminHandler) List(c *gin.Context) {
	status := int8(-1)
	if s := c.Query("status"); s != "" {
		status = int8(s[0] - '0')
	}
	page, size := getPageSize(c)
	orders, total, err := h.orderService.GetOrderList(status, page, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  orders,
		"total": total,
	})
}

func (h *OrderAdminHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.orderService.UpdateOrderStatus(id, req.Status); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}
