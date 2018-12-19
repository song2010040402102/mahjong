package logic

import (
	"config"
	"github.com/astaxie/beego/logs"
	"math"
	"sort"
	"util"
)

const ai_fact_num int32 = 10       //数量因子
const ai_fact_dis int32 = 10       //距离因子
const ai_fact_pos int32 = 1        //位置因子
const ai_rel_edg_val float32 = 5.0 //边界值，小于边界值的牌可认为没有关系的牌

type CardValue struct {
	Card int32
	Val  float32
}

func ReplaceBaiWithLZ(chiCard []*ChiCard, aiCards []AICard, lzCards []int32) {
	if len(lzCards) == 0 {
		return
	}
	for _, v := range chiCard {
		if v.CardId == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI && v.CardType == MJ_CHI_CHI {
			v.CardId = lzCards[0]
		}
	}

	for i := 0; i < len(aiCards); i++ {
		if aiCards[i].Card == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
			aiCards[i].Card = lzCards[0]
			sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })
			break
		}
	}
}

func ReplaceLZWithBai(chiCard []*ChiCard, aiCards []AICard, lzCards []int32) {
	if len(lzCards) == 0 {
		return
	}
	for _, v := range chiCard {
		if v.CardId == lzCards[0] && v.CardType == MJ_CHI_CHI {
			v.CardId = COLOR_OTHER*MAHJONG_MASK + MAHJONG_BAI
		}
	}
	for i := 0; i < len(aiCards); i++ {
		if aiCards[i].Card == lzCards[0] {
			aiCards[i].Card = COLOR_OTHER*MAHJONG_MASK + MAHJONG_BAI
			sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })
			break
		}
	}
}

func ConvertSliceToAICard(holdCards []int32, card int32, lzCards []int32) ([]AICard, int32) {
	cards := make([]int32, len(holdCards))
	copy(cards, holdCards)
	if card > 0 {
		cards = append(cards, card)
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })

	lzNum := int32(0)
	aiCards := make([]AICard, 0, len(cards))
	for i, j := 0, 0; i < len(cards); j++ {
		aiCards = append(aiCards, AICard{cards[i], 1})
		for {
			i++
			if i < len(cards) && cards[i-1] == cards[i] {
				aiCards[j].Num++
			} else {
				break
			}
		}
		isLZ := false
		if aiCards[j].Card == 0 && len(lzCards) == 0 {
			isLZ = true
		} else {
			for _, v := range lzCards {
				if aiCards[j].Card == v {
					isLZ = true
					break
				}
			}
		}
		if isLZ {
			lzNum += aiCards[j].Num
			aiCards = aiCards[:j]
			j--
		}
	}
	return aiCards, lzNum
}

//获取牌价值
func GetCardValue(rule IMahjong, holdCards []int32, moCard int32) []CardValue {
	aiCards, _ := ConvertSliceToAICard(holdCards, moCard, rule.GetLaiziCard())
	if rule.IsCanBaiReplaceLZ() == true {
		ReplaceBaiWithLZ([]*ChiCard{}, aiCards, rule.GetLaiziCard())
	}
	valCards := analyzeCardValue(aiCards)
	if rule.IsCanBaiReplaceLZ() == true {
		for i := 0; i < len(valCards); i++ {
			lzCards := rule.GetLaiziCard()
			if len(lzCards) > 0 && valCards[i].Card == lzCards[0] {
				valCards[i].Card = COLOR_OTHER*MAHJONG_MASK + MAHJONG_BAI
				break
			}
		}
	}
	return valCards
}

