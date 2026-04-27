package service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"time"
)

type TarotService struct {
	cardDAO    *dao.TarotCardDAO
	readingDAO *dao.TarotReadingDAO
	mappingDAO *dao.TarotDrinkMappingDAO
	drinkDAO   *dao.DrinkDAO
}

func NewTarotService(cardDAO *dao.TarotCardDAO, readingDAO *dao.TarotReadingDAO, mappingDAO *dao.TarotDrinkMappingDAO, drinkDAO *dao.DrinkDAO) *TarotService {
	return &TarotService{
		cardDAO:    cardDAO,
		readingDAO: readingDAO,
		mappingDAO: mappingDAO,
		drinkDAO:   drinkDAO,
	}
}

type DivineReq struct {
	Question   string `json:"question,omitempty"`
	SpreadType int8   `json:"spread_type"`
}

type DivineResp struct {
	ReadingID uint64          `json:"reading_id"`
	Cards     []CardResult    `json:"cards"`
	Recommend *DrinkRecommend `json:"recommend"`
}

type CardResult struct {
	CardNo     int    `json:"card_no"`
	Name       string `json:"name"`
	IsReversed bool   `json:"is_reversed"`
	ImageURL   string `json:"image_url"`
	Meaning    string `json:"meaning"`
}

type DrinkRecommend struct {
	DrinkID     uint64   `json:"drink_id"`
	Name        string   `json:"name"`
	EnglishName *string  `json:"english_name,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty"`
	Price       float64  `json:"price"`
	Reason      string   `json:"reason"`
	Alcohol     *float64 `json:"alcohol,omitempty"`
}

func (s *TarotService) Divine(userID uint64, req *DivineReq) (*DivineResp, error) {
	// 1. 生成随机种子
	seed := time.Now().UnixNano() + int64(userID)
	r := rand.New(rand.NewSource(seed))

	// 2. 抽取1张牌（0-77）
	cardNo := r.Intn(78)

	// 3. 20%概率逆位
	isReversed := r.Float32() < 0.2

	// 4. 查询牌信息
	card, err := s.cardDAO.GetByCardNo(cardNo)
	if err != nil {
		return nil, err
	}

	// 5. 查询映射关系
	isRevInt := int8(0)
	if isReversed {
		isRevInt = 1
	}
	mapping, err := s.mappingDAO.GetByCardNo(cardNo, isRevInt)
	if err != nil {
		return nil, err
	}

	// 6. 查询酒水
	drink, err := s.drinkDAO.GetByID(mapping.DrinkID)
	if err != nil {
		return nil, err
	}

	// 7. 生成含义
	meaning := ""
	if card.UprightMeaning != nil {
		meaning = *card.UprightMeaning
	}
	if isReversed && card.ReversedMeaning != nil {
		meaning = *card.ReversedMeaning
	}

	// 8. 生成推荐理由
	reason := ""
	if mapping.ReasonTemplate != nil {
		reason = *mapping.ReasonTemplate
	} else {
		position := "正位"
		if isReversed {
			position = "逆位"
		}
		keywords := ""
		if card.Keywords != nil {
			keywords = *card.Keywords
		}
		reason = fmt.Sprintf("【%s】%s代表%s。这杯%s正适合当下的你。", card.Name, position, keywords, drink.Name)
	}

	// 9. 保存记录
	cardsJSON, _ := json.Marshal([]map[string]interface{}{
		{
			"card_no":     cardNo,
			"is_reversed": isReversed,
			"position":    "今日运势",
		},
	})

	reading := &model.TarotReading{
		UserID:             userID,
		Question:           &req.Question,
		SpreadType:         req.SpreadType,
		Cards:              string(cardsJSON),
		RecommendedDrinkID: drink.ID,
		RecommendedReason:  &reason,
		ReadingDate:        time.Now().Format("2006-01-02"),
	}
	if err := s.readingDAO.Create(reading); err != nil {
		return nil, err
	}

	return &DivineResp{
		ReadingID: reading.ID,
		Cards: []CardResult{
			{
				CardNo:     cardNo,
				Name:       card.Name,
				IsReversed: isReversed,
				ImageURL:   card.ImageURL,
				Meaning:    meaning,
			},
		},
		Recommend: &DrinkRecommend{
			DrinkID:     drink.ID,
			Name:        drink.Name,
			EnglishName: drink.EnglishName,
			ImageURL:    drink.ImageURL,
			Price:       drink.Price,
			Reason:      reason,
			Alcohol:     drink.Alcohol,
		},
	}, nil
}

func (s *TarotService) GetHistory(userID uint64, offset, limit int) ([]*model.TarotReading, int64, error) {
	return s.readingDAO.ListByUser(userID, offset, limit)
}

func (s *TarotService) GetCards() ([]*model.TarotCard, error) {
	return s.cardDAO.List()
}

func (s *TarotService) GetAllMappings() ([]*model.TarotDrinkMapping, error) {
	return s.mappingDAO.ListAll()
}

type MappingItem struct {
	CardNo         int    `json:"card_no"`
	DrinkID        uint64 `json:"drink_id"`
	IsReversed     bool   `json:"is_reversed"`
	MatchScore     int    `json:"match_score"`
	ReasonTemplate string `json:"reason_template"`
}

func (s *TarotService) UpdateMappings(items []MappingItem) error {
	for _, item := range items {
		isRev := int8(0)
		if item.IsReversed {
			isRev = 1
		}
		reason := item.ReasonTemplate
		mapping := &model.TarotDrinkMapping{
			CardNo:         item.CardNo,
			DrinkID:        item.DrinkID,
			IsReversed:     isRev,
			MatchScore:     item.MatchScore,
			ReasonTemplate: &reason,
		}
		// 查找是否已存在
		existing, err := s.mappingDAO.GetByCardNoAndReversed(item.CardNo, isRev)
		if err == nil {
			// 更新
			mapping.ID = existing.ID
			if err := s.mappingDAO.Update(mapping); err != nil {
				return err
			}
		} else {
			// 创建
			if err := s.mappingDAO.Create(mapping); err != nil {
				return err
			}
		}
	}
	return nil
}
