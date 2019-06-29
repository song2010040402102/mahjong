package logic

import (
	//"github.com/astaxie/beego/logs"
	"sort"
	"util"
)

type RuleMahjong struct {
	MahjongBase
	cardDeck  *MahjongDeck
	huHandler IHuHandler
	lzCards   []int32
	switchSet uint64
}

func NewRuleMahjong(mjType int32) *RuleMahjong {
	rule := &RuleMahjong{}
	rule.MjType = mjType
	rule.PlayerLimit = 4
	return rule
}

func (rule *RuleMahjong) Init(playerLimit int32, switchSet uint64) {
	rule.switchSet = switchSet & MJ_SETTING_RULE
	if rule.MjType >= RULE_SC_MAHJONG_HEAD && rule.MjType < RULE_SC_MAHJONG_TAIL {
		rule.cardDeck = NewMahjongDeck(WITH_DONG | WITH_ZHONG | WITH_SPRING | WITH_MEI)
		rule.huHandler = New_SCMJ_HuHandler()
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_HY ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_JJ ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LQ {
		rule.cardDeck = NewMahjongDeck(WITH_SPRING | WITH_MEI)
		rule.huHandler = New_ZJMJ_TAIZHOU_HuHandler()
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LH {
		rule.cardDeck = NewMahjongDeck(WITH_SPRING | WITH_MEI)
		rule.huHandler = New_ZJMJ_TAIZHOU_LH_HuHandler()
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_SM {
		rule.cardDeck = NewMahjongDeck(0)
		rule.huHandler = New_ZJMJ_TAIZHOU_SM_HuHandler()
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_YH {
		rule.cardDeck = NewMahjongDeck(WITH_SPRING | WITH_MEI)
		rule.huHandler = New_ZJMJ_TAIZHOU_YH_HuHandler()
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_TT {
		rule.PlayerLimit = 3
		rule.cardDeck = NewMahjongDeck(WITH_TIAO)
		rule.huHandler = New_ZJMJ_TAIZHOU_TT_HuHandler()
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_DQ {
		rule.cardDeck = NewMahjongDeck(0)
		rule.huHandler = New_ZJMJ_TAIZHOU_DQ_HuHandler()
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL {
		if rule.switchSet&MJ_SETTING_NAO_SMALL != 0 || rule.switchSet&MJ_SETTING_NAO_BIG != 0 {
			rule.cardDeck = NewMahjongDeck(0)
		} else {
			rule.cardDeck = NewMahjongDeck(WITH_SPRING | WITH_MEI)
		}
		rule.huHandler = New_ZJMJ_TAIZHOU_WL_HuHandler()
	} else if rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW {
		rule.cardDeck = NewMahjongDeck(WITH_SPRING | WITH_MEI)
		rule.huHandler = New_ZJMJ_JINHUA_HuHandler()
	} else if rule.MjType == RULE_ZJ_MAHJONG_NINGBO {
		rule.cardDeck = NewMahjongDeck(0)
		rule.huHandler = New_ZJMJ_NINGBO_HuHandler()
	} else if rule.MjType == RULE_ZJ_MAHJONG_JIAXING {
		rule.cardDeck = NewMahjongDeck(WITH_SPRING | WITH_MEI)
		rule.huHandler = New_ZJMJ_JIAXING_HuHandler()
	}
	if playerLimit > 0 {
		rule.PlayerLimit = playerLimit
	}
}

func (rule *RuleMahjong) GameOver() {
	rule.lzCards = []int32{}
}

func (rule *RuleMahjong) GetDeck() *MahjongDeck {
	return rule.cardDeck
}

func (rule *RuleMahjong) GetHuHandler() IHuHandler {
	return rule.huHandler
}

func (rule *RuleMahjong) GetPlayerLimit() int32 {
	return rule.PlayerLimit
}

func (rule *RuleMahjong) GetSwitchSet() uint64 {
	return rule.switchSet
}

func (rule *RuleMahjong) GetPlayerDoor(index int32) int32 {
	if rule.PlayerLimit == 2 { //两人麻将东西位
		if index == 2 {
			index = index + 1
		}
	} else if rule.PlayerLimit == 3 { //三人麻将东南北位
		if index == 3 {
			index = index + 1
		}
	}
	return MAHJONG_DONG + index - 1
}

func (rule *RuleMahjong) GetPlayerCardNum() int32 {
	if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_YH || rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL {
		return int32(16)
	}
	return int32(13)
}

func (rule *RuleMahjong) GetRemainCardNum() int32 {
	if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_YH || rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL {
		return int32(16)
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_TT || rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_DQ {
		return int32(14)
	}
	return int32(0)
}

func (rule *RuleMahjong) GetUnknownCardNum() int32 {
	if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_YH {
		return int32(10)
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_DQ || rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL {
		return int32(14)
	}
	return int32(0)
}

func (rule *RuleMahjong) GetFanCardSkipNum() int32 {
	if rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW {
		return 2 //义乌麻将要跳过财神牌的那一蹲牌
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LH || rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_SM ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_TT {
		return 0 //对于固定白搭牌或者翻牌可摸不需要跳过
	}
	return 1
}

func (rule *RuleMahjong) CreateLaiziCard(fanCard int32) []int32 {
	lzCards := make([]int32, 0)
	fanSeq := GetMahjongSeq(fanCard)
	if rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_HY ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_JJ ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LQ ||
		rule.MjType == RULE_ZJ_MAHJONG_NINGBO ||
		rule.MjType == RULE_ZJ_MAHJONG_JIAXING {
		lzCards = append(lzCards, fanCard)
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL {
		start := int32(0)
		if c := fanCard % MAHJONG_MASK; c >= MAHJONG_SPRING && c <= MAHJONG_WINTER {
			start = MAHJONG_MASK*COLOR_OTHER + MAHJONG_SPRING
		} else if c >= MAHJONG_MEI && c <= MAHJONG_JU {
			start = MAHJONG_MASK*COLOR_OTHER + MAHJONG_MEI
		}
		if start > 0 {
			for i := 0; i < 4; i++ {
				lzCards = append(lzCards, start+int32(i))
			}
		} else {
			lzCards = append(lzCards, fanCard)
		}
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_YH {
		baidaSeq := fanSeq
		if fanSeq >= MAHJONG_1 && fanSeq <= MAHJONG_9 {
			baidaSeq--
			if baidaSeq <= 0 {
				baidaSeq = MAHJONG_9
			}
		} else if fanSeq >= MAHJONG_DONG && fanSeq <= MAHJONG_BAI {
			baidaSeq++
			if baidaSeq > MAHJONG_BAI {
				baidaSeq = MAHJONG_DONG
			}
		}
		lzCards = append(lzCards, fanCard-fanSeq+baidaSeq)
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LH || rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_SM ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_TT {
		lzCards = append(lzCards, MAHJONG_MASK*COLOR_OTHER+MAHJONG_BAI)
	}
	return lzCards
}

func (rule *RuleMahjong) GetLaiziCard() []int32 {
	return rule.lzCards
}

func (rule *RuleMahjong) SetLaiziCard(lzCards []int32) {
	rule.lzCards = lzCards
}

func (rule *RuleMahjong) IsLaiziCard(card int32) bool {
	for _, v := range rule.lzCards {
		if v == card {
			return true
		}
	}
	return false
}

func (rule *RuleMahjong) IsHasLaiziHu() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_HY ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_JJ ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LQ ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LH ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_SM ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_YH ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_TT ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL ||
		rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW ||
		rule.MjType == RULE_ZJ_MAHJONG_NINGBO ||
		rule.MjType == RULE_ZJ_MAHJONG_JIAXING
}

func (rule *RuleMahjong) IsHasQidui() bool {
	return rule.MjType >= RULE_SC_MAHJONG_HEAD && rule.MjType < RULE_SC_MAHJONG_TAIL ||
		rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW
}

func (rule *RuleMahjong) IsHasQuanBuKao() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW
}

func (rule *RuleMahjong) IsHas4LaiziHu() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LH ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_TT
}

func (rule *RuleMahjong) IsHas3LaiziHu() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL
}

func (rule *RuleMahjong) IsHas8HuaHu() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_TT ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_DQ ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL ||
		rule.MjType == RULE_ZJ_MAHJONG_NINGBO
}

func (rule *RuleMahjong) IsHas4HuaHu() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL
}

func (rule *RuleMahjong) IsHas4ZiSelfHu() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL
}

func (rule *RuleMahjong) IsCanLaiziSelf() bool {
	//癞子牌是否可以作为原始牌进行吃碰杠胡
	return rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LH ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_SM ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_YH
}

func (rule *RuleMahjong) IsCanBaiReplaceLZ() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_HY ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_JJ ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LQ ||
		rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW
}

func (rule *RuleMahjong) IsCanMultiPlayerHu() bool {
	return rule.MjType >= RULE_SC_MAHJONG_HEAD && rule.MjType < RULE_SC_MAHJONG_TAIL ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LH ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_SM ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_TT
}

func (rule *RuleMahjong) IsCanPlay(card int32) bool {
	if rule.IsHuaCard(card) || rule.IsLaiziCard(card) && rule.MjType != RULE_ZJ_MAHJONG_TAIZHOU_LH && rule.MjType != RULE_ZJ_MAHJONG_TAIZHOU_SM &&
		rule.MjType != RULE_ZJ_MAHJONG_TAIZHOU_YH && rule.MjType != RULE_ZJ_MAHJONG_JINHUA_YW && rule.MjType != RULE_ZJ_MAHJONG_NINGBO {
		return false
	}
	return true
}

func (rule *RuleMahjong) IsCanGang(holdCards []int32, card int32) bool {
	if rule.IsHuaCard(card) || rule.IsLaiziCard(card) && !rule.IsCanLaiziSelf() {
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
	if rule.IsLaiziCard(card) && !rule.IsCanLaiziSelf() {
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
	if rule.MjType >= RULE_SC_MAHJONG_HEAD && rule.MjType < RULE_SC_MAHJONG_TAIL || rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_TT {
		return false, 0
	}
	if rule.IsLaiziCard(card) && !rule.IsCanLaiziSelf() {
		return false, 0
	}

	lzCard := int32(0)
	if len(rule.lzCards) > 0 {
		lzCard = rule.lzCards[0]
	}
	if rule.IsCanBaiReplaceLZ() == true {
		if card == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
			card = lzCard
		}
	}

	seq := GetMahjongSeq(card)
	if seq < int32(MAHJONG_1) || seq > int32(MAHJONG_9) {
		return false, 0
	}
	groups := make(map[int32]bool)
	for _, v := range holdCards {
		if rule.IsLaiziCard(v) && !rule.IsCanLaiziSelf() {
			continue //自己牌中存在癞子牌不能参与吃牌
		}
		if rule.IsCanBaiReplaceLZ() == true {
			if v == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
				v = lzCard
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

func (rule *RuleMahjong) IsCanPlayAfterHu() bool {
	return rule.MjType >= RULE_SC_MAHJONG_HEAD && rule.MjType < RULE_SC_MAHJONG_TAIL
}

func (rule *RuleMahjong) IsCanPlayForWinner() bool {
	return rule.MjType == RULE_SC_MAHJONG_XUELIU
}

func (rule *RuleMahjong) IsCanPlaySameAsChi() bool {
	return rule.MjType != RULE_ZJ_MAHJONG_NINGBO
}

func (rule *RuleMahjong) IsCanPlaySameAsChiExceptDiao() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_NINGBO
}

func (rule *RuleMahjong) IsRealTimeCalc() bool {
	return rule.MjType >= RULE_SC_MAHJONG_HEAD && rule.MjType < RULE_SC_MAHJONG_TAIL
}

func (rule *RuleMahjong) IsEqualCalc() bool {
	return rule.MjType >= RULE_SC_MAHJONG_HEAD && rule.MjType < RULE_SC_MAHJONG_TAIL
}

func (rule *RuleMahjong) IsHuaCard(card int32) bool {
	if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_YH {
		if len(rule.lzCards) > 0 && rule.lzCards[0] == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
			if card == COLOR_OTHER*MAHJONG_MASK+MAHJONG_HONGZHONG {
				return true
			}
		} else {
			if card == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
				return true
			}
		}
	} else if rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_WL {
		if rule.switchSet&MJ_SETTING_NAO_SMALL == 0 && rule.switchSet&MJ_SETTING_NAO_BIG == 0 {
			if len(rule.lzCards) > 0 && rule.lzCards[0] == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
				if c := card % MAHJONG_MASK; c == MAHJONG_HONGZHONG || c == MAHJONG_LVFA {
					return true
				}
			} else {
				if card == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
					return true
				}
			}
		} else if !rule.IsLaiziCard(card) {
			if len(rule.lzCards) > 0 {
				if rule.lzCards[0] == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI ||
					rule.switchSet&MJ_SETTING_NAO_BIG != 0 && (rule.switchSet&MJ_SETTING_NAO_SMALL != 0 ||
						rule.lzCards[0]%MAHJONG_MASK >= MAHJONG_DONG && rule.lzCards[0]%MAHJONG_MASK >= MAHJONG_BEI) {
					if c := card % MAHJONG_MASK; c >= MAHJONG_HONGZHONG && c <= MAHJONG_BAI || c >= MAHJONG_SPRING && c <= MAHJONG_JU {
						return true
					}
				} else {
					if c := card % MAHJONG_MASK; c == MAHJONG_BAI || c >= MAHJONG_SPRING && c <= MAHJONG_JU {
						return true
					}
				}
			}
		}
	} else if c := card % MAHJONG_MASK; c >= MAHJONG_SPRING && c <= MAHJONG_JU {
		return true
	}
	return false
}

func (rule *RuleMahjong) IsLunZhuang() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_HY ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_JJ ||
		rule.MjType == RULE_ZJ_MAHJONG_TAIZHOU_LQ ||
		rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW
}

func (rule *RuleMahjong) IsNeedSwapCard() bool {
	return rule.MjType >= RULE_SC_MAHJONG_HEAD && rule.MjType < RULE_SC_MAHJONG_TAIL
}

func (rule *RuleMahjong) IsNeedChooseColor() bool {
	return rule.MjType >= RULE_SC_MAHJONG_HEAD && rule.MjType < RULE_SC_MAHJONG_TAIL
}

func (rule *RuleMahjong) IsNeedBet() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_JINHUA_YW
}

func (rule *RuleMahjong) IsNeedSameColorForChi() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_JIAXING
}

func (rule *RuleMahjong) IsNeedNotifyBaoFor3Feed() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_JIAXING
}

