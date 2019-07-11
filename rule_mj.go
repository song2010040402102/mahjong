package logic

import (
	"sort"
	"util"	
)

func (rule *RuleMahjong) CheckHuHua(index int32, huaCards []int32, flag uint32) uint64 {
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
	} else if rule.IsHas7HuaHu() && realHuaNum >= 7 {
		huType |= MJ_HU_TYPE_7HUA
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

func (rule *RuleMahjong) CheckHuType(index int32, chiCards []*ChiCard, huaCards []int32, holdCards []int32, card int32, flag uint32) (uint64, [][][3]int32) {
	var groups [][][3]int32
	if card != 0 && len(holdCards)%3 != 1 || card == 0 && len(holdCards)%3 != 2 {
		return MJ_HU_TYPE_NONE, groups
	}
	if flag&MJ_CHECK_HU_FLAG_MOBAO != 0 {
		return MJ_HU_TYPE_COMMON, groups
	}
	if flag&MJ_CHECK_HU_FLAG_3BAO != 0 {
		return MJ_HU_TYPE_3BAO, groups
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

		if rule.IsHasSiXiHu() {
			for _, v := range aiCards {
				if v.Num >= 4 {
					huType |= MJ_HU_TYPE_DASIXI
					break
				}
			}
		}
		if rule.IsHasBanBanHu() && CheckBanBanHu(chiCards, aiCards) {
			huType |= MJ_HU_TYPE_BANBAN
		}
		if rule.IsHasLiuLiuShunHu() && CheckLiuLiuShun(chiCards, aiCards) {
			huType |= MJ_HU_TYPE_LIULIUSHUN
		}
		if rule.IsHasQueYiSeHu() && CheckQueYiMen(chiCards, aiCards) {
			huType |= MJ_HU_TYPE_QUEYISE
		}
	}

	if len(chiCards) == 0 {
		if rule.IsHasQidui() {
			if ok, group := CheckQiDuiHu(aiCards, lzNum, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
				huType |= MJ_HU_TYPE_QIDUI
				groups = append(groups, group)
			} else if rule.MjType == RULE_ZJ_MAHJONG_WENZHOU {
				if CheckSoftNDuiHu(aiCards, lzNum, 8) {
					huType |= MJ_HU_TYPE_QIDUI
				}
			}
		}
		if lzNum <= 0 || lzNum == 1 && rule.IsLaiziCard(card) && flag&MJ_CHECK_HU_FLAG_ZIMO != 0 {
			if rule.IsHas7Star13BuKao() {
				if ok, group := Check7Star13BuKaoHu(aiCards, lzNum, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
					huType |= MJ_HU_TYPE_7STAR13BUKAO
					groups = append(groups, group)
				}
			}
			if huType&MJ_HU_TYPE_7STAR13BUKAO == 0 && rule.IsHas13BuKao() {
				if ok, group := Check13BuKaoHu(aiCards, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
					huType |= MJ_HU_TYPE_13BUKAO
					groups = append(groups, group)
				}
			}
		}
		if lzNum == 0 {
			if rule.IsHas7StarBuKao() {
				if ok, group := Check7StarBuKaoHu(aiCards, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
					huType |= MJ_HU_TYPE_7STARBUKAO
					groups = append(groups, group)
				}
			}
			if huType&MJ_HU_TYPE_7STARBUKAO == 0 && rule.IsHasQuanBuKao() {
				if ok, group := CheckQuanBuKaoHu(aiCards, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
					huType |= MJ_HU_TYPE_QUANBUKAO
					groups = append(groups, group)
				}
			}
			if rule.IsHas13Yao() {
				if ok, group := Check13YaoHu(aiCards, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
					huType |= MJ_HU_TYPE_13YAO
					groups = append(groups, group)
				}
			}
		}
	}
	if lzNum == 0 {
		if rule.IsHasZuHeLong() {
			if ok, group := CheckZuHeLongHu(chiCards, aiCards, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
				huType |= MJ_HU_TYPE_ZUHELONG
				groups = append(groups, group...)
			}
		}
		if rule.IsHasLuanSanFeng() {
			if ok, group := CheckLuanSanFengHu(aiCards, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
				huType |= MJ_HU_TYPE_LUANSANFENG
				groups = append(groups, group...)
			}
		}
	}
	if rule.IsExtraQingYiSeHu() && CheckQingYiSe(chiCards, aiCards) {
		huType |= MJ_HU_TYPE_QINGYISE
	} else if rule.IsExtraZiYiSeHu() && CheckZiYiSe(chiCards, aiCards) {
		huType |= MJ_HU_TYPE_ZIYISE
	} else if rule.IsHasJiangJiangHu() && CheckJiangJiangHu(chiCards, aiCards) {
		huType |= MJ_HU_TYPE_JIANGJIANG
	}
	if huType&(MJ_HU_TYPE_QINGYISE|MJ_HU_TYPE_ZIYISE|MJ_HU_TYPE_JIANGJIANG) != 0 {
		group := make([][3]int32, len(holdCards))
		for _, v := range holdCards {
			group = append(group, [3]int32{v})
		}
		group = append(group, [3]int32{card})
		groups = append(groups, group)
	}

	//后检查普通胡法
	if c := rule.GetCardForReplaceLZ(); c > 0 {
		ReplaceCardWithLZ([]*ChiCard{}, aiCards, rule.lzCards, c)
	}
	if ok, group := CheckCommonHu(aiCards, lzNum, flag&MJ_CHECK_HU_FLAG_GROUP != 0); ok {
		huType |= MJ_HU_TYPE_COMMON
		groups = append(groups, group...)
	}
	return huType, groups
}

func (rule *RuleMahjong) CheckTing(chiCards []*ChiCard, holdCards []int32, all bool) []int32 {
	tingInfo := make([]int32, 0, 34)

	if len(holdCards)%3 != 1 {
		return tingInfo
	}

	aiCards, lzNum := ConvertSliceToAICard(holdCards, 0, rule.lzCards)
	if lzNum <= 0 {
		if len(chiCards) == 0 {
			if rule.IsHasQuanBuKao() {
				tingInfo = append(tingInfo, CheckQuanBuKaoTing(aiCards)...)
			} else if rule.IsHas7StarBuKao() {
				tingInfo = append(tingInfo, Check7StarBuKaoTing(aiCards)...)
			}
			if rule.IsHas13BuKao() {
				tingInfo = append(tingInfo, Check13BuKaoTing(aiCards)...)
			} else if rule.IsHas7Star13BuKao() {
				tingInfo = append(tingInfo, Check7Star13BuKaoTing(aiCards)...)
			}
			if rule.IsHas13Yao() {
				tingInfo = append(tingInfo, Check13YaoTing(aiCards)...)
			}
		}
		if rule.IsHasZuHeLong() {
			tingInfo = append(tingInfo, CheckZuHeLongTing(chiCards, aiCards)...)
		}
	}
	if rule.IsExtraQingYiSeHu() {
		tingInfo = append(tingInfo, CheckQingYiSeTing(chiCards, aiCards)...)
	}
	if rule.IsExtraZiYiSeHu() {
		tingInfo = append(tingInfo, CheckZiYiSeTing(chiCards, aiCards)...)
	}
	if rule.IsHasJiangJiangHu() {
		tingInfo = append(tingInfo, CheckJiangJiangTing(chiCards, aiCards)...)
	}
	if !all && len(tingInfo) > 0 {
		return tingInfo
	}
	tingInfo = append(tingInfo, rule.doCheckTing(chiCards, holdCards, all)...)
	if len(tingInfo) > 0 {
		if tingInfo[0] == MAHJONG_ANY {
			return tingInfo[:1]
		}
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
				if card := rule.GetCardForReplaceLZ(); card > 0 {
					tingInfo = append(tingInfo, card)
				}
			}
			tingInfo = append(tingInfo, rule.GetLaiziCard()...)
		}
		sort.Slice(tingInfo, func(i, j int) bool { return tingInfo[i] < tingInfo[j] })
		tingInfo = util.UniqueSlice(tingInfo).([]int32)
		if all {
			tingInfo = rule.filterTingInfo(tingInfo)
		}
	}
	return tingInfo
}

