package logic

import (
	"github.com/astaxie/beego/logs"
	"math"
	"sort"
	"util"
)

type AIMjSC struct {
	AIMjBase
	colors map[int32]int32
}

func (ai *AIMjSC) SetColors(colors map[int32]int32) {
	ai.colors = colors
	logs.Info("[AIMjSC]SetColors, colors: ", colors)
}

func (ai *AIMjSC) GetDealCards() []int32 {
	return []int32{}
}

func (ai *AIMjSC) GetNextCards(index int32) []int32 {
	if ai.level == ROBOT_LEVEL_NOOB {
		if ai.isOneTing(index) {
			return ai.AIMjBase.getUnusedCards(index)
		}
	} else if ai.level == ROBOT_LEVEL_MASTER {
		if util.GetRandomRate() < 0.5 { //一半的概率能摸到想要的牌
			aiCards, _ := ConvertSliceToAICard(ai.holdCards[index], 0, ai.rule.GetLaiziCard())
			aiCards = util.RemoveSliceElem2(aiCards, func(i int) bool { return aiCards[i].Card/MAHJONG_MASK == ai.colors[index] }, true).([]AICard)
			sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Num > aiCards[j].Num })

			nxtCards := make([]int32, 0, len(aiCards))
			for _, v := range aiCards {
				nxtCards = append(nxtCards, v.Card)
			}
			return nxtCards
		}
	}
	return []int32{}
}

func (ai *AIMjSC) GetCardForRobot(index int32, moCard int32) int32 {
	aiCards, _ := ConvertSliceToAICard(ai.holdCards[index], moCard, ai.rule.GetLaiziCard())
	for _, v := range aiCards {
		if v.Card/MAHJONG_MASK == ai.colors[index] { //优先打定缺
			return v.Card
		}
	}
	if ai.level == ROBOT_LEVEL_NOOB || ai.level == ROBOT_LEVEL_AMATEUR {
		return ai.AIMjBase.getCardForWeight(ai.holdCards[index], moCard)
	} else if ai.level == ROBOT_LEVEL_MAJOR || ai.level == ROBOT_LEVEL_MASTER {
		tendColor, _ := ai.getTendColor(index, moCard)
		if tendColor > 0 { //往清一色方向打牌
			valCards := GetCardValue(ai.rule, ai.holdCards[index], moCard)
			for _, v := range valCards {
				if v.Card/MAHJONG_MASK != tendColor {
					return v.Card
				}
			}
		}
		duiNum := ai.getTendDuiNum(index, moCard)
		if duiNum >= 5 { //往七对和对对胡方向打牌
			twoCards, threeCards := []int32{}, []int32{}
			for _, v := range aiCards {
				if v.Num == 1 {
					return v.Card
				} else if v.Num == 2 {
					twoCards = append(twoCards, v.Card)
				} else if v.Num == 3 {
					threeCards = append(threeCards, v.Card)
				}
			}
			if len(ai.chiCards[index]) == 0 { //七对方向
				if len(threeCards) > 0 { //七对优先拆刻子
					return threeCards[0]
				}
			} else { //对对胡方向
				if len(twoCards) > 0 { //对对胡优先拆对子
					return twoCards[0]
				} else if len(threeCards) > 0 {
					return threeCards[0]
				}
			}
		}
		return ai.AIMjBase.getCardForTingNum(ai.holdCards[index], moCard, index, ai.chiCards[index]) //其它情况暂时用向听分析，实际可能采用数量分析更好，但有可能赢不了
	} else {
		logs.Error("[AIMjSC]GetCardForRobot, level invalid!")
	}
	return 0
}