func analyzeCardValue(aiCards []AICard) []CardValue {
	//对每张牌进行价值评估，由位置、距离、数量三个指标构成
	valCards := make([]CardValue, len(aiCards))
	for i := 0; i < len(aiCards); i++ {
		val := float32(0)
		if aiCards[i].Num > 1 {
			nv := float32(ai_fact_num * (aiCards[i].Num - 1)) //数量评估值
			val += nv
		}
		if c := aiCards[i].Card % MAHJONG_MASK; c >= MAHJONG_1 && c <= MAHJONG_9 {
			if i < len(aiCards)-1 {
				d := aiCards[i+1].Card - aiCards[i].Card
				if d == 1 && (aiCards[i].Card%MAHJONG_MASK == MAHJONG_1 || aiCards[i+1].Card%MAHJONG_MASK == MAHJONG_9) {
					d = 2
				}
				if d > 0 && d <= MAHJONG_9-MAHJONG_1 {
					if d > 2 {
						d *= 2 //这里可避免22588牌型和12牌型先打1而不打5
					}
					dv := float32(ai_fact_dis) / float32(d) //距离评估值
					val += dv                               //右距离值
					valCards[i+1].Val = dv                  //左距离值
				}
			}
			pv := float32(float64(ai_fact_pos) / (math.Abs(float64(aiCards[i].Card%MAHJONG_MASK-MAHJONG_5)) + 1)) //位置评估值
			val += pv
		}
		valCards[i].Card = aiCards[i].Card
		valCards[i].Val += val
	}
	sort.Slice(valCards, func(i, j int) bool { return valCards[i].Val <= valCards[j].Val })
	return valCards
}

func removeSolKeSeq(rule IMahjong, cards []int32) []int32 {
	lzCards := rule.GetLaiziCard()
	if len(lzCards) > 0 {
		for i := 0; i < len(cards); i++ {
			if rule.IsLaiziCard(cards[i]) == true {
				cards[i] = MAHJONG_LZ
			} else if rule.IsCanBaiReplaceLZ() == true && cards[i] == MAHJONG_MASK*COLOR_OTHER+MAHJONG_BAI {
				cards[i] = lzCards[0]
			}
		}
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })
	for i := 0; i < len(cards)-2; {
		if cards[i] == MAHJONG_LZ || cards[i+1] == MAHJONG_LZ || cards[i+2] == MAHJONG_LZ {
			i++
			continue
		}
		sol := (i == 0 || cards[i-1] == MAHJONG_LZ || cards[i-1] < cards[i]-2) && (i == len(cards)-3 || cards[i+3] == MAHJONG_LZ || cards[i+3] > cards[i+2]+2)
		if cards[i] == cards[i+1] && cards[i+1] == cards[i+2] && (cards[i]%MAHJONG_MASK >= MAHJONG_DONG || sol) {
			cards = append(cards[:i], cards[i+3:]...) //移除独立的刻子
		} else if c := cards[i] % MAHJONG_MASK; c >= MAHJONG_1 && c <= MAHJONG_7 && cards[i] == cards[i+1]-1 && cards[i+1] == cards[i+2]-1 && sol {
			cards = append(cards[:i], cards[i+3:]...) //移除独立的顺子
		} else {
			i++
		}
	}
	if len(lzCards) > 0 {
		for i := 0; i < len(cards); i++ {
			if cards[i] == MAHJONG_LZ {
				cards[i] = lzCards[0]
			} else if rule.IsCanBaiReplaceLZ() == true && rule.IsLaiziCard(cards[i]) == true {
				cards[i] = MAHJONG_MASK*COLOR_OTHER + MAHJONG_BAI
			}
		}
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })
	return cards
}

type IAIMj interface {
	//设置规则
	SetRule(rule IMahjong)
	//设置所有牌
	SetAllCard(holdCards map[int32][]int32, chiCards map[int32][]*ChiCard, usedCards map[int32][]int32, restCards []int32)
	//设置级别
	SetLevel(level int32)
	//获取发牌，可控制机器人起手牌
	GetDealCards() []int32
	//获取下一张牌，可控制机器人摸牌
	GetNextCards(index int32) []int32
	//托管出牌
	GetCardForMandate(index int32, moCard int32) int32
	//机器人出牌
	GetCardForRobot(index int32, moCard int32) int32
	//托管吃牌
	GetChiForMandate(index int32, moCard int32, chiCards []*ChiCard) *ChiCard
	//机器人吃牌
	GetChiForRobot(index int32, moCard int32, chiCards []*ChiCard) *ChiCard
}

type AIMjBase struct {
	rule      IMahjong
	holdCards map[int32][]int32
	chiCards  map[int32][]*ChiCard
	usedCards map[int32][]int32
	restCards []int32
	level     int32
}

func (ai *AIMjBase) SetRule(rule IMahjong) {
	ai.rule = rule
}

