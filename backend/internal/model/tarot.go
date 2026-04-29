package model

import (
	"time"
)

// TarotCard 塔罗牌库
type TarotCard struct {
	ID              uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	CardNo          int       `gorm:"not null;uniqueIndex:uk_card_no" json:"card_no"`
	Name            string    `gorm:"type:varchar(32);not null" json:"name"`
	NameEn          *string   `gorm:"type:varchar(64)" json:"name_en,omitempty"`
	ArcanaType      int8      `gorm:"type:smallint;not null" json:"arcana_type"`
	Suit            *string   `gorm:"type:varchar(16)" json:"suit,omitempty"`
	ImageURL        string    `gorm:"type:varchar(255);not null" json:"image_url"`
	UprightMeaning  *string   `gorm:"type:text" json:"upright_meaning,omitempty"`
	ReversedMeaning *string   `gorm:"type:text" json:"reversed_meaning,omitempty"`
	Keywords        *string   `gorm:"type:varchar(255)" json:"keywords,omitempty"`
	Element         *string   `gorm:"type:varchar(16)" json:"element,omitempty"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (TarotCard) TableName() string {
	return "tarot_cards"
}

// TarotReading 占卜记录
type TarotReading struct {
	ID                 uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID             uint64    `gorm:"not null;index:idx_user_id" json:"user_id"`
	Question           *string   `gorm:"type:varchar(255)" json:"question,omitempty"`
	SpreadType         int8      `gorm:"type:smallint;not null;default:1" json:"spread_type"`
	Cards              string    `gorm:"type:json;not null" json:"cards"`
	RecommendedDrinkID uint64    `gorm:"not null" json:"recommended_drink_id"`
	RecommendedReason  *string   `gorm:"type:varchar(255)" json:"recommended_reason,omitempty"`
	ReadingDate        string    `gorm:"type:date;not null;index:idx_reading_date" json:"reading_date"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (TarotReading) TableName() string {
	return "tarot_readings"
}

// TarotDrinkMapping 塔罗-酒水映射
// 业务语义：每张牌 (CardNo) 每个朝向 (IsReversed) 只保留一条映射
// 唯一键：(card_no, is_reversed)
type TarotDrinkMapping struct {
	ID             uint32    `gorm:"primaryKey;autoIncrement" json:"id"`
	CardNo         int       `gorm:"not null;uniqueIndex:uk_card_reversed,priority:1" json:"card_no"`
	DrinkID        uint64    `gorm:"not null" json:"drink_id"`
	IsReversed     int8      `gorm:"type:smallint;not null;default:0;uniqueIndex:uk_card_reversed,priority:2" json:"is_reversed"`
	MatchScore     int       `gorm:"not null;default:100" json:"match_score"`
	ReasonTemplate *string   `gorm:"type:varchar(255)" json:"reason_template,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (TarotDrinkMapping) TableName() string {
	return "tarot_drink_mappings"
}
