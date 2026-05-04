package service

import (
	"errors"
	"nosocial/internal/bootstrap"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/jwt"
	"nosocial/internal/pkg/utils"

	"go.uber.org/zap"
)

type AdminService struct {
	adminDAO  *dao.AdminUserDAO
	userDAO   *dao.UserDAO
	orderDAO  *dao.DrinkOrderDAO
	reviewDAO *dao.ReviewDAO
}

func NewAdminService(adminDAO *dao.AdminUserDAO, userDAO *dao.UserDAO, orderDAO *dao.DrinkOrderDAO, reviewDAO *dao.ReviewDAO) *AdminService {
	return &AdminService{
		adminDAO:  adminDAO,
		userDAO:   userDAO,
		orderDAO:  orderDAO,
		reviewDAO: reviewDAO,
	}
}

type AdminLoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminLoginResp struct {
	Token    string `json:"token"`
	Expire   int    `json:"expire"`
	Nickname string `json:"nickname,omitempty"`
	Role     int8   `json:"role"`
}

func (s *AdminService) Login(req *AdminLoginReq) (*AdminLoginResp, error) {
	admin, err := s.adminDAO.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if admin.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	if !utils.CheckPassword(req.Password, admin.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := jwt.GenerateAdminToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		return nil, err
	}

	// 异步更新最后登录时间（非关键操作，不影响登录结果）
	go func() {
		if err := s.adminDAO.UpdateLoginTime(admin.ID); err != nil {
			if bootstrap.Log != nil {
				bootstrap.Log.Warn("admin update login time failed",
					zap.Uint32("admin_id", admin.ID), zap.Error(err))
			}
		}
	}()

	return &AdminLoginResp{
		Token:    token,
		Expire:   86400,
		Nickname: *admin.Nickname,
		Role:     admin.Role,
	}, nil
}

func (s *AdminService) GetDashboardStats() (map[string]interface{}, error) {
	todayAmount, todayCount, _ := s.orderDAO.GetTodayStats()

	_, newShareholderCount, _ := s.userDAO.ListShareholders(0, 9999, "")

	_, pendingReviewCount, _ := s.reviewDAO.List(0, 0, 9999)

	// 近7天营业额趋势
	revenueTrend := []dao.RevenueTrend{}
	if trend, err := s.orderDAO.GetRevenueTrend(7); err == nil && len(trend) > 0 {
		revenueTrend = trend
	}

	return map[string]interface{}{
		"today_amount":          todayAmount,
		"today_order_count":     todayCount,
		"new_shareholder_count": newShareholderCount,
		"pending_review_count":  pendingReviewCount,
		"revenue_trend":         revenueTrend,
	}, nil
}

func (s *AdminService) GetAdminList(page, size int) ([]*model.AdminUser, int64, error) {
	offset := (page - 1) * size
	return s.adminDAO.List(offset, size)
}

func (s *AdminService) CreateAdmin(admin *model.AdminUser) error {
	hashed, err := utils.HashPassword(admin.Password)
	if err != nil {
		return err
	}
	admin.Password = hashed
	return s.adminDAO.Create(admin)
}

func (s *AdminService) UpdateAdmin(admin *model.AdminUser) error {
	return s.adminDAO.Update(admin)
}

func (s *AdminService) DeleteAdmin(id uint32) error {
	return s.adminDAO.Delete(id)
}
