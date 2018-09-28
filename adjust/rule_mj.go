package logic

import (
	//"github.com/astaxie/beego/logs"
	"sort"
	"util"
)

type RuleMahjong struct {
	MahjongBase
	Condition int32
	cardDeck  *MahjongDeck
	Checker   IHuChecker
	LaiziCard int32
}

func NewRuleMahjong(mjType int32) *RuleMahjong {
	rule := &RuleMahjong{
		Condition: 0,
		LaiziCard: 0,
	}
	rule.MjType = mjType
	if rule.MjType == RULE_SC_MAHJONG_TWO_TWO {
		rule.PlayerLimit = 2
		rule.cardDeck = NewMahjongDeck(true, 0)
		rule.Checker = New_SCMJ_HuChecker()
	} else if rule.MjType == RULE_SC_MAHJONG_TWO_THREE {
		rule.PlayerLimit = 2
		rule.cardDeck = NewMahjongDeck(false, 0)
		rule.Checker = New_SCMJ_HuChecker()
	} else if rule.MjType == RULE_SC_MAHJONG_THREE_TWO {
		rule.PlayerLimit = 3
		rule.cardDeck = NewMahjongDeck(true, 0)
		rule.Checker = New_SCMJ_HuChecker()
	} else if rule.MjType == RULE_SC_MAHJONG_THREE_THREE {
		rule.PlayerLimit = 3
		rule.cardDeck = NewMahjongDeck(false, 0)
		rule.Checker = New_SCMJ_HuChecker()
	} else if rule.MjType == RULE_SC_MAHJONG_XUEZHAN {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, 0)
		rule.Checker = New_SCMJ_HuChecker()
	} else if rule.MjType == RULE_SC_MAHJONG_XUELIU {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, 0)
		rule.Checker = New_SCMJ_HuChecker()
	} else if rule.MjType == RULE_HN_MAHJONG_HONGZHONG {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, WITH_LAIZI_HONGZHONG)
		rule.Checker = New_HNMJ_HONGZHONG_HuChecker()
		rule.LaiziCard = COLOR_OTHER*MAHJONG_MASK + int32(MAHJONG_HONGZHONG)
	} else if rule.MjType == RULE_JS_MAHJONG_SUZHOU {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, WITH_DONG|WITH_SPRING|WITH_MEI|WITH_BAIDA|WITH_BAIBAN)
		rule.Checker = New_HNMJ_HONGZHONG_HuChecker()
	} else if rule.MjType == RULE_JS_MAHJONG_KUNSHAN {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, WITH_DONG|WITH_SPRING|WITH_MEI|WITH_BAIDA|WITH_BAIBAN)
		rule.Checker = New_JSMJ_KUNSHAN_HuChecker()
	} else if rule.MjType == RULE_JS_MAHJONG_QIDONG {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, WITH_DONG|WITH_SPRING|WITH_MEI|WITH_BAIDA|WITH_BAIBAN)
		rule.Checker = New_JSMJ_QIDONG_HuChecker()
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_HY ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_SJ ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LQ {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, WITH_DONG)
		rule.Checker = New_ZJMJ_TAIZHOU_HuChecker()
	} else if rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, WITH_DONG)
		rule.Checker = New_ZJMJ_JINHUA_HuChecker()
	} else if rule.MjType == RULE_ZJ_MAHJONG_HANGZHOU {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, WITH_DONG|WITH_SPRING|WITH_MEI)
		rule.Checker = New_ZJMJ_HANGZHOU_HuChecker()
	} else if rule.MjType == RULE_ZJ_MAHJONG_HUZHOU {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, WITH_DONG)
		rule.Checker = New_ZJMJ_HUZHOU_HuChecker()
	} else if rule.MjType == RULE_ZJ_MAHJONG_JIAXING {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, WITH_DONG)
		rule.Checker = New_ZJMJ_JIAXING_HuChecker()
	} else if rule.MjType == RULE_ZJ_MAHJONG_LISHUI {
		rule.PlayerLimit = 4
		rule.cardDeck = NewMahjongDeck(false, WITH_DONG)
		rule.Checker = New_ZJMJ_LISHUI_HuChecker()
	}

	return rule
}

