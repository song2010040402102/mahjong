package mahjong

import (
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
	valCards := analyzeCardValue(aiCards)
	lzCards := []int32{}
	for _, v := range holdCards {
		if rule.IsLaiziCard(v) {
			lzCards = append(lzCards, v)
		}
	}
	if rule.IsLaiziCard(moCard) {
		lzCards = append(lzCards, moCard)
	}
	lzCards = util.UniqueSlice(lzCards, false).([]int32)
	for _, v := range lzCards {
		valCards = append(valCards, CardValue{v, 1e9})
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

type IAIMj interface {
	//设置规则
	SetRule(rule IMahjong)
	//设置所有牌
	SetAllCard(holdCards map[int32][]int32, chiCards map[int32][]*ChiCard, usedCards map[int32][]int32, restCards []int32)
	//设置级别
	SetLevel(level int32)
	//获取发牌，可控制机器人起手牌
	GetDealCards(cards []int32, cardNum int32, level int32, retry bool) int32
	//获取下一张牌，可控制机器人摸牌
	GetNextCards(index int32, level int32, robot bool) []int32
	//托管出牌
	GetCardForMandate(index int32, moCard int32, cards []int32) int32
	//机器人出牌
	GetCardForRobot(index int32, moCard int32, cards []int32) int32
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
}

func (ai *AIMjBase) GetCardForMandate(index int32, moCard int32, cards []int32) int32 {
	moAllow := false
	for _, v := range cards {
		if v == moCard {
			moAllow = true
			break
		}
	}
	if moCard == 0 || !moAllow || ai.rule.IsLaiziCard(moCard) {
		return ai.getCardForWeight(ai.holdCards[index], moCard, cards)
	}
	return moCard
}

func (ai *AIMjBase) GetCardForRobot(index int32, moCard int32, cards []int32) int32 {
	if ai.level == ROBOT_LEVEL_AMATEUR {
		return ai.getCardForWeight(ai.holdCards[index], moCard, cards)
	} else if ai.level == ROBOT_LEVEL_MAJOR {
		return ai.getCardForTingNum(ai.holdCards[index], moCard, cards, index, ai.chiCards[index])
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
	chiCard.SelectFirstChi()
	if ai.level == ROBOT_LEVEL_AMATEUR { //业余级机器人按胡>杠>碰>吃优先级

	} else if ai.level == ROBOT_LEVEL_MAJOR { //专家级
		/* 对于杠碰吃过的选择应该由向听数是否减少来决定
		** 胡牌和独立杠基本不依赖特定规则，所以可以直接处理
		** 碰和吃的独立依赖于特定规则，例如碰可能会破坏将或七对
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
			playC := ai.getCardForTingNum(preCards, 0, preCards, index, ai.chiCards[index])
			preCards = util.RemoveSliceElem(preCards, playC, false).([]int32)
		}
		preN, delN := ai.getMainCardTingNum(ai.chiCards[index], preCards, 0)

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
			playC := int32(0)
			if v.CardType == MJ_CHI_PENG || v.CardType == MJ_CHI_CHI {
				playC = ai.getCardForTingNum(curCards, 0, curCards, index, curChis)
				curCards = util.RemoveSliceElem(curCards, playC, false).([]int32)
			}
			if playC != v.CardId { //避免吃碰的牌和打牌相等
				curN, _ := ai.getMainCardTingNum(curChis, curCards, delN)
				if curN < preN || curN == preN && v.CardType >= MJ_CHI_GANG && v.CardType <= MJ_CHI_GANG_AN {
					//杠不会进听，所以如果没有退听，可以优先杠，对于吃和碰，可能退换进听，对于换听，理想情况下应该考虑换听后的听牌总数，但意义不大会有损效率
					preN = curN
					*chiCard = *v
				}
			}
		}
	} else {
		return nil
	}
	return chiCard
}

//权重出牌
func (ai *AIMjBase) getCardForWeight(holdCards []int32, moCard int32, allowCards []int32) int32 {
	valCards := GetCardValue(ai.rule, holdCards, moCard)
	for _, v := range valCards {
		for _, vv := range allowCards {
			if v.Card == vv {
				return v.Card
			}
		}
	}
	return 0
}

//向听数出牌
func (ai *AIMjBase) getCardForTingNum(holdCards []int32, moCard int32, allowCards []int32, index int32, chiCards []*ChiCard) int32 {
	cards := make([]int32, len(holdCards))
	copy(cards, holdCards)
	if moCard > 0 {
		cards = append(cards, moCard)
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })

	valCards := GetCardValue(ai.rule, cards, 0)
	for i := 0; i < len(valCards); i++ {
		exist := false
		for _, v := range allowCards {
			if valCards[i].Card == v {
				exist = true
				break
			}
		}
		if !exist {
			valCards = append(valCards[:i], valCards[i+1:]...)
			i--
		}
	}
	if len(valCards) == 0 {
		return 0
	}

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

//获取主牌的向听数
func (ai *AIMjBase) getMainCardTingNum(chiCards []*ChiCard, holdCards []int32, delN int32) (int32, int32) {
	cards := make([]int32, len(holdCards))
	copy(cards, holdCards)
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })

	valCards := GetCardValue(ai.rule, cards, 0)
	if delN == 0 {
		singleN := 0
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
				delN++
			}
		}
	} else {
		for i := int32(0); i < delN; i++ {
			cards = util.RemoveSliceElem(cards, valCards[i].Card, false).([]int32)
		}
	}

	//尝试的最大向听数，不减去癞子牌数量，因为会存在单牌
	maxN := int32(len(cards) / 3)
	if maxN < 0 {
		maxN = 0
	}

	realN, _ := ai.rule.CheckNTing(chiCards, cards, maxN, false)
	if realN < 0 {
		realN = int32(len(cards) * 2)
	}
	return realN, delN
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
	return rNum
}
