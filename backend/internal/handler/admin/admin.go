package admin

import (
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// getPageSize 从 query 参数解析 page/size，默认 1/20
func getPageSize(c *gin.Context) (int, int) {
	page, size := 1, 20
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if s := c.Query("size"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			size = v
		}
	}
	return page, size
}

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var req service.AdminLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	resp, err := h.adminService.Login(&req)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, resp)
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	stats, err := h.adminService.GetDashboardStats()
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, stats)
}

func (h *AdminHandler) ListAdmins(c *gin.Context) {
	page, size := 1, 20
	admins, total, err := h.adminService.GetAdminList(page, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  admins,
		"total": total,
	})
}

func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var admin struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
		Role     int8   `json:"role"`
	}
	if err := c.ShouldBindJSON(&admin); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	response.Error(c, 1, "功能开发中")
}

func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	response.Error(c, 1, "功能开发中")
}

func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	response.Error(c, 1, "功能开发中")
}

// UserHandler 用户管理
type UserAdminHandler struct {
	userService *service.UserService
}

func NewUserAdminHandler(userService *service.UserService) *UserAdminHandler {
	return &UserAdminHandler{userService: userService}
}

func (h *UserAdminHandler) List(c *gin.Context) {
	page, size := getPageSize(c)
	users, total, err := h.userService.GetUserList((page-1)*size, size)
	if err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, map[string]interface{}{
		"list":  users,
		"total": total,
	})
}

func (h *UserAdminHandler) UpdateStatus(c *gin.Context) {
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err = h.userService.UpdateUserStatus(id, req.Status); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}