func (rule *RuleMahjong) GetDeck() *MahjongDeck {
	return rule.cardDeck
}

func (rule *RuleMahjong) GetChecker() IHuChecker {
	return rule.Checker
}

func (rule *RuleMahjong) GetPlayerLimit() int32 {
	return rule.PlayerLimit
}

func (rule *RuleMahjong) SetPlayerLimit(val int32) {
	rule.PlayerLimit = val
}

func (rule *RuleMahjong) SetCondition(val int32) {
	rule.Condition = val
}

func (rule *RuleMahjong) GetLaiziCard() int32 {
	return rule.LaiziCard
}

func (rule *RuleMahjong) SetLaiziCard(val int32) {
	rule.LaiziCard = val
}

func (rule *RuleMahjong) IsHasLaiziHu() bool {
	return rule.MjType == RULE_HN_MAHJONG_HONGZHONG ||
		rule.MjType == RULE_JS_MAHJONG_KUNSHAN ||
		rule.MjType == RULE_JS_MAHJONG_QIDONG ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_HY ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_SJ ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LQ ||
		rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW ||
		rule.MjType == RULE_ZJ_MAHJONG_HANGZHOU ||
		rule.MjType == RULE_ZJ_MAHJONG_HUZHOU ||
		rule.MjType == RULE_ZJ_MAHJONG_JIAXING ||
		rule.MjType == RULE_ZJ_MAHJONG_LISHUI
}

func (rule *RuleMahjong) IsHasQidui() bool {
	return (rule.MjType >= RULE_SC_MAHJONG_TWO_TWO && rule.MjType <= RULE_SC_MAHJONG_TOP) ||
		(rule.MjType >= RULE_JS_MAHJONG_SUZHOU && rule.MjType <= RULE_JS_MAHJONG_TOP) ||
		rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW ||
		rule.MjType == RULE_ZJ_MAHJONG_HANGZHOU ||
		rule.MjType == RULE_ZJ_MAHJONG_LISHUI
}

func (rule *RuleMahjong) IsHasQuanBuKao() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW ||
		rule.MjType == RULE_ZJ_MAHJONG_LISHUI
}

func (rule *RuleMahjong) CanReplaceForBai() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW ||
		rule.MjType == RULE_ZJ_MAHJONG_HANGZHOU ||
		rule.MjType == RULE_ZJ_MAHJONG_LISHUI ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_HY ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_SJ ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LQ
}

func (rule *RuleMahjong) IsCanGang(holdCards []int32, card int32) bool {
	if card == rule.GetLaiziCard() {
		return false
	}
	cnt := 0
	for _, v := range holdCards {
		if v == card {
			cnt++
		}
	}
	if cnt >= 3 {
		return true
	}
	return false
}

func (rule *RuleMahjong) IsCanPeng(holdCards []int32, card int32) bool {
	if card == rule.GetLaiziCard() {
		return false
	}
	cnt := 0
	for _, v := range holdCards {
		if v == card {
			cnt++
		}
	}
	if cnt >= 2 {
		return true
	}
	return false
}