func (ai *AIMjBase) SetAllCard(holdCards map[int32][]int32, chiCards map[int32][]*ChiCard, usedCards map[int32][]int32, restCards []int32) {
	ai.holdCards = holdCards
	ai.chiCards = chiCards
	ai.usedCards = usedCards
	ai.restCards = restCards
}

func (ai *AIMjBase) SetLevel(level int32) {
	ai.level = level	
	logs.Info("AI SetLevel, level: ", level)
}

func (ai *AIMjBase) GetDealCards() []int32 {
	if ai.level == ROBOT_LEVEL_MASTER {
		cards := []int32{}
		if util.GetRandomRate() < 0.7 { //70%概率混一色加箭刻
			jian := MAHJONG_MASK*COLOR_OTHER + MAHJONG_HONGZHONG + util.GetRandom(0, 2)
			cards = append(cards, []int32{jian, jian, jian}...)
		}

		colors := []int32{}
		if ai.rule.GetDeck().GetFlag()&WITH_WAN != 0 {
			colors = append(colors, COLOR_WAN)
		}
		if ai.rule.GetDeck().GetFlag()&WITH_TONG != 0 {
			colors = append(colors, COLOR_TONG)
		}
		if ai.rule.GetDeck().GetFlag()&WITH_TIAO != 0 {
			colors = append(colors, COLOR_TIAO)
		}
		color := colors[util.GetRandom(0, int32(len(colors)-1))]

		mc := make(map[int32]int32)
		cardNum := ai.rule.GetPlayerCardNum()
		for {
			c := color*MAHJONG_MASK + util.GetRandom(MAHJONG_1, MAHJONG_9)
			if mc[c] >= 4 {
				continue
			}
			mc[c]++
			cards = append(cards, c)
			if int32(len(cards)) >= cardNum {
				break
			}
		}
		logs.Info("[AIMjBase]GetDealCards, cards: ", cards)
		return cards
	}
	return []int32{}
}

func (ai *AIMjBase) GetNextCards(index int32) []int32 {
	if ai.level == ROBOT_LEVEL_NOOB {
		_, tings := ai.rule.CheckNTing(ai.chiCards[index], ai.holdCards[index], 1, false)
		if len(tings) > 0 {
			return ai.getUnusedCards(index)
		}
	} else if ai.level == ROBOT_LEVEL_MASTER {
		if util.GetRandomRate() < 0.5 { //一半的概率能摸到想要的牌
			return ai.getNeedCards(index)
		}
	}
	return []int32{}
}

func (ai *AIMjBase) GetCardForMandate(index int32, moCard int32) int32 {
	if *config.FLAG_DEBUG == 1 || moCard == 0 || ai.rule.IsLaiziCard(moCard) == true {
		return ai.getCardForWeight(ai.holdCards[index], moCard)
	}
	return moCard
}

func (ai *AIMjBase) GetCardForRobot(index int32, moCard int32) int32 {
	if ai.level == ROBOT_LEVEL_NOOB || ai.level == ROBOT_LEVEL_AMATEUR {
		return ai.getCardForWeight(ai.holdCards[index], moCard)
	} else if ai.level == ROBOT_LEVEL_MAJOR || ai.level == ROBOT_LEVEL_MASTER {
		return ai.getCardForTingNum(ai.holdCards[index], moCard, index, ai.chiCards[index])
	} else {
		logs.Error("GetCardForRobot, level invalid!")
	}
	return 0
}

func (ai *AIMjBase) GetChiForMandate(index int32, moCard int32, chiCards []*ChiCard) *ChiCard {
	for _, v := range chiCards {
		if v != nil && v.CardType == MJ_CHI_HU {
			return v
		}
	}
	if len(chiCards) > 0 && chiCards[0] != nil {
		chiCard := NewChiCard()
		*chiCard = *chiCards[0]
		chiCard.CardType = MJ_CHI_PASS
		return chiCard
	}
	return nil
}

