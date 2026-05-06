package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"nosocial/config"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/jwt"
	"nosocial/internal/pkg/utils"
	"nosocial/internal/pkg/wx"
	"strings"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	userDAO *dao.UserDAO
	wxCfg   *config.WXConfig
}

func NewUserService(userDAO *dao.UserDAO, wxCfg *config.WXConfig) *UserService {
	return &UserService{userDAO: userDAO, wxCfg: wxCfg}
}

// WXLoginReq 微信登录请求
type WXLoginReq struct {
	Code       string `json:"code"`
	InviteCode string `json:"invite_code,omitempty"`
}

// WXLoginResp 微信登录响应
type WXLoginResp struct {
	Token    string    `json:"token"`
	Expire   int       `json:"expire"`
	UserInfo *UserInfo `json:"user_info"`
}

type UserInfo struct {
	ID            uint64  `json:"id"`
	Nickname      *string `json:"nickname,omitempty"`
	Avatar        *string `json:"avatar,omitempty"`
	IsShareholder int8    `json:"is_shareholder"`
	IsNewUser     bool    `json:"is_new_user"`
}

func (s *UserService) WXLogin(req *WXLoginReq) (*WXLoginResp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sess, err := wx.JSCode2Session(ctx, s.wxCfg.AppID, s.wxCfg.Secret, req.Code)
	if err != nil {
		log.Printf("[UserService.WXLogin] code2session failed: %v", err)
		return nil, fmt.Errorf("invalid wx code")
	}
	openid := sess.OpenID

	user, err := s.userDAO.GetByOpenid(openid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	isNewUser := false
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 新用户注册
		isNewUser = true
		user = &model.User{
			Openid: openid,
			Status: 1,
		}
		if err := s.userDAO.Create(user); err != nil {
			return nil, err
		}

		// 处理邀请码绑定
		if req.InviteCode != "" {
			s.bindInviteCode(user.ID, req.InviteCode)
		}
	}

	token, err := jwt.GenerateWXToken(user.ID, user.Openid)
	if err != nil {
		return nil, err
	}

	return &WXLoginResp{
		Token:  token,
		Expire: config.C.JWT.WXExpire,
		UserInfo: &UserInfo{
			ID:            user.ID,
			Nickname:      user.Nickname,
			Avatar:        user.Avatar,
			IsShareholder: user.IsShareholder,
			IsNewUser:     isNewUser,
		},
	}, nil
}

func (s *UserService) bindInviteCode(userID uint64, inviteCode string) {
	// 已绑定不再重复绑定
	self, err := s.userDAO.GetByID(userID)
	if err != nil {
		return
	}
	if self.ParentID != nil && *self.ParentID > 0 {
		return
	}
	parent, err := s.userDAO.GetByInviteCode(inviteCode)
	if err != nil {
		return
	}
	// 环检测：禁止自我邻请和间接循环（parentPath 包含自身）
	if parent.ID == userID {
		return
	}
	selfMarker := fmt.Sprintf("/%d/", userID)
	if parent.ParentPath != nil && strings.Contains(*parent.ParentPath, selfMarker) {
		return
	}
	parentPath := fmt.Sprintf("/%d/", parent.ID)
	if parent.ParentPath != nil && *parent.ParentPath != "" {
		parentPath = *parent.ParentPath + fmt.Sprintf("%d/", parent.ID)
	}
	// 长度保护（数据库字段为 varchar(500)）
	if len(parentPath) > 500 {
		return
	}
	s.userDAO.UpdateParent(userID, parent.ID, parentPath)
}

func (s *UserService) GetProfile(userID uint64) (*model.User, error) {
	return s.userDAO.GetByID(userID)
}

func (s *UserService) UpdateBirthday(userID uint64, birthday string) error {
	// 严格校验 YYYY-MM-DD 格式，避免写入非法字符串
	t, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return fmt.Errorf("生日格式错误，应为 YYYY-MM-DD")
	}
	// 合理性校验：不允许未来日期，且不允许早于 1900-01-01
	now := time.Now()
	if t.After(now) {
		return fmt.Errorf("生日不能是未来日期")
	}
	if t.Year() < 1900 {
		return fmt.Errorf("生日年份不合法")
	}
	return s.userDAO.UpdateBirthday(userID, birthday)
}

func (s *UserService) GetUserList(offset, limit int, search string) ([]*model.User, int64, error) {
	return s.userDAO.List(offset, limit, search)
}

func (s *UserService) UpdateUserStatus(userID uint64, status int8) error {
	return s.userDAO.UpdateStatus(userID, status)
}

// ClaimFreeDrink 股东领取免费酒水（并发安全 CAS）
func (s *UserService) ClaimFreeDrink(userID uint64) error {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		return err
	}
	if user.IsShareholder != 1 {
		return fmt.Errorf("only shareholder can claim free drink")
	}
	if user.ShareholderExpireAt != nil && user.ShareholderExpireAt.Before(time.Now()) {
		return fmt.Errorf("股东身份已过期，请续费后再领取")
	}
	if user.FreeDrinkUsed == 1 {
		return fmt.Errorf("free drink already claimed")
	}
	affected, err := s.userDAO.ClaimFreeDrinkCAS(userID)
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("free drink already claimed")
	}
	return nil
}

// SetAsShareholder 管理员手动设置用户为共享股东
func (s *UserService) SetAsShareholder(userID uint64) error {
	user, err := s.userDAO.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("用户不存在")
		}
		return err
	}
	if user.IsShareholder == 1 {
		return fmt.Errorf("该用户已是股东")
	}

	// 生成唯一邀请码（重试最多5次）
	var inviteCode string
	for i := 0; i < 5; i++ {
		code := utils.GenerateInviteCode()
		_, err := s.userDAO.GetByInviteCode(code)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			inviteCode = code
			break
		}
	}
	if inviteCode == "" {
		return fmt.Errorf("邀请码生成失败，请重试")
	}

	return s.userDAO.SetAsShareholder(userID, inviteCode)
}