func (rule *RuleMahjong) IsCanGiveupForBao() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_JIAXING
}

func (rule *RuleMahjong) IsPriorityChiForBao() bool {
	return rule.MjType == RULE_ZJ_MAHJONG_JIAXING
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
			lzCards := rule.GetLaiziCard()
			if len(lzCards) > 0 {
				c = lzCards[0]
			}
		}
		c1, c2, c3 := int32(0), int32(0), int32(0)
		if chiCard.ChiPosBit&(0x01<<1) != 0 {
			c1, c2, c3 = c-2, c-1, chiCard.CardId
		} else if chiCard.ChiPosBit&(0x01<<2) != 0 {
			c1, c2, c3 = c-1, chiCard.CardId, c+1
		} else if chiCard.ChiPosBit&(0x01<<3) != 0 {
			c1, c2, c3 = chiCard.CardId, c+1, c+2
		}
		if rule.IsLaiziCard(c2) {
			c2 = COLOR_OTHER*MAHJONG_MASK + MAHJONG_BAI
		}
		if rule.IsLaiziCard(c3) {
			c3 = COLOR_OTHER*MAHJONG_MASK + MAHJONG_BAI
		}
		cards = []int32{c1, c2, c3}
	}
	return cards
}

func (rule *RuleMahjong) CheckHuHua(index int32, huaCards []int32, flag uint32) uint32 {
	realHuaNum := 0
	otherHua := make(map[int32]int32)
	for _, v := range huaCards {
		if c := v % MAHJONG_MASK; c >= MAHJONG_SPRING && c <= MAHJONG_JU {
			realHuaNum++
		} else {
			otherHua[v]++
		}
	}
	huType := MJ_HU_TYPE_NONE
	if rule.IsHas8HuaHu() && realHuaNum == 8 {
		huType |= MJ_HU_TYPE_8HUA
	} else if rule.IsHas4HuaHu() {
		if CheckFourHua(huaCards) {
			huType |= MJ_HU_TYPE_4HUA
		} else if realHuaNum >= 4 && flag&MJ_CHECK_HU_FLAG_CANHUNHUA != 0 {
			huType |= MJ_HU_TYPE_4HUNHUA
		}
	}
	if rule.IsHas4ZiSelfHu() {
		for k, v := range otherHua {
			if v >= 4 {
				if c := k % MAHJONG_MASK; c >= MAHJONG_HONGZHONG && c <= MAHJONG_BAI || c == rule.GetPlayerDoor(index) {
					huType |= MJ_HU_TYPE_4ZISELF
					break
				}
			}
		}
	}
	return huType
}