func (ai *AIMjBase) GetChiForRobot(index int32, moCard int32, chiCards []*ChiCard) *ChiCard {
	if len(chiCards) == 0 || chiCards[0] == nil {
		return nil
	}
	chiCard := NewChiCard()
	*chiCard = *chiCards[0]
	chiCard.ChiPosBit = 0
	if ai.level == ROBOT_LEVEL_NOOB { //菜鸟级机器人在一上一听阶段截住吃和碰
		if chiCard.CardType == MJ_CHI_PENG || chiCard.CardType == MJ_CHI_CHI {
			_, tings := ai.rule.CheckNTing(ai.chiCards[index], ai.holdCards[index], 1, false)
			if len(tings) > 0 {
				chiCard.CardType = MJ_CHI_PASS
			}
		}
	} else if ai.level == ROBOT_LEVEL_AMATEUR { //业余级机器人按胡>杠>碰>吃优先级

	} else if ai.level == ROBOT_LEVEL_MAJOR || ai.level == ROBOT_LEVEL_MASTER { //专家级和大师级向听分析胡杠碰吃过
		/* 对于杠碰吃过的选择应该由向听数是否减少来决定
		** 胡牌和独立杠基本不依赖特定规则，所以可以直接处理
		** 碰和吃的独立依赖于特定规则，例如碰可能会破坏将或七对，吃可能会破坏全不靠
		 */
		if chiCard.CardType == MJ_CHI_HU {
			return chiCard
		}
		alGang := ai.getAloneGang(index, moCard, chiCards)
		if alGang != nil {
			return alGang
		}

		newChiCards := make([]*ChiCard, len(chiCards))
		copy(newChiCards, chiCards)
		if chiCards[len(chiCards)-1].CardType == MJ_CHI_CHI {
			newChiCards = newChiCards[:len(newChiCards)-1]
			for i := uint(1); i <= uint(3); i++ {
				if chiCards[len(chiCards)-1].ChiPosBit&(0x01<<i) != 0 {
					chi := NewChiCard()
					*chi = *chiCards[len(chiCards)-1]
					chi.ChiPosBit = int32(0x01 << i)
					newChiCards = append(newChiCards, chi)
				}
			}
		}

		chiCard.CardType = MJ_CHI_PASS
		preCards := make([]int32, len(ai.holdCards[index]))
		copy(preCards, ai.holdCards[index])

		if moCard > 0 {
			preCards = append(preCards, moCard)
			playC := ai.getCardForTingNum(preCards, 0, index, ai.chiCards[index])
			preCards = util.RemoveSliceElem(preCards, playC, false).([]int32)
		}
		preN := ai.getMainCardTingNum(ai.chiCards[index], preCards)
		logs.Info("GetChiForRobot, pass, preN: ", preN)

		for _, v := range newChiCards {
			curCards := make([]int32, len(ai.holdCards[index]))
			curChis := make([]*ChiCard, len(ai.chiCards[index]))
			copy(curCards, ai.holdCards[index])
			copy(curChis, ai.chiCards[index])
			if moCard > 0 {
				curCards = append(curCards, moCard)
			}
			curChis = append(curChis, v)
			if v.CardType == MJ_CHI_GANG_AN || v.CardType == MJ_CHI_GANG {
				curCards = util.RemoveSliceElem(curCards, v.CardId, false).([]int32)
				curCards = util.RemoveSliceElem(curCards, v.CardId, false).([]int32)
				curCards = util.RemoveSliceElem(curCards, v.CardId, false).([]int32)
				if v.CardType == MJ_CHI_GANG_AN {
					curCards = util.RemoveSliceElem(curCards, v.CardId, false).([]int32)
				}
			} else if v.CardType == MJ_CHI_GANG_WAN {
				for i := 0; i < len(curChis); i++ {
					if curChis[i].CardType == MJ_CHI_PENG && curChis[i].CardId == v.CardId {
						curChis = append(curChis[:i], curChis[i+1:]...)
						break
					}
				}
				curCards = util.RemoveSliceElem(curCards, v.CardId, false).([]int32)
			} else if v.CardType == MJ_CHI_PENG {
				curCards = util.RemoveSliceElem(curCards, v.CardId, false).([]int32)
				curCards = util.RemoveSliceElem(curCards, v.CardId, false).([]int32)
			} else if v.CardType == MJ_CHI_CHI {
				chiC := ai.rule.GetChiCard(v)
				for _, vv := range chiC {
					if vv != v.CardId {
						curCards = util.RemoveSliceElem(curCards, vv, false).([]int32)
					}
				}
			}
			if v.CardType == MJ_CHI_PENG || v.CardType == MJ_CHI_CHI {
				playC := ai.getCardForTingNum(curCards, 0, index, curChis)
				curCards = util.RemoveSliceElem(curCards, playC, false).([]int32)
			}
			curN := ai.getMainCardTingNum(curChis, curCards)
			if curN < preN || curN == preN && v.CardType >= MJ_CHI_GANG && v.CardType <= MJ_CHI_GANG_AN {
				//杠不会进听，所以如果没有退听，可以优先杠，对于吃和碰，可能退换进听，对于换听，理想情况下应该考虑换听后的听牌总数，但意义不大会有损效率
				preN = curN
				*chiCard = *v
			}
			logs.Info("GetChiForRobot, type: ", v.CardType, " curN: ", curN)
		}
	} else {
		logs.Error("GetChiForRobot, level invalid!")
		return nil
	}
	return chiCard
}

