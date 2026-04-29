package service

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"nosocial/internal/dao"
	"nosocial/internal/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TarotService struct {
	cardDAO    *dao.TarotCardDAO
	readingDAO *dao.TarotReadingDAO
	mappingDAO *dao.TarotDrinkMappingDAO
	drinkDAO   *dao.DrinkDAO
	db         *gorm.DB
}

func NewTarotService(cardDAO *dao.TarotCardDAO, readingDAO *dao.TarotReadingDAO, mappingDAO *dao.TarotDrinkMappingDAO, drinkDAO *dao.DrinkDAO, db *gorm.DB) *TarotService {
	return &TarotService{
		cardDAO:    cardDAO,
		readingDAO: readingDAO,
		mappingDAO: mappingDAO,
		drinkDAO:   drinkDAO,
		db:         db,
	}
}

// secureIntn 安全随机数，相对于 math/rand 不可预测。fallback 到 time 仅当 crypto 不可用。
func secureIntn(n int) int {
	if n <= 0 {
		return 0
	}
	var b [8]byte
	if _, err := cryptorand.Read(b[:]); err != nil {
		return int(time.Now().UnixNano() % int64(n))
	}
	v := binary.BigEndian.Uint64(b[:])
	return int(v % uint64(n))
}

// secureBoolWithProb 按指定概率返回 true（只接受 0-1）
func secureBoolWithProb(p float32) bool {
	if p <= 0 {
		return false
	}
	if p >= 1 {
		return true
	}
	var b [4]byte
	if _, err := cryptorand.Read(b[:]); err != nil {
		return false
	}
	u := binary.BigEndian.Uint32(b[:])
	threshold := uint32(p * float32(^uint32(0)))
	return u < threshold
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
	// 1. 安全随机：0-77 卡号 + 20% 逆位概率
	cardNo := secureIntn(78)
	isReversed := secureBoolWithProb(0.2)

	// 2. 查询牌信息
	card, err := s.cardDAO.GetByCardNo(cardNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("塔罗牌库尚未初始化，请联系管理员")
		}
		return nil, err
	}

	// 3. 查询映射关系（容错：无映射时降级返回牌片 + 空推荐）
	isRevInt := int8(0)
	if isReversed {
		isRevInt = 1
	}
	mapping, err := s.mappingDAO.GetByCardNo(cardNo, isRevInt)
	var drink *model.Drink
	mappingOK := false
	if err == nil {
		mappingOK = true
		d, derr := s.drinkDAO.GetByID(mapping.DrinkID)
		if derr == nil && d != nil && d.Status == 1 {
			drink = d
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 4. 生成含义
	meaning := ""
	if card.UprightMeaning != nil {
		meaning = *card.UprightMeaning
	}
	if isReversed && card.ReversedMeaning != nil {
		meaning = *card.ReversedMeaning
	}

	resp := &DivineResp{
		Cards: []CardResult{
			{
				CardNo:     cardNo,
				Name:       card.Name,
				IsReversed: isReversed,
				ImageURL:   card.ImageURL,
				Meaning:    meaning,
			},
		},
	}

	// 5. 若映射或酒水缺失：不写入 reading（RecommendedDrinkID NOT NULL），返回空推荐
	if !mappingOK || drink == nil {
		return resp, nil
	}

	// 6. 生成推荐理由
	reason := ""
	if mapping.ReasonTemplate != nil && *mapping.ReasonTemplate != "" {
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

	// 7. 保存占卜记录
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

	resp.ReadingID = reading.ID
	resp.Recommend = &DrinkRecommend{
		DrinkID:     drink.ID,
		Name:        drink.Name,
		EnglishName: drink.EnglishName,
		ImageURL:    drink.ImageURL,
		Price:       drink.Price,
		Reason:      reason,
		Alcohol:     drink.Alcohol,
	}
	return resp, nil
}

func (s *TarotService) GetHistory(userID uint64, offset, limit int) ([]*model.TarotReading, int64, error) {
	return s.readingDAO.ListByUser(userID, offset, limit)
}

func (s *TarotService) GetCards() ([]*model.TarotCard, error) {
	return s.cardDAO.List()
}

func (s *TarotService) UpdateCard(id uint32, imageURL string) error {
	res := s.db.Model(&model.TarotCard{}).Where("id = ?", id).Update("image_url", imageURL)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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
	if len(items) == 0 {
		return nil
	}
	// 1. 基础校验：CardNo 范围 + DrinkID 非零 + MatchScore 范围
	drinkIDSet := make(map[uint64]struct{}, len(items))
	for _, it := range items {
		if it.CardNo < 0 || it.CardNo > 77 {
			return fmt.Errorf("非法卡号 card_no=%d", it.CardNo)
		}
		if it.DrinkID == 0 {
			return fmt.Errorf("card_no=%d 的 drink_id 不能为空", it.CardNo)
		}
		if it.MatchScore < 0 || it.MatchScore > 100 {
			return fmt.Errorf("card_no=%d 的 match_score 必须在 0-100 之间", it.CardNo)
		}
		drinkIDSet[it.DrinkID] = struct{}{}
	}
	// 2. 批量校验 DrinkID 存在性且为上架状态
	drinkIDs := make([]uint64, 0, len(drinkIDSet))
	for id := range drinkIDSet {
		drinkIDs = append(drinkIDs, id)
	}
	drinks, err := s.drinkDAO.ListByIDs(drinkIDs)
	if err != nil {
		return err
	}
	existing := make(map[uint64]struct{}, len(drinks))
	for _, d := range drinks {
		existing[d.ID] = struct{}{}
	}
	for id := range drinkIDSet {
		if _, ok := existing[id]; !ok {
			return fmt.Errorf("drink_id=%d 不存在或已下架", id)
		}
	}

	// 3. 事务内 UPSERT（ON CONFLICT(card_no, is_reversed) DO UPDATE）
	return s.db.Transaction(func(tx *gorm.DB) error {
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
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "card_no"}, {Name: "is_reversed"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"drink_id", "match_score", "reason_template",
				}),
			}).Create(mapping).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