//这里仅检测大概率出现的听牌牌型，例如七对和普通听，若检查类似13幺、全不靠、4癞子胡等极难出现的听牌，会使机器人打牌误入歧途
func (rule *RuleMahjong) doCheckTing(chiCards []*ChiCard, holdCards []int32, all bool) []int32 {
	tingInfo := make([]int32, 0, 34)

	aiCards, lzNum := ConvertSliceToAICard(holdCards, 0, rule.lzCards)

	if rule.IsHasQidui() && len(chiCards) <= 0 {
		tingInfo = append(tingInfo, CheckQiDuiTing(aiCards, lzNum)...)
	}
	if len(tingInfo) > 0 && (!all || tingInfo[0] == MAHJONG_ANY) {
		return tingInfo
	}

	if card := rule.GetCardForReplaceLZ(); card > 0 && len(rule.lzCards) > 0 && rule.lzCards[0]%MAHJONG_MASK >= MAHJONG_1 && rule.lzCards[0]%MAHJONG_MASK <= MAHJONG_9 {
		//若白板可以替癞子牌，应提前换掉
		for i := 0; i < len(aiCards); i++ {
			if aiCards[i].Card == card {
				aiCards[i].Card = rule.lzCards[0]
				sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })
				break
			}
		}
	}
	tingInfo = append(tingInfo, CheckCommonTing(aiCards, lzNum, all)...)
	if all {
		tingInfo = rule.filterTingInfo(tingInfo)
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
	tings := rule.doCheckTing(chiCards, cards, all)
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
		if baseN > 0 { //加快向听搜索
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

//对于麻将牌中缺少序数牌的情况，听牌信息需统一过滤一下，若缺少2到8的序数牌还需要考虑胡牌及是否听的检查
func (rule *RuleMahjong) filterTingInfo(tingInfo []int32) []int32 {
	cards := rule.cardDeck.GetCardKind(false)
	return util.RemoveSliceElem2(tingInfo, func(i int) bool {
		if tingInfo[i] == MAHJONG_ANY {
			return false
		}
		for _, v := range cards {
			if v == tingInfo[i] {
				return false
			}
		}
		return true
	}, true).([]int32)
}
