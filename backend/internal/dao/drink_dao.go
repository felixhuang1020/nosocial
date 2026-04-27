package dao

import (
	"nosocial/internal/model"

	"gorm.io/gorm"
)

type DrinkCategoryDAO struct {
	db *gorm.DB
}

func NewDrinkCategoryDAO(db *gorm.DB) *DrinkCategoryDAO {
	return &DrinkCategoryDAO{db: db}
}

func (d *DrinkCategoryDAO) Create(category *model.DrinkCategory) error {
	return d.db.Create(category).Error
}

func (d *DrinkCategoryDAO) GetByID(id uint32) (*model.DrinkCategory, error) {
	var category model.DrinkCategory
	err := d.db.First(&category, id).Error
	return &category, err
}

func (d *DrinkCategoryDAO) List() ([]*model.DrinkCategory, error) {
	var categories []*model.DrinkCategory
	err := d.db.Where("status = ?", 1).Order("sort_order ASC").Find(&categories).Error
	return categories, err
}

func (d *DrinkCategoryDAO) Update(category *model.DrinkCategory) error {
	return d.db.Model(&model.DrinkCategory{}).Where("id = ?", category.ID).Omit("created_at").Updates(category).Error
}

func (d *DrinkCategoryDAO) Delete(id uint32) error {
	return d.db.Delete(&model.DrinkCategory{}, id).Error
}

// DrinkDAO 酒水DAO
type DrinkDAO struct {
	db *gorm.DB
}

func NewDrinkDAO(db *gorm.DB) *DrinkDAO {
	return &DrinkDAO{db: db}
}

func (d *DrinkDAO) Create(drink *model.Drink) error {
	return d.db.Create(drink).Error
}

func (d *DrinkDAO) GetByID(id uint64) (*model.Drink, error) {
	var drink model.Drink
	err := d.db.Preload("Category").First(&drink, id).Error
	return &drink, err
}

func (d *DrinkDAO) List(categoryID uint32, status int8, offset, limit int) ([]*model.Drink, int64, error) {
	var drinks []*model.Drink
	var total int64
	query := d.db.Model(&model.Drink{}).Preload("Category")
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id DESC").Find(&drinks).Error
	return drinks, total, err
}

func (d *DrinkDAO) ListRecommended() ([]*model.Drink, error) {
	var drinks []*model.Drink
	err := d.db.Where("is_recommended = ? AND status = ?", 1, 1).
		Preload("Category").
		Order("sort_order ASC").Find(&drinks).Error
	return drinks, err
}

func (d *DrinkDAO) Update(drink *model.Drink) error {
	return d.db.Model(&model.Drink{}).Where("id = ?", drink.ID).Omit("created_at").Updates(drink).Error
}

func (d *DrinkDAO) Delete(id uint64) error {
	return d.db.Delete(&model.Drink{}, id).Error
}
