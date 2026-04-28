package dao

import (
	"nosocial/internal/model"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (d *UserDAO) Create(user *model.User) error {
	return d.db.Create(user).Error
}

func (d *UserDAO) GetByID(id uint64) (*model.User, error) {
	var user model.User
	err := d.db.First(&user, id).Error
	return &user, err
}

func (d *UserDAO) GetByOpenid(openid string) (*model.User, error) {
	var user model.User
	err := d.db.Where("openid = ?", openid).First(&user).Error
	return &user, err
}

func (d *UserDAO) GetByInviteCode(code string) (*model.User, error) {
	var user model.User
	err := d.db.Where("invite_code = ?", code).First(&user).Error
	return &user, err
}

func (d *UserDAO) Update(user *model.User) error {
	return d.db.Model(&model.User{}).Where("id = ?", user.ID).Omit("created_at").Updates(user).Error
}

func (d *UserDAO) UpdateBirthday(userID uint64, birthday string) error {
	return d.db.Model(&model.User{}).Where("id = ?", userID).Update("birthday", birthday).Error
}

func (d *UserDAO) UpdateShareholderStatus(userID uint64, isShareholder int8, inviteCode string) error {
	updates := map[string]interface{}{
		"is_shareholder":    isShareholder,
		"register_fee_paid": 1,
		"invite_code":       inviteCode,
	}
	return d.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

func (d *UserDAO) UpdateParent(userID uint64, parentID uint64, parentPath string) error {
	return d.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"parent_id":   parentID,
		"parent_path": parentPath,
	}).Error
}

func (d *UserDAO) UpdateBalance(userID uint64, amount float64) error {
	return d.db.Model(&model.User{}).Where("id = ?", userID).Update("balance", gorm.Expr("balance + ?", amount)).Error
}

func (d *UserDAO) UpdateTotalEarning(userID uint64, amount float64) error {
	return d.db.Model(&model.User{}).Where("id = ?", userID).Update("total_earning", gorm.Expr("total_earning + ?", amount)).Error
}

func (d *UserDAO) ListShareholders(offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64
	query := d.db.Where("is_shareholder = ?", 1)
	query.Model(&model.User{}).Count(&total)
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (d *UserDAO) GetChildren(parentID uint64) ([]*model.User, error) {
	var users []*model.User
	err := d.db.Where("parent_id = ?", parentID).Find(&users).Error
	return users, err
}

func (d *UserDAO) GetChildrenCount(parentID uint64) (int64, error) {
	var count int64
	err := d.db.Model(&model.User{}).Where("parent_id = ?", parentID).Count(&count).Error
	return count, err
}

func (d *UserDAO) GetByBirthday(monthDay string) ([]*model.User, error) {
	var users []*model.User
	err := d.db.Where("TO_CHAR(birthday, 'MM-DD') = ? AND is_shareholder = 1", monthDay).Find(&users).Error
	return users, err
}

func (d *UserDAO) List(offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64
	d.db.Model(&model.User{}).Count(&total)
	err := d.db.Offset(offset).Limit(limit).Order("id DESC").Find(&users).Error
	return users, total, err
}

func (d *UserDAO) UpdateStatus(userID uint64, status int8) error {
	return d.db.Model(&model.User{}).Where("id = ?", userID).Update("status", status).Error
}

func (d *UserDAO) UpdateFreeDrinkUsed(userID uint64, used int8) error {
	return d.db.Model(&model.User{}).Where("id = ?", userID).Update("free_drink_used", used).Error
}

// ClaimFreeDrinkCAS 领取免费酒水的原子 CAS：仅在 is_shareholder=1 且 free_drink_used=0 时置 1
func (d *UserDAO) ClaimFreeDrinkCAS(userID uint64) (int64, error) {
	res := d.db.Model(&model.User{}).
		Where("id = ? AND is_shareholder = 1 AND free_drink_used = 0", userID).
		Update("free_drink_used", 1)
	return res.RowsAffected, res.Error
}
