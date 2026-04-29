package main

// 塔罗牌 78 张基础数据
// 大阿卡纳 22 张 (card_no 0-21)
// 小阿卡纳 56 张 (card_no 22-77) = 4 花色 × 14 张
// 幂等：已存在则跳过；image_url 初始为空，由管理员在后台上传

import (
	"log"
	"nosocial/internal/model"

	"gorm.io/gorm"
)

type tarotSeed struct {
	CardNo     int
	Name       string
	NameEn     string
	ArcanaType int8 // 1=大阿卡纳 2=小阿卡纳
	Suit       string
	Element    string
	Keywords   string
	Upright    string
	Reversed   string
}

// 大阿卡纳 0-21
var majorArcana = []tarotSeed{
	{0, "愚者", "The Fool", 1, "", "风", "冒险,纯真,新开始", "怀抱无限可能迈出新的一步，保持天真与好奇心。", "鲁莽冲动，逃避现实，决策失误。"},
	{1, "魔术师", "The Magician", 1, "", "风", "创造,意志,行动", "将意志化为现实，拥有实现目标的全部资源。", "操弄他人，浪费才华，自我怀疑。"},
	{2, "女祭司", "The High Priestess", 1, "", "水", "直觉,神秘,内在", "倾听直觉与潜意识，静观其变。", "忽视直觉，被秘密困扰，表里不一。"},
	{3, "女皇", "The Empress", 1, "", "土", "丰饶,母性,感性", "感官富足与创造力涌现，关系圆融。", "过度依赖，创造力停滞，情感窒息。"},
	{4, "皇帝", "The Emperor", 1, "", "火", "权威,秩序,稳固", "建立稳定架构，展现领导力与纪律。", "专制独裁，顽固不化，失去控制。"},
	{5, "教皇", "The Hierophant", 1, "", "土", "传统,信仰,传承", "遵循既有规则与传统智慧，寻求精神指引。", "墨守成规，挑战权威，教条主义。"},
	{6, "恋人", "The Lovers", 1, "", "风", "选择,结合,价值观", "真心相契的结合或价值观一致的重要抉择。", "关系失衡，错误选择，三心二意。"},
	{7, "战车", "The Chariot", 1, "", "水", "意志,胜利,掌控", "以坚定意志驾驭对立力量，迈向胜利。", "失去方向，冲动妄为，自我放纵。"},
	{8, "力量", "Strength", 1, "", "火", "勇气,柔韧,驯服", "以温柔与内在勇气驾驭原始力量。", "自我怀疑，情绪失控，软弱退缩。"},
	{9, "隐士", "The Hermit", 1, "", "土", "内省,智慧,独处", "独处沉思，向内寻找真正的答案。", "孤立封闭，拒绝建议，迷失方向。"},
	{10, "命运之轮", "Wheel of Fortune", 1, "", "火", "循环,机遇,转折", "命运之轮运转，迎接生命的重要转机。", "厄运循环，抗拒变化，错失机会。"},
	{11, "正义", "Justice", 1, "", "风", "公平,真相,责任", "因果平衡，真相浮现，承担应有责任。", "偏见不公，推卸责任，判断失误。"},
	{12, "倒吊人", "The Hanged Man", 1, "", "水", "暂停,换位,牺牲", "暂停等待，转换视角看见新的可能。", "拖延僵化，无谓牺牲，抗拒改变。"},
	{13, "死神", "Death", 1, "", "水", "结束,转变,新生", "彻底的终结带来真正的新生与蜕变。", "抗拒结束，停滞不前，恐惧改变。"},
	{14, "节制", "Temperance", 1, "", "火", "平衡,调和,耐心", "耐心调和对立，寻得和谐中道。", "失衡偏激，急躁冒进，格格不入。"},
	{15, "恶魔", "The Devil", 1, "", "土", "欲望,束缚,诱惑", "直面欲望与阴暗面，觉察心中枷锁。", "挣脱束缚，戒除瘾症，恢复自由。"},
	{16, "塔", "The Tower", 1, "", "火", "突变,崩坏,觉醒", "旧结构骤然崩塌，迎来觉醒真相。", "逃避变故，延迟崩塌，惊魂未定。"},
	{17, "星星", "The Star", 1, "", "风", "希望,灵感,疗愈", "风雨后的希望之光，心灵得到抚慰。", "失去信心，灵感枯竭，自我怀疑。"},
	{18, "月亮", "The Moon", 1, "", "水", "幻象,潜意识,迷惘", "潜意识浮现，直面内在恐惧与幻象。", "走出迷雾，释放恐惧，拨云见日。"},
	{19, "太阳", "The Sun", 1, "", "火", "喜悦,成功,活力", "阳光普照，尽情享受成功与生命的喜悦。", "暂时阴霾，过度乐观，延迟的喜悦。"},
	{20, "审判", "Judgement", 1, "", "火", "觉醒,重生,决断", "内在召唤响起，做出重要的人生决断。", "逃避审视，自我否定，错过召唤。"},
	{21, "世界", "The World", 1, "", "土", "完成,圆满,成就", "阶段圆满完成，迎来全新的整合与成就。", "尚未完成，半途而废，错失成果。"},
}

