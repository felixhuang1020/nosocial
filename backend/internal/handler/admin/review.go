package admin

import (
	"nosocial/internal/middleware"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewAdminHandler struct {
	reviewService *service.ReviewService
}

func NewReviewAdminHandler(reviewService *service.ReviewService) *ReviewAdminHandler {
	return &ReviewAdminHandler{reviewService: reviewService}
}

func (h *ReviewAdminHandler) List(c *gin.Context) {
	status := int8(-1)
	if s := c.Query("status"); s != "" {
		if v, err := strconv.Atoi(s); err == nil {
			status = int8(v)
		}
	}
	page, size := getPageSize(c)
	reviews, total, err := h.reviewService.GetReviewList(status, page, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  reviews,
		"total": total,
	})
}

func (h *ReviewAdminHandler) Audit(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Pass         bool   `json:"pass"`
		RejectReason string `json:"reject_reason,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// 驳回理由长度限制，避免数据库字段溢出或恶意超长输入
	if len(req.RejectReason) > 500 {
		response.BadRequest(c, "驳回理由过长（上限 500 字）")
		return
	}

	adminID := middleware.GetAdminID(c)
	if err := h.reviewService.AuditReview(id, adminID, req.Pass, req.RejectReason); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}
