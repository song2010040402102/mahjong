package logic

import (
	"fmt"
	"sort"
	"testing"
	"time"
	"util"
)

func TestBatchGetCardForRobot(t *testing.T) {
	rule := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	rule.Init(0, 0)
	rule.SetLaiziCard([]int32{101})
	ai := &AIMjBase{}
	ai.SetRule(rule)
	ai.SetLevel(ROBOT_LEVEL_MAJOR)
	ai.SetAllCard(make(map[int32][]int32), make(map[int32][]*ChiCard), make(map[int32][]int32), []int32{})
	for i := 0; i < 1; i++ {
		cards := CreateRandCards(14)
		sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })

		ai.holdCards[1] = cards
		now := time.Now().UnixNano()
		aic := ai.GetCardForRobot(1, 0, cards)

		strRes := ""
		for _, vv := range cards {
			strRes += card2str(vv) + " "
		}
		strRes += fmt.Sprintf("-->> %s  %d\n\n", card2str(aic), time.Now().UnixNano()-now)
		appendToFile("ai.txt", strRes)
	}
}

func TestGetChiForRobot(t *testing.T) {
	rule := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	rule.Init(0, 0)
	rule.SetLaiziCard([]int32{101})
	ai := &AIMjBase{}
	ai.SetRule(rule)
	ai.SetLevel(ROBOT_LEVEL_MAJOR)
	ai.SetAllCard(make(map[int32][]int32), make(map[int32][]*ChiCard), make(map[int32][]int32), []int32{})

	chiCard := NewChiCard()
	chiCard.CardType = MJ_CHI_PENG
	chiCard.CardId = 201
	ai.chiCards[1] = []*ChiCard{chiCard}

	anChi := NewChiCard()
	anChi.CardType = MJ_CHI_GANG_AN
	anChi.CardId = 202
	ai.holdCards[1] = []int32{202, 202, 202, 204}
	ret := ai.GetChiForRobot(1, 202, []*ChiCard{anChi})
	fmt.Println("TestGetChiForRobot, GANG_AN, cards: ", ai.holdCards[1], " ret: ", *ret)

	wanChi := NewChiCard()
	wanChi.CardType = MJ_CHI_GANG_WAN
	wanChi.CardId = 201
	ai.holdCards[1] = []int32{202, 203, 205, 206}
	ret = ai.GetChiForRobot(1, 201, []*ChiCard{wanChi})
	fmt.Println("TestGetChiForRobot, GANG_WAN, cards: ", ai.holdCards[1], " ret: ", *ret)

	gangChi, pengChi, chiChi := NewChiCard(), NewChiCard(), NewChiCard()
	gangChi.CardType, pengChi.CardType, chiChi.CardType = MJ_CHI_GANG, MJ_CHI_PENG, MJ_CHI_CHI
	gangChi.CardId, pengChi.CardId, chiChi.CardId = 201, 201, 201
	ai.holdCards[1] = []int32{201, 201, 201, 202, 203, 205, 206}
	ret = ai.GetChiForRobot(1, 0, []*ChiCard{gangChi, pengChi, chiChi})
	fmt.Println("TestGetChiForRobot, GANG_PENG_CHI, cards: ", ai.holdCards[1], " ret: ", *ret)
}

func CreateRandCards(num int32) []int32 {
	cc := []int32{101, 102, 103, 104, 105, 106, 107, 108, 109, 201, 202, 203, 204, 205, 206, 207, 208, 209, 301, 302, 303, 304, 305, 306, 307, 308, 309, 411, 412, 413, 414, 415, 416, 417}
	mc := make(map[int32]int32)
	cards := make([]int32, 0, num)
	for {
		//c := cc[util.GetRandom(0, int32(len(cc)-1))]
		c := cc[util.GetRandom(18, 28)]
		if mc[c] >= 4 {
			continue
		}
		mc[c]++
		cards = append(cards, c)
		if int32(len(cards)) == num {
			break
		}
	}
	return cards
}
