package wx

import (
	"nosocial/internal/middleware"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) ClaimFreeDrink(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	if err := h.userService.ClaimFreeDrink(userID); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req service.WXLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	resp, err := h.userService.WXLogin(&req)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetWXUserID(c)
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) SetBirthday(c *gin.Context) {
	var req struct {
		Birthday string `json:"birthday" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetWXUserID(c)
	if err := h.userService.UpdateBirthday(userID, req.Birthday); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}