func (rule *RuleMahjong) CheckHuType(index int32, chiCards []*ChiCard, huaCards []int32, holdCards []int32, card int32, flag uint32) (uint32, [][][3]int32) {
	var groups [][][3]int32
	if card != 0 && len(holdCards)%3 != 1 || card == 0 && len(holdCards)%3 != 2 {
		return MJ_HU_TYPE_NONE, groups
	}

	if rule.IsHuaCard(card) {
		return MJ_HU_TYPE_NONE, groups
	} else {
		for _, v := range holdCards {
			if rule.IsHuaCard(v) {
				return MJ_HU_TYPE_NONE, groups
			}
		}
	}

	lzNum := int32(0)
	var aiCards []AICard
	if rule.IsLaiziCard(card) && flag&MJ_CHECK_HU_FLAG_ZIMO == 0 {
		if rule.IsCanLaiziSelf() { //处理可作为原始牌的癞子牌的点炮情况
			aiCards, lzNum = ConvertSliceToAICard(holdCards, 0, rule.lzCards)
			aiCards = append(aiCards, AICard{card, 1})
		} else {
			return MJ_HU_TYPE_NONE, groups
		}
	} else {
		aiCards, lzNum = ConvertSliceToAICard(holdCards, card, rule.lzCards)
	}

	//先检查特殊胡法，因为情况单一，所以不需要groups信息
	huType := MJ_HU_TYPE_NONE
	if flag&MJ_CHECK_HU_FLAG_ZIMO != 0 { //必须自摸才可胡
		huType |= rule.CheckHuHua(index, huaCards, flag) //胡花
		if rule.IsHas3LaiziHu() && lzNum >= 3 {          //三癞子胡
			huType |= MJ_HU_TYPE_3LAIZI
		}
		if rule.IsHas4LaiziHu() && lzNum >= 4 { //四癞子胡
			huType |= MJ_HU_TYPE_4LAIZI
		}
		if rule.IsHas4ZiSelfHu() { //胡门风箭
			if lzNum >= 4 { //4个癞子牌充当4个门风箭牌胡牌
				huType |= MJ_HU_TYPE_4ZISELF
			} else {
				for _, v := range chiCards {
					if v.CardType == MJ_CHI_GANG_AN {
						if c := v.CardId % MAHJONG_MASK; c >= MAHJONG_HONGZHONG && c <= MAHJONG_BAI || c == rule.GetPlayerDoor(index) {
							huType |= MJ_HU_TYPE_4ZISELF
							break
						}
					}
				}
				if huType&MJ_HU_TYPE_4ZISELF == 0 {
					for _, v := range aiCards {
						if v.Num >= 4 {
							if c := v.Card % MAHJONG_MASK; c >= MAHJONG_HONGZHONG && c <= MAHJONG_BAI || c == rule.GetPlayerDoor(index) {
								huType |= MJ_HU_TYPE_4ZISELF
								break
							}
						}
					}
				}
			}
		}
	}
	if rule.IsHasQidui() && len(chiCards) <= 0 {
		if CheckQiDuiHuForLZ(aiCards, lzNum) == true {
			huType |= MJ_HU_TYPE_QIDUI
		}
	}
	if rule.IsHasQuanBuKao() && len(chiCards) <= 0 && lzNum <= 0 && CheckQuanBuKaoHu(aiCards) == true {
		huType |= MJ_HU_TYPE_QUANBUKAO
	}

	//后检查普通胡法
	ok := false
	if rule.IsCanBaiReplaceLZ() == true {
		ReplaceBaiWithLZ([]*ChiCard{}, aiCards, rule.lzCards)
	}
	if ok, groups = CheckCommonHuForLZ(aiCards, lzNum, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
		huType |= MJ_HU_TYPE_COMMON
	}
	return huType, groups
}

