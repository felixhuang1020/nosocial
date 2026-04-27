package service

import (
	"nosocial/config"
	"nosocial/internal/dao"
	"strconv"
)

type SystemService struct {
	settingDAO *dao.SettingDAO
}

func NewSystemService(settingDAO *dao.SettingDAO) *SystemService {
	return &SystemService{settingDAO: settingDAO}
}

// GetConfig 获取业务配置
func (s *SystemService) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"shareholder_fee":          config.C.Business.ShareholderFee,
		"commission_rate":          config.C.Business.CommissionRate,
		"free_drink_id":            config.C.Business.FreeDrinkID,
		"review_coupon_amount":     config.C.Business.ReviewCouponAmount,
		"review_coupon_min_order":  config.C.Business.ReviewCouponMinOrder,
		"review_coupon_valid_days": config.C.Business.ReviewCouponValidDays,
	}
}

// GetShopInfo 获取门店信息
func (s *SystemService) GetShopInfo() map[string]interface{} {
	// 优先从settings表读取，不存在则使用默认值
	shop := map[string]interface{}{
		"name":           "NoSocial 酒吧",
		"address":        "北京市朝阳区三里屯太古里北区 N8-20",
		"phone":          "010-8888-6666",
		"business_hours": "周一至周日 18:00 - 04:00",
		"wifi_name":      "NoSocial_Free",
		"wifi_password":  "nosocial888",
		"latitude":       39.934,
		"longitude":      116.455,
	}

	settings, err := s.settingDAO.List()
	if err != nil {
		return shop
	}

	for _, st := range settings {
		switch st.Key {
		case "shop_name":
			shop["name"] = st.Value
		case "shop_address":
			shop["address"] = st.Value
		case "shop_phone":
			shop["phone"] = st.Value
		case "shop_business_hours":
			shop["business_hours"] = st.Value
		case "shop_wifi_name":
			shop["wifi_name"] = st.Value
		case "shop_wifi_password":
			shop["wifi_password"] = st.Value
		case "shop_latitude":
			if v, err := strconv.ParseFloat(st.Value, 64); err == nil {
				shop["latitude"] = v
			}
		case "shop_longitude":
			if v, err := strconv.ParseFloat(st.Value, 64); err == nil {
				shop["longitude"] = v
			}
		}
	}

	return shop
}

// GetAllSettings 获取所有系统配置（管理端用）
func (s *SystemService) GetAllSettings() map[string]interface{} {
	// 默认值（来自 config.yaml）
	result := s.GetConfig()
	shop := s.GetShopInfo()
	for k, v := range shop {
		result[k] = v
	}

	// 从数据库覆盖（业务配置）
	settings, err := s.settingDAO.List()
	if err == nil {
		for _, st := range settings {
			switch st.Key {
			case "shareholder_fee":
				if v, err := strconv.ParseFloat(st.Value, 64); err == nil {
					result["shareholder_fee"] = v
				}
			case "commission_rate":
				if v, err := strconv.ParseFloat(st.Value, 64); err == nil {
					result["commission_rate"] = v
				}
			case "free_drink_id":
				if v, err := strconv.ParseUint(st.Value, 10, 64); err == nil {
					result["free_drink_id"] = v
				}
			case "review_coupon_amount":
				if v, err := strconv.ParseFloat(st.Value, 64); err == nil {
					result["review_coupon_amount"] = v
				}
			case "review_coupon_min_order":
				if v, err := strconv.ParseFloat(st.Value, 64); err == nil {
					result["review_coupon_min_order"] = v
				}
			case "review_coupon_valid_days":
				if v, err := strconv.Atoi(st.Value); err == nil {
					result["review_coupon_valid_days"] = v
				}
			case "wx_appid":
				result["wx_appid"] = st.Value
			case "wx_mch_id":
				result["wx_mch_id"] = st.Value
			}
		}
	}

	return result
}

// UpdateSettings 批量更新系统配置（管理端用）
func (s *SystemService) UpdateSettings(updates map[string]string) error {
	for key, val := range updates {
		if err := s.settingDAO.Set(key, val); err != nil {
			return err
		}
	}
	return nil
}

// InitDefaultSettings 初始化默认配置
func (s *SystemService) InitDefaultSettings() error {
	defaults := map[string]string{
		"shop_name":           "NoSocial 酒吧",
		"shop_address":        "北京市朝阳区三里屯太古里北区 N8-20",
		"shop_phone":          "010-8888-6666",
		"shop_business_hours": "周一至周日 18:00 - 04:00",
		"shop_wifi_name":      "NoSocial_Free",
		"shop_wifi_password":  "nosocial888",
		"shop_latitude":       "39.934",
		"shop_longitude":      "116.455",
	}

	for key, value := range defaults {
		_, err := s.settingDAO.GetByKey(key)
		if err != nil {
			// 不存在则创建
			if err := s.settingDAO.Set(key, value); err != nil {
				return err
			}
		}
	}

	return nil
}