// 小阿卡纳花色定义
type suitDef struct {
	SuitEn  string // Wands/Cups/Swords/Pentacles
	SuitZh  string // 权杖/圣杯/宝剑/星币
	Element string // 火/水/风/土
}

var suits = []suitDef{
	{"Wands", "权杖", "火"},
	{"Cups", "圣杯", "水"},
	{"Swords", "宝剑", "风"},
	{"Pentacles", "星币", "土"},
}

// 每花色 14 张：Ace、2-10、Page、Knight、Queen、King
// 序号命名
var minorRanks = []struct {
	ZhLabel string
	EnLabel string
	// 正位关键字（各花色共通大方向，再拼花色语义）
}{
	{"一", "Ace"},
	{"二", "Two"},
	{"三", "Three"},
	{"四", "Four"},
	{"五", "Five"},
	{"六", "Six"},
	{"七", "Seven"},
	{"八", "Eight"},
	{"九", "Nine"},
	{"十", "Ten"},
	{"侍从", "Page"},
	{"骑士", "Knight"},
	{"王后", "Queen"},
	{"国王", "King"},
}

// 小阿卡纳寓意简表（14 张 × 4 花色 = 56 张）
// 为避免文件过长，寓意以各花色核心主题 + 阶段叙事生成
var suitThemes = map[string]struct {
	Core     string // 花色核心
	Keywords []string
}{
	"Wands":     {"激情与行动", []string{"起心动念,灵感", "合作,规划", "扩张,先见", "稳固,庆祝", "竞争,冲突", "胜利,喝彩", "坚守,挑战", "疾速,进展", "警觉,韧性", "重担,完成", "热情,探索", "冲劲,远行", "自信,魅力", "远见,领导"}},
	"Cups":      {"情感与直觉", []string{"情感涌现,爱", "情投意合,交流", "欢庆,友谊", "倦怠,内省", "失落,遗憾", "回忆,纯真", "幻想,抉择", "放下,追寻", "满足,愿望", "圆满,家庭", "浪漫,创意", "追求,理想", "温柔,共情", "沉稳,包容"}},
	"Swords":    {"思想与冲突", []string{"洞察,真相", "僵局,两难", "心痛,悲伤", "休养,沉思", "冲突,失和", "过渡,启程", "诡计,谨慎", "束缚,自限", "焦虑,忧心", "绝境,转机", "敏捷,好奇", "果敢,行动", "敏锐,理性", "权威,公正"}},
	"Pentacles": {"物质与务实", []string{"机会,种子", "平衡,兼顾", "协作,精进", "持守,稳固", "困顿,互助", "馈赠,分享", "评估,耐心", "专注,精益", "独立,丰盛", "传承,富足", "学习,踏实", "勤勉,守成", "滋养,安全", "丰饶,经营"}},
}

func buildMinorArcana() []tarotSeed {
	var out []tarotSeed
	cardNo := 22
	for _, s := range suits {
		theme := suitThemes[s.SuitEn]
		for i, r := range minorRanks {
			name := s.SuitZh + r.ZhLabel
			nameEn := r.EnLabel + " of " + s.SuitEn
			kw := theme.Keywords[i]
			upright := s.SuitZh + "(" + theme.Core + ") - " + kw + "，正位带来正向能量。"
			reversed := s.SuitZh + "(" + theme.Core + ") - " + kw + "，逆位提示阻滞或反面。"
			out = append(out, tarotSeed{
				CardNo:     cardNo,
				Name:       name,
				NameEn:     nameEn,
				ArcanaType: 2,
				Suit:       s.SuitEn,
				Element:    s.Element,
				Keywords:   kw,
				Upright:    upright,
				Reversed:   reversed,
			})
			cardNo++
		}
	}
	return out
}

