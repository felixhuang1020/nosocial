package wx

import (
	"nosocial/internal/middleware"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
}

func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

func (h *ReviewHandler) Submit(c *gin.Context) {
	var req struct {
		ScreenshotURL string   `json:"screenshot_url,omitempty"`
		ImageURLs     []string `json:"image_urls,omitempty"`
		Content       string   `json:"content,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 兼容旧协议：若未传 image_urls，用 screenshot_url 填充
	imgs := req.ImageURLs
	if len(imgs) == 0 && req.ScreenshotURL != "" {
		imgs = []string{req.ScreenshotURL}
	}
	if len(imgs) == 0 {
		response.BadRequest(c, "请上传至少一张截图")
		return
	}
	if len(imgs) > 9 {
		response.BadRequest(c, "最多上传 9 张截图")
		return
	}

	userID := middleware.GetWXUserID(c)
	review, err := h.reviewService.SubmitReview(userID, imgs, req.Content)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"review_id": review.ID,
		"status":    review.Status,
		"message":   "审核通过后将自动发放优惠券",
	})
}

func (h *ReviewHandler) Status(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	review, err := h.reviewService.GetReviewStatus(userID)
	if err != nil {
		response.Error(c, 1, "获取审核状态失败")
		return
	}
	if review == nil {
		response.Success(c, map[string]interface{}{
			"has_review": false,
		})
		return
	}
	response.Success(c, map[string]interface{}{
		"has_review": true,
		"review":     review,
	})
}
