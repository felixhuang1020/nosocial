package dao

import (
	"nosocial/internal/model"

	"gorm.io/gorm"
)

// BannerDAO
type BannerDAO struct {
	db *gorm.DB
}

func NewBannerDAO(db *gorm.DB) *BannerDAO {
	return &BannerDAO{db: db}
}

func (d *BannerDAO) Create(banner *model.Banner) error {
	return d.db.Create(banner).Error
}

func (d *BannerDAO) GetByID(id uint32) (*model.Banner, error) {
	var banner model.Banner
	err := d.db.First(&banner, id).Error
	return &banner, err
}

func (d *BannerDAO) ListByPosition(position int8) ([]*model.Banner, error) {
	var banners []*model.Banner
	err := d.db.Where("position = ? AND status = ?", position, 1).Order("sort_order ASC").Find(&banners).Error
	return banners, err
}

func (d *BannerDAO) ListAll(offset, limit int) ([]*model.Banner, int64, error) {
	var banners []*model.Banner
	var total int64
	d.db.Model(&model.Banner{}).Count(&total)
	err := d.db.Offset(offset).Limit(limit).Order("sort_order ASC").Find(&banners).Error
	return banners, total, err
}

func (d *BannerDAO) Update(banner *model.Banner) error {
	return d.db.Model(&model.Banner{}).Where("id = ?", banner.ID).Omit("created_at").Updates(banner).Error
}

func (d *BannerDAO) Delete(id uint32) error {
	return d.db.Delete(&model.Banner{}, id).Error
}

// AdminUserDAO
type AdminUserDAO struct {
	db *gorm.DB
}

func NewAdminUserDAO(db *gorm.DB) *AdminUserDAO {
	return &AdminUserDAO{db: db}
}

func (d *AdminUserDAO) Create(admin *model.AdminUser) error {
	return d.db.Create(admin).Error
}

func (d *AdminUserDAO) GetByID(id uint32) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := d.db.First(&admin, id).Error
	return &admin, err
}

func (d *AdminUserDAO) GetByUsername(username string) (*model.AdminUser, error) {
	var admin model.AdminUser
	err := d.db.Where("username = ?", username).First(&admin).Error
	return &admin, err
}

func (d *AdminUserDAO) UpdateLoginTime(id uint32) error {
	return d.db.Model(&model.AdminUser{}).Where("id = ?", id).Update("last_login_at", gorm.Expr("NOW()")).Error
}

func (d *AdminUserDAO) List(offset, limit int) ([]*model.AdminUser, int64, error) {
	var admins []*model.AdminUser
	var total int64
	d.db.Model(&model.AdminUser{}).Count(&total)
	err := d.db.Offset(offset).Limit(limit).Order("id DESC").Find(&admins).Error
	return admins, total, err
}

func (d *AdminUserDAO) Update(admin *model.AdminUser) error {
	return d.db.Model(&model.AdminUser{}).Where("id = ?", admin.ID).Omit("created_at").Updates(admin).Error
}

func (d *AdminUserDAO) Delete(id uint32) error {
	return d.db.Delete(&model.AdminUser{}, id).Error
}

// SettingDAO
type SettingDAO struct {
	db *gorm.DB
}

func NewSettingDAO(db *gorm.DB) *SettingDAO {
	return &SettingDAO{db: db}
}

func (d *SettingDAO) GetByKey(key string) (*model.Setting, error) {
	var setting model.Setting
	err := d.db.Where(map[string]interface{}{"key": key}).First(&setting).Error
	return &setting, err
}

func (d *SettingDAO) Set(key, value string) error {
	// 使用 PostgreSQL UPSERT 避免 check-then-act 竞态条件
	return d.db.Exec(`
		INSERT INTO settings (key, value, created_at, updated_at)
		VALUES (?, ?, NOW(), NOW())
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`, key, value).Error
}

func (d *SettingDAO) List() ([]*model.Setting, error) {
	var settings []*model.Setting
	err := d.db.Find(&settings).Error
	return settings, err
}