func seedTarotCards(db *gorm.DB) {
	all := append([]tarotSeed{}, majorArcana...)
	all = append(all, buildMinorArcana()...)

	created := 0
	skipped := 0
	for _, c := range all {
		var cnt int64
		db.Model(&model.TarotCard{}).Where("card_no = ?", c.CardNo).Count(&cnt)
		if cnt > 0 {
			skipped++
			continue
		}
		nameEn := c.NameEn
		kw := c.Keywords
		elem := c.Element
		up := c.Upright
		rev := c.Reversed
		var suitPtr *string
		if c.Suit != "" {
			s := c.Suit
			suitPtr = &s
		}
		card := &model.TarotCard{
			CardNo:          c.CardNo,
			Name:            c.Name,
			NameEn:          &nameEn,
			ArcanaType:      c.ArcanaType,
			Suit:            suitPtr,
			ImageURL:        "", // 由管理员在后台上传
			UprightMeaning:  &up,
			ReversedMeaning: &rev,
			Keywords:        &kw,
			Element:         &elem,
		}
		if err := db.Create(card).Error; err != nil {
			log.Printf("  - 创建塔罗牌失败 card_no=%d: %v", c.CardNo, err)
			continue
		}
		created++
	}
	log.Printf("  - 塔罗牌: 总 %d 张（新增 %d，已存在 %d）", len(all), created, skipped)
}

// migrateTarotMappingUniqueIndex 将 tarot_drink_mappings 的唯一键从 (card_no,drink_id,is_reversed)
// 改为业务正确的 (card_no,is_reversed)。幂等：既存索引名匹配则直接跳过。
func migrateTarotMappingUniqueIndex(db *gorm.DB) {
	// 清理可能存在的重复映射（同 card_no+is_reversed 有多条），保留 match_score 最高 / id 最小的一条
	if err := db.Exec(`
		DELETE FROM tarot_drink_mappings a
		USING tarot_drink_mappings b
		WHERE a.card_no = b.card_no
		  AND a.is_reversed = b.is_reversed
		  AND (a.match_score < b.match_score
		       OR (a.match_score = b.match_score AND a.id > b.id))
	`).Error; err != nil {
		log.Printf("  - 清理重复 tarot_drink_mappings 失败: %v", err)
	}
	// 删除旧唯一索引
	if err := db.Exec(`DROP INDEX IF EXISTS uk_card_drink_reversed`).Error; err != nil {
		log.Printf("  - DROP 旧索引失败: %v", err)
	}
	// 创建新唯一索引（GORM AutoMigrate 对已存在表的 tag 新索引不一定会创建，显式执行更可靠）
	if err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS uk_card_reversed ON tarot_drink_mappings(card_no, is_reversed)`).Error; err != nil {
		log.Printf("  - CREATE 新唯一索引失败: %v", err)
		return
	}
	log.Println("  - tarot_drink_mappings 唯一索引已迁移为 (card_no, is_reversed)")
}

// seedTarotDrinkMappings 为 78 张牌 × 2 朝向 = 156 条默认映射提供兑底，
// 使塔罗占卜功能在管理员手动配置前即可工作。
// 选题规则：取上架的首张 drink（id 最小）作为兑底；仅在对应映射不存在时插入。
func seedTarotDrinkMappings(db *gorm.DB) {
	var defaultDrink model.Drink
	if err := db.Where("status = ?", 1).Order("id ASC").First(&defaultDrink).Error; err != nil {
		log.Printf("  - 跳过塔罗映射种子：没有上架酒水可供兑底 (err=%v)", err)
		return
	}

	created := 0
	skipped := 0
	for cardNo := 0; cardNo <= 77; cardNo++ {
		for _, rev := range []int8{0, 1} {
			var cnt int64
			db.Model(&model.TarotDrinkMapping{}).
				Where("card_no = ? AND is_reversed = ?", cardNo, rev).
				Count(&cnt)
			if cnt > 0 {
				skipped++
				continue
			}
			reason := "默认推荐，管理员后续在后台调优"
			m := &model.TarotDrinkMapping{
				CardNo:         cardNo,
				DrinkID:        defaultDrink.ID,
				IsReversed:     rev,
				MatchScore:     50,
				ReasonTemplate: &reason,
			}
			if err := db.Create(m).Error; err != nil {
				log.Printf("  - 创建塔罗映射失败 card_no=%d rev=%d: %v", cardNo, rev, err)
				continue
			}
			created++
		}
	}
	log.Printf("  - 塔罗-酒水映射: 总 156 条（新增 %d，已存在 %d），兑底 drink_id=%d", created, skipped, defaultDrink.ID)
}
