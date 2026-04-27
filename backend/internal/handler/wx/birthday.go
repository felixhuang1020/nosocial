package wx

import (
	"nosocial/internal/middleware"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
)

type BirthdayHandler struct {
	birthdayService *service.BirthdayService
}

func NewBirthdayHandler(birthdayService *service.BirthdayService) *BirthdayHandler {
	return &BirthdayHandler{birthdayService: birthdayService}
}

func (h *BirthdayHandler) CheckGift(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	gift, err := h.birthdayService.CheckGift(userID)
	if err != nil {
		response.Success(c, map[string]interface{}{"has_gift": false})
		return
	}
	response.Success(c, map[string]interface{}{"has_gift": true, "gift": gift})
}

func (h *BirthdayHandler) Claim(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	if err := h.birthdayService.ClaimGift(userID); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *BirthdayHandler) History(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	history, err := h.birthdayService.GetHistory(userID)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, history)
}
