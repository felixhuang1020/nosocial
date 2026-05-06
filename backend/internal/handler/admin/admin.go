package admin

import (
	"context"
	"net/http"
	"nosocial/config"
	"nosocial/internal/pkg/jwt"
	"nosocial/internal/pkg/response"
	"nosocial/internal/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// getPageSize 从 query 参数解析 page/size，默认 1/20，上限 100
// 封顶是为了防止恶意传入 size=1000000 抜库 / 打爆 Response。
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
	if size > 100 {
		size = 100
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

	// 设置 httpOnly Cookie，开发环境不强制 HTTPS
	cfg := config.C
	isSecure := cfg.App.Mode == "release"
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("admin_token", resp.Token, cfg.JWT.AdminExpire, "/api/v1/admin", "", isSecure, true)

	response.Success(c, resp)
}

func (h *AdminHandler) Logout(c *gin.Context) {
	// 1. 尝试获取当前 token 进行吊销
	var tokenString string
	if cookie, err := c.Cookie("admin_token"); err == nil && cookie != "" {
		tokenString = cookie
	} else {
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		}
	}

	// 2. 吊销 token（解析出 jti 后加入黑名单）
	if tokenString != "" {
		if claims, err := jwt.ParseAdminToken(tokenString); err == nil {
			cfg := config.C
			ttl := time.Duration(cfg.JWT.AdminExpire) * time.Second
			_ = jwt.Revoke(context.Background(), claims.ID, ttl)
		}
	}

	// 3. 清除 Cookie
	isSecure := config.C.App.Mode == "release"
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("admin_token", "", -1, "/api/v1/admin", "", isSecure, true)

	response.Success(c, nil)
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
	search := c.Query("search")
	users, total, err := h.userService.GetUserList((page-1)*size, size, search)
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

func (h *UserAdminHandler) SetShareholder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err = h.userService.SetAsShareholder(id); err != nil {
		response.Error(c, 1, err.Error())
		return
	}
	response.Success(c, nil)
}
