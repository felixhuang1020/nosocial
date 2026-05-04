package wx

import (
	"html"
	"net/url"
	"nosocial/internal/middleware"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
)

// validateImageURL 验证图片 URL 合法性：仅允许 https + 指定 OSS 域名
func validateImageURL(rawURL string) bool {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	if parsed.Scheme != "https" {
		return false
	}
	if parsed.Host != "nosocial.oss-cn-beijing.aliyuncs.com" {
		return false
	}
	return true
}

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

	// 验证图片 URL 合法性
	for _, imgURL := range imgs {
		if !validateImageURL(imgURL) {
			response.BadRequest(c, "图片地址不合法")
			return
		}
	}

	// 评价内容长度限制
	if len([]rune(req.Content)) > 500 {
		response.BadRequest(c, "评价内容不能超过500字")
		return
	}

	// XSS 防护：转义 HTML 特殊字符
	safeContent := html.EscapeString(req.Content)

	userID := middleware.GetWXUserID(c)
	review, err := h.reviewService.SubmitReview(userID, imgs, safeContent)
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
