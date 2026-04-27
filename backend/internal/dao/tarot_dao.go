package dao

import (
	"nosocial/internal/model"

	"gorm.io/gorm"
)

type TarotCardDAO struct {
	db *gorm.DB
}

func NewTarotCardDAO(db *gorm.DB) *TarotCardDAO {
	return &TarotCardDAO{db: db}
}

func (d *TarotCardDAO) List() ([]*model.TarotCard, error) {
	var cards []*model.TarotCard
	err := d.db.Order("card_no ASC").Find(&cards).Error
	return cards, err
}

func (d *TarotCardDAO) GetByCardNo(cardNo int) (*model.TarotCard, error) {
	var card model.TarotCard
	err := d.db.Where("card_no = ?", cardNo).First(&card).Error
	return &card, err
}

func (d *TarotCardDAO) Update(card *model.TarotCard) error {
	return d.db.Model(&model.TarotCard{}).Where("id = ?", card.ID).Omit("created_at").Updates(card).Error
}

// TarotReadingDAO 占卜记录DAO
type TarotReadingDAO struct {
	db *gorm.DB
}

func NewTarotReadingDAO(db *gorm.DB) *TarotReadingDAO {
	return &TarotReadingDAO{db: db}
}

func (d *TarotReadingDAO) Create(reading *model.TarotReading) error {
	return d.db.Create(reading).Error
}

func (d *TarotReadingDAO) ListByUser(userID uint64, offset, limit int) ([]*model.TarotReading, int64, error) {
	var readings []*model.TarotReading
	var total int64
	query := d.db.Where("user_id = ?", userID)
	query.Model(&model.TarotReading{}).Count(&total)
	err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&readings).Error
	return readings, total, err
}

// TarotDrinkMappingDAO 塔罗酒水映射DAO
type TarotDrinkMappingDAO struct {
	db *gorm.DB
}

func NewTarotDrinkMappingDAO(db *gorm.DB) *TarotDrinkMappingDAO {
	return &TarotDrinkMappingDAO{db: db}
}

func (d *TarotDrinkMappingDAO) GetByCardNo(cardNo int, isReversed int8) (*model.TarotDrinkMapping, error) {
	var mapping model.TarotDrinkMapping
	err := d.db.Where("card_no = ? AND is_reversed = ?", cardNo, isReversed).
		Order("match_score DESC").
		First(&mapping).Error
	return &mapping, err
}

func (d *TarotDrinkMappingDAO) ListByCardNo(cardNo int) ([]*model.TarotDrinkMapping, error) {
	var mappings []*model.TarotDrinkMapping
	err := d.db.Where("card_no = ?", cardNo).Order("match_score DESC").Find(&mappings).Error
	return mappings, err
}

func (d *TarotDrinkMappingDAO) ListAll() ([]*model.TarotDrinkMapping, error) {
	var mappings []*model.TarotDrinkMapping
	err := d.db.Find(&mappings).Error
	return mappings, err
}

func (d *TarotDrinkMappingDAO) GetByCardNoAndReversed(cardNo int, isReversed int8) (*model.TarotDrinkMapping, error) {
	var mapping model.TarotDrinkMapping
	err := d.db.Where("card_no = ? AND is_reversed = ?", cardNo, isReversed).First(&mapping).Error
	return &mapping, err
}

func (d *TarotDrinkMappingDAO) Update(mapping *model.TarotDrinkMapping) error {
	return d.db.Model(&model.TarotDrinkMapping{}).Where("id = ?", mapping.ID).Omit("created_at").Updates(mapping).Error
}

func (d *TarotDrinkMappingDAO) Create(mapping *model.TarotDrinkMapping) error {
	return d.db.Create(mapping).Error
}

func (d *TarotDrinkMappingDAO) Delete(id uint32) error {
	return d.db.Delete(&model.TarotDrinkMapping{}, id).Error
}