func (rule *RuleMahjong) IsCanChi(holdCards []int32, card int32) (bool, int32) {
	if card == rule.GetLaiziCard() {
		return false, 0
	}
	if rule.CanReplaceForBai() == true {
		if card == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
			card = rule.LaiziCard
		}
	}

	seq := GetMahjongSeq(card)
	if seq < int32(MAHJONG_1) || seq > int32(MAHJONG_9) {
		return false, 0
	}
	groups := make(map[int32]bool)
	for _, v := range holdCards {
		if v == rule.LaiziCard {
			continue //自己牌中存在癞子牌不能参与吃牌
		}
		if rule.CanReplaceForBai() == true {
			if v == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
				v = rule.LaiziCard
			}
		}
		if v == card-1 || v == card-2 || v == card+1 || v == card+2 {
			groups[v] = true
		}
	}
	if len(groups) <= 1 {
		return false, 0
	}
	var chiPosBit int32 = 0
	var found1, found2, found3 bool = true, true, true
	for i := card - 2; i <= card-1; i++ {
		if _, ok := groups[i]; !ok {
			found1 = false
		}
	}
	for i := card - 1; i <= card+1; i++ {
		if i == card {
			continue
		}
		if _, ok := groups[i]; !ok {
			found2 = false
		}
	}
	for i := card + 1; i <= card+2; i++ {
		if _, ok := groups[i]; !ok {
			found3 = false
		}
	}
	if found1 {
		chiPosBit |= 0x01 << 1
	}
	if found2 {
		chiPosBit |= 0x01 << 2
	}
	if found3 {
		chiPosBit |= 0x01 << 3
	}
	return found1 || found2 || found3, chiPosBit
}

func (rule *RuleMahjong) GetChiCard(chiCard *ChiCard) []int32 {
	var cards []int32
	if chiCard.CardType == MJ_CHI_GANG || chiCard.CardType == MJ_CHI_GANG_WAN || chiCard.CardType == MJ_CHI_GANG_AN {
		cards = []int32{chiCard.CardId, chiCard.CardId, chiCard.CardId, chiCard.CardId}
	} else if chiCard.CardType == MJ_CHI_PENG {
		cards = []int32{chiCard.CardId, chiCard.CardId, chiCard.CardId}
	} else if chiCard.CardType == MJ_CHI_CHI {
		c := chiCard.CardId
		if chiCard.CardId == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
			c = rule.GetLaiziCard()
		}
		c1, c2, c3 := int32(0), int32(0), int32(0)
		if chiCard.ChiPosBit&(0x01<<1) != 0 {
			c1, c2, c3 = c-2, c-1, chiCard.CardId
		} else if chiCard.ChiPosBit&(0x01<<2) != 0 {
			c1, c2, c3 = c-1, chiCard.CardId, c+1
		} else if chiCard.ChiPosBit&(0x01<<3) != 0 {
			c1, c2, c3 = chiCard.CardId, c+1, c+2
		}
		if c2 == rule.GetLaiziCard() {
			c2 = COLOR_OTHER*MAHJONG_MASK + MAHJONG_BAI
		}
		if c3 == rule.GetLaiziCard() {
			c3 = COLOR_OTHER*MAHJONG_MASK + MAHJONG_BAI
		}
		cards = []int32{c1, c2, c3}
	}
	return cards
}

func (rule *RuleMahjong) CheckHuType(chiCards []*ChiCard, holdCards []int32, card int32, all bool) (uint8, [][][3]int32) {
	var groups [][][3]int32
	if card != 0 && len(holdCards)%3 != 1 || card == 0 && len(holdCards)%3 != 2 {
		return MJ_HU_TYPE_NONE, groups
	}

	aiCards, lzNum := ConvertSliceToAICard(holdCards, card, rule.LaiziCard)

	//确定胡法
	b_lz := rule.IsHasLaiziHu()     //是否允许癞子胡
	b_bt := rule.CanReplaceForBai() //是否允许白板替换癞子
	b_qd := rule.IsHasQidui()       //是否允许七对胡
	b_bk := rule.IsHasQuanBuKao()   //是否允许全不靠胡

	//先检查特殊胡法，因为情况单一，所以不需要groups信息
	if b_qd == true && len(chiCards) <= 0 {
		if b_lz == true && lzNum > 0 {
			if CheckQiDuiHuForLZ(aiCards, lzNum) == true {
				return MJ_HU_TYPE_QIDUI, groups
			}
		} else {
			if CheckQiDuiHu(aiCards) == true {
				return MJ_HU_TYPE_QIDUI, groups
			}
		}
	}
	if b_bk == true && len(chiCards) <= 0 && lzNum <= 0 && CheckQuanBuKaoHu(aiCards) == true {
		return MJ_HU_TYPE_QUANBUKAO, groups
	}

	//后检查普通胡法
	if b_lz == true {
		if b_bt == true {
			ReplaceBaiWithLZ([]*ChiCard{}, aiCards, rule.LaiziCard)
		}
		//logs.Debug("CheckHuType, b_bt: ", b_bt, "holdCards: ", holdCards, "card: ", card, "aiCards: ", aiCards, "lzNum: ", lzNum, "lzCard: ", rule.LaiziCard)
		if lzNum > 0 {
			if ok, groups := CheckCommonHuForLZ(aiCards, lzNum, all); ok {
				return MJ_HU_TYPE_COMMON, groups
			}
		} else {
			if ok, groups := CheckCommonHu(aiCards, all); ok {
				return MJ_HU_TYPE_COMMON, groups
			}
		}
	} else {
		if ok, groups := CheckCommonHu(aiCards, all); ok {
			return MJ_HU_TYPE_COMMON, groups
		}
	}
	return MJ_HU_TYPE_NONE, groups
}