//随机出牌
func (ai *AIMjBase) getCardForRand(holdCards []int32, moCard int32) int32 {
	aiCards, _ := ConvertSliceToAICard(holdCards, moCard, ai.rule.GetLaiziCard())
	if len(aiCards) == 0 {
		logs.Error("getCardForRand, aiCards empty!")
		return 0
	}
	index := util.GetRandom(0, int32(len(aiCards)-1))
	return aiCards[index].Card
}

//权重出牌
func (ai *AIMjBase) getCardForWeight(holdCards []int32, moCard int32) int32 {
	valCards := GetCardValue(ai.rule, holdCards, moCard)
	if len(valCards) == 0 {
		logs.Error("getCardForWeight, valCards empty!")
		return 0
	}
	return valCards[0].Card
}

//向听数出牌
func (ai *AIMjBase) getCardForTingNum(holdCards []int32, moCard int32, index int32, chiCards []*ChiCard) int32 {
	cards := make([]int32, len(holdCards))
	copy(cards, holdCards)
	if moCard > 0 {
		cards = append(cards, moCard)
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })
	if len(cards) > 3 {
		cards = removeSolKeSeq(ai.rule, cards)
	}

	valCards := GetCardValue(ai.rule, cards, 0)
	if len(valCards) == 0 {
		logs.Error("getCardForTingNum, valCards empty!")
		return 0
	}
	logs.Info("getCardForTingNum, cards: ", cards, " valCards: ", valCards)

	//直接打出唯一牌
	if len(valCards) == 1 {
		return valCards[0].Card
	}

	//在牌数至少为5的情况下，直接打出最差的无关联牌
	if len(cards) > 2 && valCards[0].Val <= ai_rel_edg_val {
		return valCards[0].Card
	}

	//若仅剩两张不同的非癞子牌时，直接选择剩余多的那张牌作为胡牌
	if len(cards) == 2 {
		num0 := ai.getCardRemainNum(cards, chiCards, index, cards[0])
		num1 := ai.getCardRemainNum(cards, chiCards, index, cards[1])
		if num0 >= num1 {
			return cards[1]
		} else {
			return cards[0]
		}
	}

	//获取癞子牌数量
	lzNum := int32(0)
	for _, v := range cards {
		if ai.rule.IsLaiziCard(v) == true {
			lzNum++
		}
	}

	//尝试的最大向听数，没有单牌的情况下，可减去癞子牌数量
	maxN := int32(len(cards)/3) - lzNum
	if maxN < 0 {
		maxN = 0
	}

	//按权重依次尝试
	maxTingC, minTingN, maxTingV := int32(0), maxN, int32(0)
	bkCards := make([]int32, len(cards))
	copy(bkCards, cards)
	for _, vc := range valCards {
		cards = util.RemoveSliceElem(cards, vc.Card, false).([]int32)
		realN, tings := ai.rule.CheckNTing(chiCards, cards, minTingN, true)
		if len(tings) > 0 {
			logs.Info("getCardForTingNum, cards: ", cards, " c: ", vc.Card, " realN: ", realN, " tings: ", tings)
			tingV := int32(0)
			for _, v := range tings {
				tingV += ai.getCardRemainNum(cards, chiCards, index, v)
			}
			if realN < minTingN || tingV > maxTingV {
				maxTingC, minTingN, maxTingV = vc.Card, realN, tingV
			}
		}
		cards = make([]int32, len(bkCards))
		copy(cards, bkCards)
	}
	if maxTingC == 0 {
		maxTingC = valCards[0].Card
	}
	logs.Info("getCardForTingNum, maxTingC: ", maxTingC, " minTingN: ", minTingN, " maxTingV: ", maxTingV)
	return maxTingC
}