func (ai *AIMjSC) GetChiForRobot(index int32, moCard int32, chiCards []*ChiCard) *ChiCard {
	if len(chiCards) == 0 || chiCards[0] == nil {
		return nil
	}
	chiCard := NewChiCard()
	*chiCard = *chiCards[0]
	chiCard.ChiPosBit = 0
	if ai.level == ROBOT_LEVEL_NOOB { //菜鸟级机器人在一上一听阶段截住吃和碰
		if chiCard.CardType == MJ_CHI_PENG || chiCard.CardType == MJ_CHI_CHI {
			if ai.isOneTing(index) {
				chiCard.CardType = MJ_CHI_PASS
			}
		}
	} else if ai.level == ROBOT_LEVEL_AMATEUR { //业余级机器人按胡>杠>碰>吃优先级

	} else if ai.level == ROBOT_LEVEL_MAJOR || ai.level == ROBOT_LEVEL_MASTER { //专业和大师级机器人根据牌的状态进行分析
		aiCards, _ := ConvertSliceToAICard(ai.holdCards[index], moCard, ai.rule.GetLaiziCard())
		longNum := GetLongNum(ai.chiCards[index], aiCards)
		duiNum := ai.getTendDuiNum(index, moCard)
		tendColor, tendRate := ai.getTendColor(index, moCard)
		if chiCard.CardType == MJ_CHI_HU {
			if moCard > 0 { //自摸直接胡牌
				return chiCard
			}
			if ai.level == ROBOT_LEVEL_MASTER && ai.getTingNumForRestCard(index) == 0 { //大师级机器人判断剩余牌是否有听牌，没有则胡牌
				return chiCard
			}
			if len(ai.restCards) < 10 && chiCard.CardType == MJ_CHI_HU { //剩余牌少于10张直接胡牌
				return chiCard
			}
			if longNum > 0 { //有根直接胡牌
				return chiCard
			}
			if duiNum >= 5 { //七对或对对胡直接胡牌
				return chiCard
			}
			if tendRate == 1 { //清一色直接胡牌
				return chiCard
			}
		}
		//太冒险先屏蔽
		/*if len(ai.chiCards[index]) == 0 && duiNum >= 5 { //避免破坏七对
			chiCard.CardType = MJ_CHI_PASS
			return chiCard
		}*/
		selChiCards := []*ChiCard{}
		for _, v := range chiCards {
			if tendColor != 0 && tendColor != v.CardId/MAHJONG_MASK { //避免破坏清一色
				continue
			}
			selChiCards = append(selChiCards, v)
		}
		if len(selChiCards) == 0 {
			chiCard.CardType = MJ_CHI_PASS
			return chiCard
		} else {
			return selChiCards[0]
		}
	} else {
		logs.Error("[AIMjSC]GetChiForRobot, level invalid!")
		return nil
	}
	return chiCard
}

func (ai *AIMjSC) isOneTing(index int32) bool {
	colorNum := int32(0)
	for _, v := range ai.holdCards[index] {
		if v/MAHJONG_MASK == ai.colors[index] {
			colorNum++
		}
	}
	if colorNum <= 1 {
		cards := make([]int32, len(ai.holdCards[index]))
		copy(cards, ai.holdCards[index])
		if colorNum == 1 {
			for i := 0; i < len(cards); i++ {
				if cards[i]/MAHJONG_MASK == ai.colors[index] {
					cards[i] = 0
				}
			}
		}
		tings := ai.rule.CheckTing(ai.chiCards[index], cards, false)
		if len(tings) > 0 {
			return true
		}
	}
	return false
}

func (ai *AIMjSC) getTendColor(index int32, moCard int32) (int32, float32) {
	if !CheckQingYiSe(ai.chiCards[index], []AICard{}) {
		return int32(0), float32(0)
	}
	if len(ai.restCards) < 10 { //牌少于10张不考虑清一色
		return int32(0), float32(0)
	}
	maxColor, maxColorNum, mapColorNum := int32(0), int32(0), map[int32]int32{}
	aiCards, _ := ConvertSliceToAICard(ai.holdCards[index], moCard, ai.rule.GetLaiziCard())
	for _, v := range aiCards {
		mapColorNum[v.Card/MAHJONG_MASK] += v.Num
	}
	if len(ai.chiCards[index]) > 0 {
		mapColorNum[ai.chiCards[index][0].CardId/MAHJONG_MASK] += int32(3 * len(ai.chiCards[index]))
	}
	for k, v := range mapColorNum {
		if v > maxColorNum {
			maxColor, maxColorNum = k, v
		}
	}
	if maxColorNum >= ai.rule.GetPlayerCardNum()-1 {
		return maxColor, float32(1)
	}
	sameRate, growRate := float64(0.7), float64(0.3)
	sumCardNum := int32(len(ai.rule.GetDeck().GetCards())) - ai.rule.GetPlayerCardNum()*ai.rule.GetPlayerLimit() - 1
	tendRate := float32(maxColorNum) / float32(ai.rule.GetPlayerCardNum())
	if float64(tendRate) > sameRate*math.Pow(float64(float32(sumCardNum)/float32(len(ai.restCards))), growRate) {
		return maxColor, tendRate
	}
	return int32(0), float32(0)
}

func (ai *AIMjSC) getTendDuiNum(index int32, moCard int32) int32 {
	if len(ai.restCards) < 15 { //牌少于15张不考虑往对子方向发展
		return int32(0)
	}
	duiNum := int32(len(ai.chiCards[index]) * 3 / 2)
	aiCards, _ := ConvertSliceToAICard(ai.holdCards[index], moCard, ai.rule.GetLaiziCard())
	for _, v := range aiCards {
		duiNum += v.Num / 2
	}
	return duiNum
}

func (ai *AIMjSC) getTingNumForRestCard(index int32) int32 {
	tingNum := int32(0)
	tings := ai.rule.CheckTing(ai.chiCards[index], ai.holdCards[index], true)
	for _, v1 := range tings {
		for _, v2 := range ai.restCards {
			if v1 == v2 {
				tingNum++
			}
		}
	}
	return tingNum
}