func (rule *RuleMahjong) CheckTing(chiCards []*ChiCard, holdCards []int32, all bool) []int32 {
	tingInfo := make([]int32, 0, 34)
	if len(holdCards)%3 != 1 {
		return tingInfo
	}

	aiCards, lzNum := ConvertSliceToAICard(holdCards, 0, rule.LaiziCard)

	//确定胡法
	b_lz := rule.IsHasLaiziHu()     //是否允许癞子胡
	b_bt := rule.CanReplaceForBai() //是否允许白板替换癞子
	b_qd := rule.IsHasQidui()       //是否允许七对胡
	b_bk := rule.IsHasQuanBuKao()   //是否允许全不靠胡

	//先检查特殊听
	if b_qd && len(chiCards) <= 0 {
		if b_lz && lzNum > 0 {
			tingInfo = append(tingInfo, CheckQiDuiTingForLZ(aiCards, lzNum)...)
		} else {
			tingInfo = append(tingInfo, CheckQiDuiTing(aiCards)...)
		}
		if len(tingInfo) > 0 && (!all || tingInfo[0] == MAHJONG_ANY) {
			return tingInfo
		}
	}
	if b_bk && len(chiCards) <= 0 && lzNum <= 0 {
		tingInfo = append(tingInfo, CheckQuanBuKaoTing(aiCards)...)
		if len(tingInfo) > 0 && !all {
			return tingInfo
		}
	}

	//后检查普通听
	if b_lz {
		if b_bt && rule.LaiziCard%MAHJONG_MASK >= MAHJONG_1 && rule.LaiziCard%MAHJONG_MASK <= MAHJONG_9 {
			//若白板可以替癞子牌，应提前换掉
			for i := 0; i < len(aiCards); i++ {
				if aiCards[i].Card == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
					aiCards[i].Card = rule.LaiziCard
					sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })
					break
				}
			}
		}
		if lzNum > 0 {
			tingInfo = append(tingInfo, CheckCommonTingForLZ(aiCards, lzNum, all)...)
			if len(tingInfo) > 0 && tingInfo[0] == MAHJONG_ANY {
				return tingInfo
			}
		} else {
			tingInfo = append(tingInfo, CheckCommonTing(aiCards, all)...)
		}
	} else {
		tingInfo = append(tingInfo, CheckCommonTing(aiCards, all)...)
	}

	if len(tingInfo) > 0 {
		//处理最终听牌结果
		if b_lz {
			hasLZ := false
			for _, v := range tingInfo {
				if v == rule.GetLaiziCard() {
					hasLZ = true
					break
				}
			}
			if hasLZ {
				if b_bt {
					tingInfo = append(tingInfo, COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI)
				}
			} else {
				tingInfo = append(tingInfo, rule.GetLaiziCard())
			}
		}
		sort.Slice(tingInfo, func(i, j int) bool { return tingInfo[i] < tingInfo[j] })
		tingInfo = util.UniqueSlice(tingInfo)
	}
	return tingInfo
}

