package wx

import (
	"nosocial/internal/middleware"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	couponService *service.CouponService
}

func NewCouponHandler(couponService *service.CouponService) *CouponHandler {
	return &CouponHandler{couponService: couponService}
}

func (h *CouponHandler) MyCoupons(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	status := int8(-1)
	if s := c.Query("status"); s != "" {
		status = int8([]byte(s)[0] - '0')
	}
	page := 1
	size := 20
	coupons, total, err := h.couponService.GetMyCoupons(userID, status, page, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  coupons,
		"total": total,
	})
}

func (h *CouponHandler) UseCoupon(c *gin.Context) {
	var req struct {
		CouponID uint64 `json:"coupon_id" binding:"required"`
		OrderID  uint64 `json:"order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetWXUserID(c)
	if err := h.couponService.UseCoupon(userID, req.CouponID, req.OrderID); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}
