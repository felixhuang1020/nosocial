package wx

import (
	"nosocial/internal/middleware"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
)

type TarotHandler struct {
	tarotService *service.TarotService
}

func NewTarotHandler(tarotService *service.TarotService) *TarotHandler {
	return &TarotHandler{tarotService: tarotService}
}

func (h *TarotHandler) Divine(c *gin.Context) {
	var req service.DivineReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetWXUserID(c)
	result, err := h.tarotService.Divine(userID, &req)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, result)
}

func (h *TarotHandler) History(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	page := 1
	size := 20
	readings, total, err := h.tarotService.GetHistory(userID, (page-1)*size, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  readings,
		"total": total,
	})
}