//获取独立的杠
func (ai *AIMjBase) getAloneGang(index int32, moCard int32, chiCards []*ChiCard) *ChiCard {
	if len(chiCards) > 0 && chiCards[0] != nil &&
		chiCards[0].CardType == MJ_CHI_GANG_AN || chiCards[0].CardType == MJ_CHI_GANG_WAN || chiCards[0].CardType == MJ_CHI_GANG {
		cards := make([]int32, len(ai.holdCards[index]))
		copy(cards, ai.holdCards[index])
		sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })
		for _, v := range chiCards {
			card := v.CardId
			if ai.rule.IsCanBaiReplaceLZ() == true && v.CardId == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
				lzCards := ai.rule.GetLaiziCard()
				if len(lzCards) > 0 {
					card = lzCards[0]
				}
			}
			if c := card % MAHJONG_MASK; c >= MAHJONG_DONG && c <= MAHJONG_BAI {
				return v
			}
			alone := true
			for _, c := range cards {
				if c == card-1 || c == card-2 || c == card+1 || c == card+2 {
					alone = false
					break
				}
			}
			if alone {
				return v
			}
		}
	}
	return nil
}

//移除单张牌
func (ai *AIMjBase) removeSingleCard(chiCards []*ChiCard, cards []int32) []int32 {
	huaCards := []int32{}
	for _, v := range cards {
		if ai.rule.IsHuaCard(v) {
			huaCards = append(huaCards, v)
		}
	}
	for _, v := range huaCards {
		cards = util.RemoveSliceElem(cards, v, false).([]int32)
	}
	singleN := 0
	valCards := GetCardValue(ai.rule, cards, 0)
	for _, v := range valCards {
		if v.Val <= ai_rel_edg_val {
			singleN++
		} else {
			break
		}
	}
	if singleN > 3 {
		for i := 0; i < singleN-singleN%3; i++ {
			cards = util.RemoveSliceElem(cards, valCards[i].Card, false).([]int32)
		}
	}
	return cards
}

//获取主牌的向听数
func (ai *AIMjBase) getMainCardTingNum(chiCards []*ChiCard, holdCards []int32) int32 {
	cards := make([]int32, len(holdCards))
	copy(cards, holdCards)
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })
	cards = ai.removeSingleCard(chiCards, cards)
	if len(cards) > 3 {
		cards = removeSolKeSeq(ai.rule, cards)
	}

	//尝试的最大向听数，不减去癞子牌数量，因为会存在单牌
	maxN := int32(len(cards) / 3)
	if maxN < 0 {
		maxN = 0
	}

	realN, _ := ai.rule.CheckNTing(chiCards, cards, maxN, false)
	if realN < 0 {
		realN = int32(len(cards) * 2)
		logs.Error("getMainCardTingNum, maxN error!")
	}
	logs.Info("getMainCardTingNum, cards: ", cards, " maxN: ", maxN, " realN: ", realN)
	return realN
}

/* 获取可能的剩余牌数
** 这个函数主要用于计算向听总数，严格意义应该考虑碰权重是吃权重的两倍，在向听搜索中如果考虑碰或吃或者既碰又吃，会极其复杂，若需要再考虑
 */
func (ai *AIMjBase) getCardRemainNum(cards []int32, chiCards []*ChiCard, index int32, card int32) int32 {
	rNum := int32(4)
	if card == MAHJONG_ANY {
		return 10000
	} else if ai.rule.IsLaiziCard(card) == true { //听牌中必然含癞子牌，所以可统一返回，随后再考虑翻牌
		return rNum
	}
	if ai.level == ROBOT_LEVEL_MAJOR || ai.level == ROBOT_LEVEL_MASTER {
		for _, v := range cards {
			if v == card {
				rNum--
				if rNum <= 0 {
					return 0
				}
			}
		}
		for k, v := range ai.chiCards {
			chis := v
			if k == index {
				chis = chiCards
			}
			for _, vv := range chis {
				chiC := ai.rule.GetChiCard(vv)
				for _, vvv := range chiC {
					if vvv == card {
						rNum--
						if rNum <= 0 {
							return 0
						}
					}
				}
			}
		}
		for _, v := range ai.usedCards {
			for _, vv := range v {
				if vv == card {
					rNum--
					if rNum <= 0 {
						return 0
					}
				}
			}
		}
	}
	return rNum
}