/* 检查N向听牌，N: 0,1,2,3...，返回实际向听数及相应的向听牌
** N为0，等价于CheckTing函数，N为1，一上一听，依次类推，时间复杂度为组合数C{M,N}*Ot，M手牌数，Ot为检查听牌的时间复杂度
** 手牌数为14，最坏情况为四上一听，手牌数为17，最坏情况为五上一听，依次类推
 */
func (rule *RuleMahjong) CheckNTing(chiCards []*ChiCard, holdCards []int32, N int32, all bool) (int32, []int32) {
	if len(holdCards)%3 != 1 || N < 0 {
		return -1, []int32{}
	}

	cards := make([]int32, len(holdCards))
	copy(cards, holdCards)
	sort.Slice(cards, func(i, j int) bool { return cards[i] < cards[j] })

	indexs := make([]int, 0, len(cards))
	for i := 0; i < len(cards); i++ {
		if cards[i] == rule.GetLaiziCard() || i > 0 && cards[i] == cards[i-1] {
			continue
		}
		indexs = append(indexs, i)
	}

	//对indexs按权重预排序会提高剪枝效率，能加快搜索树的返回，对于实际听数小于检测听数的剪枝效率很高，其它情况效率差不多
	cardVal := GetCardValue(rule, cards, 0)
	mapVal := make(map[int32]float32)
	for _, v := range cardVal {
		mapVal[v.Card] = v.Val
	}
	sort.Slice(indexs, func(i, j int) bool { return mapVal[cards[indexs[i]]] < mapVal[cards[indexs[j]]] })

	baseN, tings := rule.doCheckNTing(chiCards, cards, N, indexs, 0, all)
	if len(tings) == 0 {
		return -1, tings
	}
	return N - baseN, tings
}

func (rule *RuleMahjong) doCheckNTing(chiCards []*ChiCard, cards []int32, N int32, indexs []int, minIndex int, all bool) (int32, []int32) {
	tings := rule.CheckTing(chiCards, cards, all)
	if len(tings) > 0 || N == 0 {
		return N, tings
	}

	//由于麻将牌型的复杂性，目前只能抽取牌，然后加上癞子牌，取听牌的并集
	maxBaseN := int32(0)
	for k := minIndex; k < len(indexs); k++ {
		bkc := cards[indexs[k]]
		cards[indexs[k]] = rule.GetLaiziCard()
		baseN, subTings := rule.doCheckNTing(chiCards, cards, N-maxBaseN-1, indexs, k+1, all)
		if baseN > 0 {
			maxBaseN += baseN
			tings = tings[0:0]
		}
		tings = append(tings, subTings...)
		cards[indexs[k]] = bkc
	}

	if maxBaseN == 0 && N == 2 {
		duiIndex := make([]int, 0, len(indexs))
		for k := 0; k < len(indexs); k++ {
			if cards[indexs[k]] != rule.GetLaiziCard() && indexs[k] < len(cards)-1 && cards[indexs[k]+1] == cards[indexs[k]] &&
				(indexs[k] > len(cards)-3 || cards[indexs[k]+1] != cards[indexs[k]+2]) {
				duiIndex = append(duiIndex, k)
			}
		}
		if len(duiIndex) >= 5 { //当对子数超过5个时，需要尝试拆对子来检查听
			for _, v := range duiIndex {
				bkc := cards[indexs[v]]
				cards[indexs[v]], cards[indexs[v]+1] = rule.GetLaiziCard(), rule.GetLaiziCard()
				_, subTings := rule.doCheckNTing(chiCards, cards, 0, indexs, 0, all)
				tings = append(tings, subTings...)
				if len(tings) > 0 && !all {
					return 0, tings
				}
				cards[indexs[v]], cards[indexs[v]+1] = bkc, bkc
			}
		}
	}

	if len(tings) > 1 {
		sort.Slice(tings, func(i, j int) bool { return tings[i] < tings[j] })
		tings = util.UniqueSlice(tings)
	}
	return maxBaseN, tings
}
