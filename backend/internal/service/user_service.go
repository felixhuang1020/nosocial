package service

import (
	"context"
	"errors"
	"fmt"
	"nosocial/config"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"nosocial/internal/pkg/jwt"
	"nosocial/internal/pkg/wx"
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
	parent, err := s.userDAO.GetByInviteCode(inviteCode)
	if err != nil {
		return
	}
	parentPath := fmt.Sprintf("/%d/", parent.ID)
	if parent.ParentPath != nil {
		parentPath = *parent.ParentPath + fmt.Sprintf("%d/", parent.ID)
	}
	s.userDAO.UpdateParent(userID, parent.ID, parentPath)
}

func (s *UserService) GetProfile(userID uint64) (*model.User, error) {
	return s.userDAO.GetByID(userID)
}

func (s *UserService) UpdateBirthday(userID uint64, birthday string) error {
	return s.userDAO.UpdateBirthday(userID, birthday)
}

func (s *UserService) GetUserList(offset, limit int) ([]*model.User, int64, error) {
	return s.userDAO.List(offset, limit)
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