//获取玩家不需要的牌
func (ai *AIMjBase) getUnusedCards(index int32) []int32 {
	unusedC := []int32{}
	lzCards := ai.rule.GetLaiziCard()
	aiCards, _ := ConvertSliceToAICard(ai.holdCards[index], 0, lzCards)
	if ai.rule.IsCanBaiReplaceLZ() == true {
		if len(lzCards) > 0 {
			for _, v := range aiCards {
				if v.Card == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
					aiCards = append(aiCards, AICard{lzCards[0], 1})
					sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })
					break
				}
			}
		}
	}
	lackColor := []int32{COLOR_WAN, COLOR_TONG, COLOR_TIAO, COLOR_OTHER}
	for i := 0; i < len(aiCards); i++ {
		c := aiCards[i].Card % MAHJONG_MASK
		lackColor = util.RemoveSliceElem(lackColor, aiCards[i].Card/MAHJONG_MASK, false).([]int32)
		if i == 0 || aiCards[i].Card-aiCards[i-1].Card > MAHJONG_9-MAHJONG_1 {
			start, end := MAHJONG_1, c-1
			if c >= MAHJONG_DONG && c <= MAHJONG_BAI {
				start, end = MAHJONG_DONG, c
			}
			for j := start; j < end; j++ {
				unusedC = append(unusedC, aiCards[i].Card-c+j)
			}
		}
		if i < len(aiCards)-1 && aiCards[i+1].Card-aiCards[i].Card <= MAHJONG_9-MAHJONG_1 {
			offset := int32(2)
			if c >= MAHJONG_DONG && c <= MAHJONG_BAI {
				offset = int32(1)
			}
			for j := aiCards[i].Card + offset; j <= aiCards[i+1].Card-offset; j++ {
				unusedC = append(unusedC, j)
			}
		}
		if i == len(aiCards)-1 || aiCards[i+1].Card-aiCards[i].Card > MAHJONG_9-MAHJONG_1 {
			start, end := c+2, MAHJONG_9
			if c >= MAHJONG_DONG && c <= MAHJONG_BAI {
				start, end = c+1, MAHJONG_BAI
			}
			for j := start; j <= end; j++ {
				unusedC = append(unusedC, aiCards[i].Card-c+j)
			}
		}
	}
	for _, v := range lackColor {
		if v == COLOR_OTHER {
			for c := MAHJONG_DONG; c <= MAHJONG_BAI; c++ {
				unusedC = append(unusedC, COLOR_OTHER*MAHJONG_MASK+c)
			}
		} else {
			for c := MAHJONG_1; c <= MAHJONG_9; c++ {
				unusedC = append(unusedC, v*MAHJONG_MASK+c)
			}
		}
	}
	for _, v := range lzCards {
		unusedC = util.RemoveSliceElem(unusedC, v, false).([]int32)
	}
	sort.Slice(unusedC, func(i, j int) bool { return unusedC[i] < unusedC[j] })
	return unusedC
}

//获取玩家需要的牌
func (ai *AIMjBase) getNeedCards(index int32) []int32 {
	needC := []int32{}
	cards := make([]int32, len(ai.holdCards[index]))
	copy(cards, ai.holdCards[index])
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })
	if len(cards) > 3 {
		cards = removeSolKeSeq(ai.rule, cards)
	}

	aiCards, lzNum := ConvertSliceToAICard(cards, 0, ai.rule.GetLaiziCard())
	maxN := int32(len(cards)/3) - lzNum
	if maxN < 0 {
		maxN = 0
	}
	_, tings := ai.rule.CheckNTing(ai.chiCards[index], cards, maxN, true)
	if len(tings) > 0 {
		needC = append(needC, tings...)
	}
	if len(needC) == 0 {
		for _, v := range aiCards {
			if ai.rule.IsLaiziCard(v.Card) == true {
				continue
			}
			needC = append(needC, v.Card)
		}
	}
	return needC
}