func (rule *RuleMahjong) CheckTing(chiCards []*ChiCard, holdCards []int32, all bool) []int32 {
	tingInfo := make([]int32, 0, 34)
	if len(holdCards)%3 != 1 {
		return tingInfo
	}

	aiCards, lzNum := ConvertSliceToAICard(holdCards, 0, rule.lzCards)

	//先检查特殊听，这里不能检查与牌型结构无关的听牌，例如四癞子胡或8花胡
	if rule.IsHasQidui() && len(chiCards) <= 0 {
		tingInfo = append(tingInfo, CheckQiDuiTingForLZ(aiCards, lzNum)...)
	}
	if rule.IsHasQuanBuKao() && len(chiCards) <= 0 && lzNum <= 0 {
		tingInfo = append(tingInfo, CheckQuanBuKaoTing(aiCards)...)
	}
	if len(tingInfo) > 0 && (!all || tingInfo[0] == MAHJONG_ANY) {
		return tingInfo
	}

	//后检查普通听
	if rule.IsCanBaiReplaceLZ() && len(rule.lzCards) > 0 && rule.lzCards[0]%MAHJONG_MASK >= MAHJONG_1 && rule.lzCards[0]%MAHJONG_MASK <= MAHJONG_9 {
		//若白板可以替癞子牌，应提前换掉
		for i := 0; i < len(aiCards); i++ {
			if aiCards[i].Card == COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI {
				aiCards[i].Card = rule.lzCards[0]
				sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })
				break
			}
		}
	}
	tingInfo = append(tingInfo, CheckCommonTingForLZ(aiCards, lzNum, all)...)
	if len(tingInfo) > 0 && tingInfo[0] == MAHJONG_ANY {
		return tingInfo
	}

	if len(tingInfo) > 0 {
		//处理最终听牌结果
		if rule.IsHasLaiziHu() {
			hasLZ := false
			for _, v := range tingInfo {
				if rule.IsLaiziCard(v) {
					hasLZ = true
					break
				}
			}
			if hasLZ {
				if rule.IsCanBaiReplaceLZ() {
					tingInfo = append(tingInfo, COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI)
				}
			}
			tingInfo = append(tingInfo, rule.GetLaiziCard()...)
		}
		sort.Slice(tingInfo, func(i, j int) bool { return tingInfo[i] < tingInfo[j] })
		tingInfo = util.UniqueSlice(tingInfo).([]int32)
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
		if rule.IsLaiziCard(cards[i]) || i > 0 && cards[i] == cards[i-1] {
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
	maxBaseN, lzCard := int32(0), int32(0)
	if len(rule.lzCards) > 0 {
		lzCard = rule.lzCards[0]
	}
	for k := minIndex; k < len(indexs); k++ {
		bkc := cards[indexs[k]]
		cards[indexs[k]] = lzCard
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
			if cards[indexs[k]] != 0 && rule.IsLaiziCard(cards[indexs[k]]) == false &&
				indexs[k] < len(cards)-1 && cards[indexs[k]+1] == cards[indexs[k]] &&
				(indexs[k] > len(cards)-3 || cards[indexs[k]+1] != cards[indexs[k]+2]) {
				duiIndex = append(duiIndex, k)
			}
		}
		if len(duiIndex) >= 5 { //当对子数超过5个时，需要尝试拆对子来检查听
			for _, v := range duiIndex {
				bkc := cards[indexs[v]]
				cards[indexs[v]], cards[indexs[v]+1] = lzCard, lzCard
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
		tings = util.UniqueSlice(tings).([]int32)
	}
	return maxBaseN, tings
}
