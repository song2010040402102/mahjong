package logic

import (
	"sort"
	"util"	
)

const (
	GT_NONE   int32 = iota
	GT_KE           //刻子
	GT_DUI          //对子
	GT_SEQ          //顺子
	GT_TWO          //二连，相邻两牌间隔为1或2
	GT_QIA          //卡子，隔张牌
	GT_DAN          //单张
	GT_3LZ          //三癞子刻
	GT_JG           //无癞子将
	GT_JG_LZ        //单癞子将
	GT_JG_2LZ       //两癞子将
)

type AICard struct {
	Card int32
	Num  int32
}

type TreeCard struct {
	gtype int32
	cards [3]int32
	child []*TreeCard
}

func NewTreeCard() *TreeCard {
	tree := &TreeCard{
		gtype: GT_NONE,
	}
	return tree
}

func createHuGroup(treeCard *TreeCard) [][][3]int32 {
	groups := [][][3]int32{}
	for _, v := range treeCard.child {
		subG := createHuGroup(v)
		for i := 0; i < len(subG); i++ {
			subG[i] = append(subG[i], treeCard.cards)
		}
		groups = append(groups, subG...)
	}
	if len(treeCard.child) == 0 {
		groups = append(groups, [][3]int32{treeCard.cards})
	}
	return groups
}

func analyzeCommonHu(aiCards []AICard, lzNum int32, all bool, treeCard *TreeCard, index int) bool {
	for {
		if index >= len(aiCards) || aiCards[index].Num > 0 {
			break
		}
		index++
	}
	if index >= len(aiCards) {
		if all {
			if lzNum == 3 {
				tree := NewTreeCard()
				tree.gtype = GT_3LZ
				tree.cards = [3]int32{MAHJONG_LZ, MAHJONG_LZ, MAHJONG_LZ}
				treeCard.child = append(treeCard.child, tree)
			}
		}
		return true
	}

	if aiCards[index].Num >= 3 { //不需要癞子的刻子
		tree := NewTreeCard()
		tree.gtype = GT_KE
		if all {
			tree.cards = [3]int32{aiCards[index].Card, aiCards[index].Card, aiCards[index].Card}
		}
		treeCard.child = append(treeCard.child, tree)
	}
	if aiCards[index].Num >= 2 && lzNum > 0 { //需要一个癞子的对子
		tree := NewTreeCard()
		tree.gtype = GT_DUI
		if all {
			tree.cards = [3]int32{aiCards[index].Card, aiCards[index].Card, MAHJONG_LZ}
		}
		treeCard.child = append(treeCard.child, tree)
	}
	if c := aiCards[index].Card % MAHJONG_MASK; c >= MAHJONG_1 && c <= MAHJONG_9 {
		if index <= len(aiCards)-3 && aiCards[index+1].Num > 0 && aiCards[index+2].Num > 0 &&
			aiCards[index].Card == aiCards[index+1].Card-1 && aiCards[index+1].Card+1 == aiCards[index+2].Card { //不需癞子的顺子
			tree := NewTreeCard()
			tree.gtype = GT_SEQ
			if all {
				tree.cards = [3]int32{aiCards[index].Card, aiCards[index+1].Card, aiCards[index+2].Card}
			}
			treeCard.child = append(treeCard.child, tree)
		}
		if lzNum > 0 { //需要一个癞子的顺子
			tree := NewTreeCard()
			if index <= len(aiCards)-3 && aiCards[index+1].Num == 0 && aiCards[index+2].Card-aiCards[index].Card <= 2 {
				tree.gtype = GT_QIA
				if all {
					tree.cards = [3]int32{aiCards[index].Card, aiCards[index+2].Card, MAHJONG_LZ}
				}
			} else if index <= len(aiCards)-2 && aiCards[index+1].Num != 0 && aiCards[index+1].Card-aiCards[index].Card <= 2 {
				tree.gtype = GT_TWO
				if all {
					tree.cards = [3]int32{aiCards[index].Card, aiCards[index+1].Card, MAHJONG_LZ}
				}
			}
			if tree.gtype != GT_NONE {
				treeCard.child = append(treeCard.child, tree)
			}
		}
	}
	if c := aiCards[index].Card % MAHJONG_MASK; lzNum > 1 && (aiCards[index].Num == 1 || c >= MAHJONG_1 && c <= MAHJONG_9) { //需要两个癞子的单张
		tree := NewTreeCard()
		tree.gtype = GT_DAN
		if all {
			tree.cards = [3]int32{aiCards[index].Card, MAHJONG_LZ, MAHJONG_LZ}
		}
		treeCard.child = append(treeCard.child, tree)
	}
	if len(treeCard.child) == 0 {
		return false
	}

	aiCardsBk := make([]AICard, len(aiCards))
	copy(aiCardsBk, aiCards)
	lzNumBk := lzNum
	for i := 0; i < len(treeCard.child); i++ {
		bHu := false
		if treeCard.child[i] != nil && treeCard.child[i].gtype >= GT_KE && treeCard.child[i].gtype <= GT_DAN {
			if treeCard.child[i].gtype == GT_KE {
				aiCards[index].Num -= 3
			} else if treeCard.child[i].gtype == GT_DUI {
				aiCards[index].Num -= 2
				lzNum--
			} else if treeCard.child[i].gtype == GT_SEQ {
				aiCards[index].Num--
				aiCards[index+1].Num--
				aiCards[index+2].Num--
			} else if treeCard.child[i].gtype == GT_TWO {
				aiCards[index].Num--
				aiCards[index+1].Num--
				lzNum--
			} else if treeCard.child[i].gtype == GT_QIA {
				aiCards[index].Num--
				aiCards[index+2].Num--
				lzNum--
			} else {
				aiCards[index].Num--
				lzNum -= 2
			}
			bHu = analyzeCommonHu(aiCards, lzNum, all, treeCard.child[i], index)
			copy(aiCards, aiCardsBk)
			lzNum = lzNumBk
		}
		if bHu {
			if !all {
				return true
			}
		} else {
			treeCard.child = append(treeCard.child[:i], treeCard.child[i+1:]...)
			i--
		}
	}
	return len(treeCard.child) > 0
}

//检查一般性胡法，可能存在两个等价的group，但顺序不同，考虑扩展性，不便于采用hash唯一化
func CheckCommonHu(cards []AICard, lzNum int32, all bool) (bool, [][][3]int32) {
	bHu := false
	var groups [][][3]int32
	if lzNum < 0 {
		return bHu, groups
	}
	if all {
		groups = make([][][3]int32, 0, len(cards))
	}

	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	sort.Slice(aiCards, func(i, index int) bool { return aiCards[i].Card < aiCards[index].Card })

	aiCardsBk := make([]AICard, len(aiCards))
	copy(aiCardsBk, aiCards)
	lzNumBk := lzNum

	for i := -1; i < len(aiCards); i++ {
		if i == -1 {
			if lzNum < 2 {
				continue
			}
		} else if aiCards[i].Num == 1 && lzNum < 1 {
			continue
		}
		treeCard := NewTreeCard()
		if all {
			if i == -1 {
				treeCard.gtype = GT_JG_2LZ
				treeCard.cards = [3]int32{MAHJONG_LZ, MAHJONG_LZ, 0}
			} else if aiCards[i].Num == 1 {
				treeCard.gtype = GT_JG_LZ
				treeCard.cards = [3]int32{aiCards[i].Card, MAHJONG_LZ, 0}
			} else {
				treeCard.gtype = GT_JG
				treeCard.cards = [3]int32{aiCards[i].Card, aiCards[i].Card, 0}
			}
		}
		if i == -1 {
			lzNum -= 2
		} else if aiCards[i].Num == 1 {
			aiCards[i].Num--
			lzNum--
		} else {
			aiCards[i].Num -= 2
		}
		if analyzeCommonHu(aiCards, lzNum, all, treeCard, 0) {
			bHu = true
			if all {
				groups = append(groups, createHuGroup(treeCard)...)
			} else {
				break
			}
		} else {
			treeCard = nil
		}
		copy(aiCards, aiCardsBk)
		lzNum = lzNumBk
	}
	return bHu, groups
}

//检查一般性听
func CheckCommonTing(cards []AICard, lzNum int32, all bool) []int32 {
	tingInfo := make([]int32, 0, 34)
	if lzNum < 0 {
		return tingInfo
	}

	//检查是否听所有牌
	huCards := make([]AICard, len(cards)+1)
	copy(huCards, cards)
	if ok, _ := CheckCommonHu(append(huCards, AICard{MAHJONG_ANY, 1}), lzNum, false); ok {
		tingInfo = append(tingInfo, MAHJONG_ANY)
		return tingInfo
	}

	if ok, groups := CheckCommonHu(cards, lzNum+1, all); ok {
		if all {
			for _, group := range groups {
				for _, v := range group {
					if v[2] == 0 {
						if v[1] == MAHJONG_LZ {
							tingInfo = append(tingInfo, v[0])
						}
					} else if v[2] == MAHJONG_LZ {
						if v[1] == MAHJONG_LZ {
							if v[0]%MAHJONG_MASK >= MAHJONG_DONG {
								tingInfo = append(tingInfo, v[0])
							} else if v[0]%MAHJONG_MASK == MAHJONG_1 {
								tingInfo = append(tingInfo, []int32{v[0], v[0] + 1, v[0] + 2}...)
							} else if v[0]%MAHJONG_MASK == MAHJONG_2 {
								tingInfo = append(tingInfo, []int32{v[0] - 1, v[0], v[0] + 1, v[0] + 2}...)
							} else if v[0]%MAHJONG_MASK == MAHJONG_8 {
								tingInfo = append(tingInfo, []int32{v[0] - 2, v[0] - 1, v[0], v[0] + 1}...)
							} else if v[0]%MAHJONG_MASK == MAHJONG_9 {
								tingInfo = append(tingInfo, []int32{v[0] - 2, v[0] - 1, v[0]}...)
							} else {
								tingInfo = append(tingInfo, []int32{v[0] - 2, v[0] - 1, v[0], v[0] + 1, v[0] + 2}...)
							}
						} else if v[0] == v[1] {
							tingInfo = append(tingInfo, v[0])
						} else if v[0]+1 == v[1] {
							if v[0]%MAHJONG_MASK == MAHJONG_1 {
								tingInfo = append(tingInfo, v[1]+1)
							} else if v[1]%MAHJONG_MASK == MAHJONG_9 {
								tingInfo = append(tingInfo, v[0]-1)
							} else {
								tingInfo = append(tingInfo, v[0]-1)
								tingInfo = append(tingInfo, v[1]+1)
							}
						} else {
							tingInfo = append(tingInfo, v[0]+1)
						}
					}
				}
			}
		} else {
			tingInfo = append(tingInfo, int32(0))
		}
	}

	if len(tingInfo) > 0 {
		sort.Slice(tingInfo, func(i, index int) bool { return tingInfo[i] < tingInfo[index] })
		tingInfo = util.UniqueSlice(tingInfo).([]int32)
	}
	return tingInfo
}

//检查七对胡
func CheckQiDuiHu(cards []AICard, lzNum int32, all bool) (bool, [][3]int32) {
	groups := [][3]int32{}
	if lzNum < 0 {
		return false, groups
	}
	for _, v := range cards {
		if v.Num%2 == 1 {
			lzNum--
			if lzNum < 0 {
				return false, groups
			}
			if all {
				groups = append(groups, [3]int32{v.Card, MAHJONG_LZ, 0})
			}
		}
		if all {
			for i := int32(0); i < v.Num/2; i++ {
				groups = append(groups, [3]int32{v.Card, v.Card, 0})
			}
		}
	}
	return true, groups
}

//检查七对听
func CheckQiDuiTing(cards []AICard, lzNum int32) []int32 {
	if lzNum < 0 {
		return []int32{}
	}
	tingInfo := make([]int32, 0, 5)
	for i := 0; i < len(cards); i++ {
		if cards[i].Num%2 != 0 {
			if lzNum >= 0 {
				lzNum--
				tingInfo = append(tingInfo, cards[i].Card)
			} else {
				return tingInfo[:0]
			}
		}
	}
	if lzNum > 0 {
		tingInfo = append(tingInfo[:0], MAHJONG_ANY)
	} else {
		sort.Slice(tingInfo, func(i, index int) bool { return tingInfo[i] < tingInfo[index] })
	}
	return tingInfo
}

//检查十三不靠胡，不能含癞子
func Check13BuKaoHu(cards []AICard, all bool) (bool, [][3]int32) {
	groups := [][3]int32{}
	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	sort.Slice(aiCards, func(i, index int) bool { return aiCards[i].Card < aiCards[index].Card })
	for i := 0; i < len(aiCards); i++ {
		if aiCards[i].Num > 1 || (aiCards[i].Card%MAHJONG_MASK < MAHJONG_DONG && i < len(aiCards)-1 && aiCards[i+1].Num > 0 && aiCards[i+1].Card-aiCards[i].Card <= 2) {
			return false, groups
		}
		if all {
			groups = append(groups, [3]int32{aiCards[i].Card})
		}
	}
	return true, groups
}

//检查7星13不靠胡
func Check7Star13BuKaoHu(cards []AICard, lzNum int32, all bool) (bool, [][3]int32) {
	if lzNum <= 1 {
		ok, groups := Check13BuKaoHu(cards, all)
		if ok {
			mask := uint64(0)
			for _, v := range cards {
				mask |= uint64(0x01) << cardId2Bit(v.Card)
			}
			if mask&(uint64(0x7f)<<27) == uint64(0x7f)<<27 {
				return ok, groups
			}
			if lzNum == 1 {
				num := 0
				for i := uint32(27); i < 34; i++ {
					if mask&(uint64(0x01)<<i) > 0 {
						num++
					}
				}
				if num == 6 {
					return ok, groups
				}
			}
		}
	}
	return false, [][3]int32{}
}

//检查十三不靠听，不能含癞子
func Check13BuKaoTing(cards []AICard) []int32 {
	tingInfo := make([]int32, 0, 34)
	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	sort.Slice(aiCards, func(i, index int) bool { return aiCards[i].Card < aiCards[index].Card })

	//先判断是否听
	for i := 0; i < len(aiCards); i++ {
		if aiCards[i].Num > 1 || (aiCards[i].Card%MAHJONG_MASK < MAHJONG_DONG && i < len(aiCards)-1 && aiCards[i+1].Num > 0 && aiCards[i+1].Card-aiCards[i].Card <= 2) {
			return tingInfo
		}
	}

	//处理左边界的听牌
	for i := COLOR_WAN*MAHJONG_MASK + MAHJONG_1; i < aiCards[0].Card-2; {
		tingInfo = append(tingInfo, i)
		if i%MAHJONG_MASK == MAHJONG_9 {
			i += MAHJONG_MASK - MAHJONG_9 + MAHJONG_1
		} else {
			i++
		}
	}

	//追加右边界，处理中间的听牌
	aiCards = append(aiCards, AICard{COLOR_OTHER*MAHJONG_MASK + MAHJONG_BAI + 1, 1})
	for i := 0; i < len(aiCards)-1; i++ {
		lb := aiCards[i].Card + 1
		if aiCards[i].Card%MAHJONG_MASK < MAHJONG_DONG {
			lb += 2
			if lb%MAHJONG_MASK > MAHJONG_9 {
				if lb/MAHJONG_MASK == COLOR_TIAO {
					lb = COLOR_OTHER*MAHJONG_MASK + MAHJONG_DONG
				} else {
					lb = (lb/MAHJONG_MASK+1)*MAHJONG_MASK + MAHJONG_1
				}
			}
		}
		rb := aiCards[i+1].Card
		if aiCards[i+1].Card%MAHJONG_MASK < MAHJONG_DONG {
			rb -= 2
		}
		for index := lb; index < rb; {
			tingInfo = append(tingInfo, index)
			if index%MAHJONG_MASK == MAHJONG_9 {
				if index/MAHJONG_MASK == COLOR_TIAO {
					index = COLOR_OTHER*MAHJONG_MASK + MAHJONG_DONG
				} else {
					index += MAHJONG_MASK - MAHJONG_9 + MAHJONG_1
				}
			} else {
				index++
			}
		}
	}
	return tingInfo
}

func Check7Star13BuKaoTing(cards []AICard) []int32 {
	tings := Check13BuKaoTing(cards)
	if len(tings) > 0 {
		zi := []int32{}
		for _, v := range tings {
			if v%MAHJONG_MASK >= MAHJONG_DONG {
				zi = append(zi, v)
			}
		}
		if len(zi) == 1 {
			return zi
		} else if len(zi) > 1 {
			return []int32{}
		}
	}
	return tings
}

func cardId2Bit(card int32) uint32 {
	if card%MAHJONG_MASK >= MAHJONG_DONG {
		return uint32((card/MAHJONG_MASK-1)*9 + card%MAHJONG_MASK - MAHJONG_DONG)
	}
	return uint32((card/MAHJONG_MASK-1)*9 + card%MAHJONG_MASK - MAHJONG_1)
}

func cardBit2Id(bit uint32) int32 {
	if bit/9 > 2 {
		return int32(bit/9+1)*MAHJONG_MASK + int32(bit%9) + MAHJONG_DONG
	}
	return int32(bit/9+1)*MAHJONG_MASK + int32(bit%9) + MAHJONG_1
}

func checkQuanBuKao(cards []AICard) []int32 {
	// 1111111 100100100 010010010 001001001
	// 1111111 100100100 001001001 010010010
	// 1111111 010010010 100100100 001001001
	// 1111111 010010010 001001001 100100100
	// 1111111 001001001 100100100 010010010
	// 1111111 001001001 010010010 100100100
	group := []uint64{17122272329, 17122235026, 17084074057, 17083962148, 17064937618, 17064863012}
	cg := uint64(0)
	for _, v := range cards {
		if v.Num > 1 {
			return []int32{}
		}
		cg |= uint64(0x01) << cardId2Bit(v.Card)
	}
	for _, v := range group {
		v ^= cg
		rest := []int32{}
		for i := uint32(0); i < 34; i++ {
			if v&(0x01<<i) != 0 {
				rest = append(rest, cardBit2Id(i))
			}
		}
		if len(rest) == 2 || len(rest) == 3 {
			return rest
		}
	}
	return []int32{}
}

//检查全不靠胡，不能含癞子
func CheckQuanBuKaoHu(cards []AICard, all bool) (bool, [][3]int32) {
	groups := [][3]int32{}
	if res := checkQuanBuKao(cards); len(res) == 2 {
		if all {
			for _, v := range cards {
				groups = append(groups, [3]int32{v.Card})
			}
		}
		return true, groups
	}
	return false, groups
}

//检查全不靠听，不能含癞子
func CheckQuanBuKaoTing(cards []AICard) []int32 {
	if res := checkQuanBuKao(cards); len(res) == 3 {
		return res
	}
	return []int32{}
}

//检查七星不靠胡，不能含癞子
func Check7StarBuKaoHu(cards []AICard, all bool) (bool, [][3]int32) {
	groups := [][3]int32{}
	res := checkQuanBuKao(cards)
	if len(res) == 2 && res[0]%MAHJONG_MASK < MAHJONG_DONG && res[1]%MAHJONG_MASK < MAHJONG_DONG {
		if all {
			for _, v := range cards {
				groups = append(groups, [3]int32{v.Card})
			}
		}
		return true, groups
	}
	return false, groups
}

//检查七星不靠听，不能含癞子
func Check7StarBuKaoTing(cards []AICard) []int32 {
	res := checkQuanBuKao(cards)
	if len(res) == 3 && res[1]%MAHJONG_MASK < MAHJONG_DONG {
		if res[2]%MAHJONG_MASK >= MAHJONG_DONG {
			return res[2:]
		} else {
			return res
		}
	}
	return []int32{}
}

//检查13幺胡，不能含癞子
func Check13YaoHu(cards []AICard, all bool) (bool, [][3]int32) {
	groups := [][3]int32{}
	my := uint64(17113154305) //1111111 100000001 100000001 100000001
	if len(cards) != 13 {
		return false, groups
	}
	cg := uint64(0)
	for _, v := range cards {
		cg |= uint64(0x01) << cardId2Bit(v.Card)
		if all {
			if v.Num == 2 {
				groups = append(groups, [3]int32{v.Card, v.Card})
			} else {
				groups = append(groups, [3]int32{v.Card})
			}
		}
	}
	return my^cg == 0, groups
}

//检查13幺听，不能含癞子
func Check13YaoTing(cards []AICard) []int32 {
	my := uint64(17113154305) //1111111 100000001 100000001 100000001
	cg := uint64(0)
	for _, v := range cards {
		cg |= uint64(0x01) << cardId2Bit(v.Card)
	}
	my ^= cg
	if len(cards) == 12 {
		rest := []int32{}
		for i := uint32(0); i < 34; i++ {
			if my&(0x01<<i) != 0 {
				rest = append(rest, cardBit2Id(i))
			}
		}
		if len(rest) == 1 {
			return rest
		}
	} else if len(cards) == 13 {
		if my == 0 {
			tings := []int32{}
			for _, v := range cards {
				tings = append(tings, v.Card)
			}
			return tings
		}
	}
	return []int32{}
}

func checkZuHeLong(cards []AICard) (int32, uint32, []AICard) {
	// 100100100 010010010 001001001
	// 100100100 001001001 010010010
	// 010010010 100100100 001001001
	// 010010010 001001001 100100100
	// 001001001 100100100 010010010
	// 001001001 010010010 100100100
	group := []uint32{76620873, 76583570, 38422601, 38310692, 19286162, 19211556}
	cg := uint32(0)
	for _, v := range cards {
		cg |= uint32(0x01) << cardId2Bit(v.Card)
	}

	lackC, mask := int32(-1), uint32(0)
	for _, v := range group {
		if v&cg == v {
			lackC, mask = 0, v
			break
		}
	}
	for _, v := range group {
		count := 0
		res := v&cg ^ v
		for i := uint32(0); i < 27; i++ {
			if res&(0x01<<i) != 0 {
				count++
				if count > 1 {
					lackC = -1
					break
				}
				lackC = cardBit2Id(i)
			}
		}
		if count == 1 {
			mask = v
			break
		}
	}
	if lackC == -1 || mask == 0 {
		return -1, mask, []AICard{}
	}

	restCards := make([]AICard, len(cards))
	copy(restCards, cards)
	for i := uint32(0); i < 27; i++ {
		if mask&(0x01<<i) != 0 {
			card := cardBit2Id(i)
			for index := 0; index < len(restCards); index++ {
				if restCards[index].Card == card {
					restCards[index].Num--
					if restCards[index].Num == 0 {
						restCards = append(restCards[:index], restCards[index+1:]...)
						index--
					}
				}
			}
		}
	}
	return lackC, mask, restCards
}

//检查不含癞子组合龙胡
func CheckZuHeLongHu(chiCards []*ChiCard, cards []AICard, all bool) (bool, [][][3]int32) {
	groups := [][][3]int32{}
	lackC, mask, restCards := checkZuHeLong(cards)
	if lackC == 0 {
		group := make([][3]int32, 3)
		if all {
			num := 0
			for i := uint32(0); i < 27; i++ {
				if mask&(uint32(0x01)<<i) != 0 {
					group[num/3][num%3] = cardBit2Id(i)
					num++
				}
			}
		}
		quanbukao := true
		for _, v := range restCards {
			if len(chiCards) > 0 || v.Num > 1 || v.Card%MAHJONG_MASK < MAHJONG_DONG {
				quanbukao = false
				break
			}
		}
		if quanbukao {
			if all {
				for _, v := range restCards {
					group = append(group, [3]int32{v.Card})
				}
				groups = append(groups, group)
			}
			return true, groups
		}
		ok, restGroups := CheckCommonHu(restCards, 0, true)
		if ok && all {
			for _, v := range restGroups {
				groups = append(groups, group)
				groups[len(groups)-1] = append(groups[len(groups)-1], v...)
			}
		}
		return ok, groups
	}
	return false, groups
}

//检查不含癞子组合龙听
func CheckZuHeLongTing(chiCards []*ChiCard, cards []AICard) []int32 {
	tings := []int32{}
	lackC, _, restCards := checkZuHeLong(cards)
	if lackC == 0 {
		if len(chiCards) == 0 {
			cg := uint64(0)
			for _, v := range restCards {
				cg |= uint64(0x01) << cardId2Bit(v.Card)
			}
			cg &= uint64(0x7f) << 27
			for i := uint32(27); i < 34; i++ {
				if cg&(0x01<<i) == 0 {
					tings = append(tings, cardBit2Id(i))
				}
			}
			if len(tings) != 3 {
				tings = tings[:0]
			}
		}
		if len(tings) == 0 {
			tings = append(tings, CheckCommonTing(restCards, 0, true)...)
		}
	} else if lackC > 0 {
		quanbukao := true
		for _, v := range restCards {
			if len(chiCards) > 0 || v.Num > 1 || v.Card%MAHJONG_MASK < MAHJONG_DONG {
				quanbukao = false
				break
			}
		}
		if quanbukao {
			tings = append(tings, lackC)
		} else if ok, _ := CheckCommonHu(restCards, 0, false); ok {
			tings = append(tings, lackC)
		}
	}
	return tings
}

func CheckHardNDuiHu(cards []AICard, lzNum int32, n int32) bool {
	num := int32(0)
	for i := 0; i < len(cards); i++ {
		num += cards[i].Num / 2
	}
	num += (lzNum / 2)

	return num >= n
}

func CheckSoftNDuiHu(cards []AICard, lzNum int32, n int32) bool {
	num := int32(0)
	for i := 0; i < len(cards); i++ {
		num += cards[i].Num / 2
	}
	num += lzNum

	return num >= n // 包含硬n对
}

/******************************以下为牌型判断**********************************/

//合并吃牌和手牌
func MergeChiCard(chiCards []*ChiCard, cards []AICard, gang4 bool) []AICard {
	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_CHI {
			c1, c2, c3 := int32(0), int32(0), int32(0)
			if v.ChiPosBit&(0x01<<1) != 0 {
				c1, c2, c3 = v.CardId-2, v.CardId-1, v.CardId
			} else if v.ChiPosBit&(0x01<<2) != 0 {
				c1, c2, c3 = v.CardId-1, v.CardId, v.CardId+1
			} else if v.ChiPosBit&(0x01<<3) != 0 {
				c1, c2, c3 = v.CardId, v.CardId+1, v.CardId+2
			}
			aiCards = append(aiCards, []AICard{{c1, 1}, {c2, 1}, {c3, 1}}...)
		} else if v.CardType == MJ_CHI_PENG {
			aiCards = append(aiCards, AICard{v.CardId, 3})
		} else {
			if gang4 {
				aiCards = append(aiCards, AICard{v.CardId, 4})
			} else {
				aiCards = append(aiCards, AICard{v.CardId, 3})
			}
		}
	}
	return aiCards
}

//检查单调，不适用于七对，因为七对的胡牌只能单调
func CheckTingTypeForDanDiao(cards []AICard) bool {
	if len(cards) == 1 && cards[0].Num == 1 {
		return true
	}
	tings := CheckCommonTing(cards, 0, true)
	if len(tings) != 1 {
		return false
	}

	//思路：去除单调的那一张牌，然后加上任意一对将，若胡，则为单调
	found := false
	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	for i := 0; i < len(aiCards); i++ {
		if aiCards[i].Card == tings[0] {
			aiCards[i].Num--
			found = true
			break
		}
	}
	if !found {
		return false
	}
	aiCards = append(aiCards, AICard{MAHJONG_ANY, 2})
	if ok, _ := CheckCommonHu(aiCards, 0, false); ok {
		return true
	}
	return false
}

//检查卡子，例：1万3万胡2万
func CheckTingTypeForQiaZi(cards []AICard) bool {
	tings := CheckCommonTing(cards, 0, true)
	if len(tings) != 1 || tings[0]%MAHJONG_MASK < MAHJONG_2 || tings[0]%MAHJONG_MASK > MAHJONG_8 {
		return false
	}

	//思路：去除两边的牌，若胡，则听卡子
	found1, found2 := false, false
	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	for i := 0; i < len(aiCards); i++ {
		if found1 && found2 {
			break
		}
		if !found1 && aiCards[i].Card == tings[0]-1 {
			aiCards[i].Num--
			found1 = true
		}
		if !found2 && aiCards[i].Card == tings[0]+1 {
			aiCards[i].Num--
			found2 = true
		}
	}
	if !found1 || !found2 {
		return false
	}
	if ok, _ := CheckCommonHu(aiCards, 0, false); ok {
		return true
	}
	return false
}

//检查边张，例：1万2万胡3万
func CheckTingTypeForBianZhang(cards []AICard) bool {
	tings := CheckCommonTing(cards, 0, true)
	if len(tings) != 1 || (tings[0]%MAHJONG_MASK != MAHJONG_3 && tings[0]%MAHJONG_MASK != MAHJONG_7) {
		return false
	}

	//思路：去除边上的牌，若胡，则听边张
	found1, found2 := false, false
	c1, c2 := tings[0]+1, tings[0]+2
	if tings[0]%MAHJONG_MASK == MAHJONG_3 {
		c1, c2 = tings[0]-1, tings[0]-2
	}
	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	for i := 0; i < len(aiCards); i++ {
		if found1 && found2 {
			break
		}
		if !found1 && aiCards[i].Card == c1 {
			aiCards[i].Num--
			found1 = true
		}
		if !found2 && aiCards[i].Card == c2 {
			aiCards[i].Num--
			found2 = true
		}
	}
	if !found1 || !found2 {
		return false
	}
	if ok, _ := CheckCommonHu(aiCards, 0, false); ok {
		return true
	}
	return false
}

//检查对倒
func CheckTingTypeForDuiDao(cards []AICard) bool {
	tings := CheckCommonTing(cards, 0, true)
	if len(tings) != 2 {
		return false
	}

	//思路：去除听的两张对牌，加上任意一对将，若胡，则听边张
	found1, found2 := false, false
	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	for i := 0; i < len(aiCards); i++ {
		if found1 && found2 {
			break
		}
		if !found1 && aiCards[i].Card == tings[0] {
			aiCards[i].Num -= 2
			found1 = true
		}
		if !found2 && aiCards[i].Card == tings[1] {
			aiCards[i].Num -= 2
			found2 = true
		}
	}
	if !found1 || !found2 {
		return false
	}
	aiCards = append(aiCards, AICard{MAHJONG_ANY, 2})
	if ok, _ := CheckCommonHu(aiCards, 0, false); ok {
		return true
	}
	return false
}

//检查清一色
func CheckQingYiSe(chiCards []*ChiCard, cards []AICard) bool {
	color := int32(0)
	//先判断吃牌花色是否相同
	for _, v := range chiCards {
		if color == 0 {
			color = v.CardId / MAHJONG_MASK
			if color > COLOR_TIAO || color < COLOR_WAN {
				return false
			}
		} else if color != v.CardId/MAHJONG_MASK {
			return false
		}
	}
	//再判断手牌花色是否相同
	for i := 0; i < len(cards); i++ {
		if color == 0 {
			color = cards[i].Card / MAHJONG_MASK
			if color > COLOR_TIAO || color < COLOR_WAN {
				return false
			}
		} else if color != cards[i].Card/MAHJONG_MASK {
			return false
		}
	}
	return true
}

//检查清一色听
func CheckQingYiSeTing(chiCards []*ChiCard, cards []AICard) []int32 {
	tings := []int32{}
	if CheckQingYiSe(chiCards, cards) && len(cards) > 0 {
		for c := MAHJONG_1; c <= MAHJONG_9; c++ {
			tings = append(tings, cards[0].Card-cards[0].Card%MAHJONG_MASK+c)
		}
	}
	return tings
}

//检查字一色
func CheckZiYiSe(chiCards []*ChiCard, cards []AICard) bool {
	color := int32(0)
	//先判断吃牌花色是否相同
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_CHI {
			return false
		}
		if color == 0 {
			color = v.CardId / MAHJONG_MASK
			if color != COLOR_OTHER {
				return false
			}
		} else if color != v.CardId/MAHJONG_MASK {
			return false
		}
	}
	//再判断手牌花色是否相同
	for i := 0; i < len(cards); i++ {
		if color == 0 {
			color = cards[i].Card / MAHJONG_MASK
			if color != COLOR_OTHER {
				return false
			}
		} else if color != cards[i].Card/MAHJONG_MASK {
			return false
		}
	}
	return true
}

//检查字一色听
func CheckZiYiSeTing(chiCards []*ChiCard, cards []AICard) []int32 {
	tings := []int32{}
	if CheckZiYiSe(chiCards, cards) && len(cards) > 0 {
		for c := MAHJONG_DONG; c <= MAHJONG_BAI; c++ {
			tings = append(tings, cards[0].Card-cards[0].Card%MAHJONG_MASK+c)
		}
	}
	return tings
}

//检查混一色
func CheckHunYiSe(chiCards []*ChiCard, cards []AICard) bool {
	color1 := int32(0)
	color2 := int32(0)
	//先判断吃牌花色
	for _, v := range chiCards {
		c := v.CardId / MAHJONG_MASK
		if color1 == 0 {
			color1 = c
		} else if color2 == 0 && color1 != c {
			color2 = c
			if color1 != COLOR_OTHER && color2 != COLOR_OTHER {
				return false
			}
		} else if color1 != c && color2 != c {
			return false
		}
	}
	//再判断手牌花色
	for i := 0; i < len(cards); i++ {
		if color1 == 0 {
			color1 = cards[i].Card / MAHJONG_MASK
		} else if color2 == 0 && color1 != cards[i].Card/MAHJONG_MASK {
			color2 = cards[i].Card / MAHJONG_MASK
			if color1 != COLOR_OTHER && color2 != COLOR_OTHER {
				return false
			}
		} else if color1 != cards[i].Card/MAHJONG_MASK && color2 != cards[i].Card/MAHJONG_MASK {
			return false
		}
	}
	return color2 != 0
}

//检查绿一色
func CheckLvYiSe(chiCards []*ChiCard, cards []AICard) bool {
	mask := uint64(4340580352) //0100000 010101110 000000000 000000000, 由23468条及发财组成的胡牌
	aiCards := MergeChiCard(chiCards, cards, false)
	for _, v := range aiCards {
		if mask&(uint64(0x01)<<cardId2Bit(v.Card)) == 0 {
			return false
		}
	}
	return true
}

func getCardMinMaxMask(cards []AICard) (int32, int32, uint32) {
	min, max, mask := int32(1000), int32(0), uint32(0)
	for _, v := range cards {
		if v.Card < min {
			min = v.Card
		}
		if v.Card > max {
			max = v.Card
		}
		mask |= uint32(v.Num) << uint32(v.Card%MAHJONG_MASK-MAHJONG_1) * 3
	}
	return min, max, mask
}

//检查九莲宝灯, cards不包含胡那张牌
func Check9LianBaoDeng(cards []AICard) bool {
	min, max, mask := getCardMinMaxMask(cards)
	if max-min == MAHJONG_9-MAHJONG_1 && mask == 52728395 { //011 001 001 001 001 001 001 001 011
		return true
	}
	return false
}

//检查一色双龙会，相同花色，两个123、两个789、对5
func CheckYiSeSameLong(chiCards []*ChiCard, cards []AICard) bool {
	aiCards := MergeChiCard(chiCards, cards, false)
	min, max, mask := getCardMinMaxMask(aiCards)
	if max-min == MAHJONG_9-MAHJONG_1 && mask == 38281362 { //010 010 010 000 010 000 010 010 010
		return true
	}
	return false
}

//检查三色双龙会，两种花色老少副，另一种花色5做将
func Check3SeSameLong(chiCards []*ChiCard, cards []AICard) bool {
	mask := uint32(0)
	aiCards := MergeChiCard(chiCards, cards, false)
	for _, v := range aiCards {
		if v.Card%MAHJONG_MASK >= MAHJONG_DONG {
			return false
		}
		mask |= uint32(0x01) << cardId2Bit(v.Card)
	}

	// 111000111 111000111 000010000
	// 111000111 000010000 111000111
	// 000010000 111000111 111000111
	if mask == 119508496 || mask == 119284167 || mask == 4427719 {
		return true
	}
	return false
}

//检查对对胡
func CheckDuiDuiHu(chiCards []*ChiCard, cards []AICard) bool {
	for _, v := range chiCards {
		if v.CardType != MJ_CHI_GANG &&
			v.CardType != MJ_CHI_GANG_WAN &&
			v.CardType != MJ_CHI_GANG_AN &&
			v.CardType != MJ_CHI_PENG {
			return false
		}
	}

	for i, jiang := 0, false; i < len(cards); i++ {
		if cards[i].Num == 1 || cards[i].Num == 4 {
			return false
		} else if cards[i].Num == 2 {
			if jiang == false {
				jiang = true
			} else {
				return false
			}
		}
	}
	return true
}

//检查对对胡，带癞子牌
func CheckDuiDuiHu2(chiCards []*ChiCard, cards []AICard, lzNum int32) bool {
	if ok, groups := CheckCommonHu(cards, lzNum, true); ok {
		for _, v := range groups {
			if CheckDuiDuiHuForLZ(chiCards, v) {
				return true
			}
		}
	}
	return false
}

//检查清龙
func CheckQingLong2(chiCards []*ChiCard, cards []AICard, lzNum int32) bool {
	if ok, groups := CheckCommonHu(cards, lzNum, true); ok {
		for _, v := range groups {
			if CheckQingLong(chiCards, v) {
				return true
			}
		}
	}
	return false
}

//用于四川麻将的将对（对对胡并且每组都有258）
func CheckDai258(chiCards []*ChiCard, cards []AICard) bool {
	if !CheckDuiDuiHu(chiCards, cards) {
		return false
	}
	for _, v := range chiCards {
		if c := v.CardId % MAHJONG_MASK; c != MAHJONG_2 && c != MAHJONG_5 && c != MAHJONG_8 {
			return false
		}
	}
	for i := 0; i < len(cards); i++ {
		if c := cards[i].Card % MAHJONG_MASK; c != MAHJONG_2 && c != MAHJONG_5 && c != MAHJONG_8 {
			return false
		}
	}
	return true
}

//用于四川麻将的中张
func CheckZhongZhang(chiCards []*ChiCard, cards []AICard) bool {
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_CHI {
			card := v.CardId
			if v.ChiPosBit&(0x01<<1) != 0 {
				card = v.CardId - 2
			} else if v.ChiPosBit&(0x01<<2) != 0 {
				card = v.CardId - 1
			}
			if c := card % MAHJONG_MASK; c == MAHJONG_1 || c == MAHJONG_7 {
				return false
			}
		} else {
			if c := v.CardId % MAHJONG_MASK; c == MAHJONG_1 || c == MAHJONG_9 || c >= MAHJONG_DONG {
				return false
			}
		}
	}
	for i := 0; i < len(cards); i++ {
		if c := cards[i].Card % MAHJONG_MASK; c == MAHJONG_1 || c == MAHJONG_9 || c >= MAHJONG_DONG {
			return false
		}
	}
	return true
}

//用于四川麻将根
func GetLongNum(chiCards []*ChiCard, cards []AICard) int32 {
	cardNum := map[int32]int32{}
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_GANG || v.CardType == MJ_CHI_GANG_WAN || v.CardType == MJ_CHI_GANG_AN {
			cardNum[v.CardId] += 4
		} else if v.CardType == MJ_CHI_PENG {
			cardNum[v.CardId] += 3
		}
	}
	for i := 0; i < len(cards); i++ {
		cardNum[cards[i].Card] += cards[i].Num
	}
	longNum := int32(0)
	for _, v := range cardNum {
		if v == 4 {
			longNum++
		}
	}
	return longNum
}

//全老头，用于嘉兴麻将
func CheckQuanLaoTou(chiCards []*ChiCard, cards []AICard) bool {
	for _, v := range chiCards {
		if c := v.CardId % MAHJONG_MASK; c < MAHJONG_DONG || c > MAHJONG_BAI {
			return false
		}
	}
	for i := 0; i < len(cards); i++ {
		if c := cards[i].Card % MAHJONG_MASK; c < MAHJONG_DONG || c > MAHJONG_BAI {
			return false
		}
	}
	return true
}

//四花齐放，同时拥有春夏秋冬或梅兰竹菊
func CheckFourHua(huaCards []int32) bool {
	if len(huaCards) >= 4 && len(huaCards) < 8 {
		bitHua := uint8(0)
		for _, v := range huaCards {
			bitHua |= 0x01 << uint8(v-(COLOR_OTHER*MAHJONG_MASK+MAHJONG_SPRING))
		}
		if bitHua&0x0f == 0x0f || bitHua&0xf0 == 0xf0 {
			return true
		}
	}
	return false
}

//检查门清
func CheckMenQing(chiCards []*ChiCard) bool {
	for _, v := range chiCards {
		if v.CardType != MJ_CHI_GANG_AN {
			return false
		}
	}
	return true
}

//检查大门清
func CheckDaMenQing(chiCards []*ChiCard) bool {
	if len(chiCards) > 0 {
		return false
	}
	return true
}

//检查全求人
func CheckQuanQiuRen(chiCards []*ChiCard, cards []AICard, ziMo bool) bool {
	if ziMo {
		return false
	}
	if len(cards) > 1 {
		return false
	}
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_GANG_AN {
			return false
		}
	}
	return true
}

//检查不求人
func CheckBuQiuRen(chiCards []*ChiCard, ziMo bool) bool {
	return ziMo && CheckMenQing(chiCards)
}

//检查全双刻，全是2、4、6、8组成的刻、杠、对
func CheckQuanPairKe(chiCards []*ChiCard, cards []AICard) bool {
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_CHI || v.CardId%MAHJONG_MASK >= MAHJONG_DONG || (v.CardId%MAHJONG_MASK-MAHJONG_1+1)%2 != 0 {
			return false
		}
	}
	for _, v := range cards {
		if v.Card%MAHJONG_MASK >= MAHJONG_DONG || (v.Card%MAHJONG_MASK-MAHJONG_1+1)%2 != 0 {
			return false
		}
	}
	return true
}

func getRangeCardNum(chiCards []*ChiCard, cards []AICard, min int32, max int32) (bool, int32) {
	all, num := true, int32(0)
	aiCards := MergeChiCard(chiCards, cards, true)
	for _, v := range aiCards {
		if c := v.Card % MAHJONG_MASK; c >= min && (max == -1 || c <= max) {
			num += v.Num
		} else {
			all = false
		}
	}
	return all, num
}

//检查全大，全是789
func CheckQuanDa(chiCards []*ChiCard, cards []AICard) bool {
	all, _ := getRangeCardNum(chiCards, cards, MAHJONG_7, MAHJONG_9)
	return all
}

//检查全中，全是456
func CheckQuanZhong(chiCards []*ChiCard, cards []AICard) bool {
	all, _ := getRangeCardNum(chiCards, cards, MAHJONG_4, MAHJONG_6)
	return all
}

//检查全小，全是123
func CheckQuanXiao(chiCards []*ChiCard, cards []AICard) bool {
	all, _ := getRangeCardNum(chiCards, cards, MAHJONG_1, MAHJONG_3)
	return all
}

//检查大于5
func CheckDaYu5(chiCards []*ChiCard, cards []AICard) bool {
	all, _ := getRangeCardNum(chiCards, cards, MAHJONG_6, MAHJONG_9)
	return all
}

//检查小于5
func CheckXiaoYu5(chiCards []*ChiCard, cards []AICard) bool {
	all, _ := getRangeCardNum(chiCards, cards, MAHJONG_1, MAHJONG_4)
	return all
}

//检查大于4
func CheckDaYu4(chiCards []*ChiCard, cards []AICard) bool {
	all, _ := getRangeCardNum(chiCards, cards, MAHJONG_5, MAHJONG_9)
	return all
}

//检查小于6
func CheckXiaoYu6(chiCards []*ChiCard, cards []AICard) bool {
	all, _ := getRangeCardNum(chiCards, cards, MAHJONG_1, MAHJONG_5)
	return all
}

//检查大于4是否超过10个
func Check10DaYu4(chiCards []*ChiCard, cards []AICard) bool {
	all, num := getRangeCardNum(chiCards, cards, MAHJONG_5, MAHJONG_9)
	return all || num >= 10
}

//检查小于6是否超过10个
func Check10XiaoYu6(chiCards []*ChiCard, cards []AICard) bool {
	all, num := getRangeCardNum(chiCards, cards, MAHJONG_1, MAHJONG_5)
	return all || num >= 10
}

//检查三风刻杠
func Check3FengKeGang(chiCards []*ChiCard, cards []AICard) bool {
	num := 0
	aiCards := MergeChiCard(chiCards, cards, false)
	for _, v := range aiCards {
		if c := v.Card % MAHJONG_MASK; v.Num >= 3 && c >= MAHJONG_DONG && c <= MAHJONG_BEI {
			num++
		}
	}
	return num == 3
}

//获取箭刻杠数
func GetJianKeGangNum(chiCards []*ChiCard, cards []AICard) int32 {
	num := int32(0)
	aiCards := MergeChiCard(chiCards, cards, false)
	for _, v := range aiCards {
		if c := v.Card % MAHJONG_MASK; v.Num >= 3 && c >= MAHJONG_HONGZHONG {
			num++
		}
	}
	return num
}

//判断是否存在指定的风刻杠
func CheckFengKeGang(chiCards []*ChiCard, cards []AICard, feng int32) bool {
	aiCards := MergeChiCard(chiCards, cards, false)
	for _, v := range aiCards {
		if v.Num >= 3 && v.Card == feng {
			return true
		}
	}
	return false
}

//判断四归一，手牌中有4个相同的牌
func Check4Gui1(cards []AICard) bool {
	for _, v := range cards {
		if v.Num == 4 {
			return true
		}
	}
	return false
}

//检查推不倒，由1、2、3、4、5、8、9饼，2、4、5、6、8、9条，白板组成
func CheckTuiBuDao(chiCards []*ChiCard, cards []AICard) bool {
	mask := uint64(8706014720) //1000000 110111010 110011111 000000000
	aiCards := MergeChiCard(chiCards, cards, false)
	for _, v := range aiCards {
		if mask&(uint64(0x01)<<cardId2Bit(v.Card)) == 0 {
			return false
		}
	}
	return true
}

//检查五门齐
func Check5MenQi(chiCards []*ChiCard, cards []AICard) bool {
	mask := uint8(0)
	aiCards := MergeChiCard(chiCards, cards, false)
	for _, v := range aiCards {
		if c := v.Card % MAHJONG_MASK; c >= MAHJONG_HONGZHONG {
			mask |= uint8(0x01) << 4
		} else if c >= MAHJONG_DONG {
			mask |= uint8(0x01) << 3
		} else {
			mask |= uint8(0x01) << uint32(v.Card/MAHJONG_MASK-COLOR_WAN)
		}
	}
	return mask == 0x1f
}

//检查缺一门
func CheckQueYiMen(chiCards []*ChiCard, cards []AICard) bool {
	mask := uint8(0x07)
	aiCards := MergeChiCard(chiCards, cards, false)
	for _, v := range aiCards {
		if c := v.Card / MAHJONG_MASK; c == COLOR_OTHER {
			continue
		} else {
			mask &^= uint8(0x01) << uint32(c-COLOR_WAN)
		}
	}
	return mask == 1 || mask == 2 || mask == 4
}

//检查清缺
func CheckQingQue(chiCards []*ChiCard, cards []AICard) bool {
	mask := uint8(0x07)
	aiCards := MergeChiCard(chiCards, cards, false)
	for _, v := range aiCards {
		if c := v.Card / MAHJONG_MASK; c == COLOR_OTHER {
			return false
		} else {
			mask &^= uint8(0x01) << uint32(c-COLOR_WAN)
		}
	}
	return mask == 1 || mask == 2 || mask == 4
}

//检查无字
func CheckWuZi(chiCards []*ChiCard, cards []AICard) bool {
	aiCards := MergeChiCard(chiCards, cards, false)
	for _, v := range aiCards {
		if v.Card/MAHJONG_MASK == COLOR_OTHER {
			return false
		}
	}
	return true
}

//检查双八支：两种颜色各占8个，其中一个颜色是两个杠
func CheckTwoBaZhi(chiCards []*ChiCard, cards []AICard) bool {
	color1, color2 := int32(0), int32(0)
	colorNum := make(map[int32]int32)
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_GANG || v.CardType == MJ_CHI_GANG_AN || v.CardType == MJ_CHI_GANG_WAN {
			color1 = v.CardId / MAHJONG_MASK
			colorNum[color1] += 4
		}
	}
	if len(colorNum) != 1 || colorNum[color1] != 8 {
		return false
	}
	for _, v := range cards {
		color2 = v.Card / MAHJONG_MASK
		if color1 == color2 {
			return false
		}
		colorNum[color2] += v.Num
	}
	if len(colorNum) != 2 || colorNum[color2] != 8 {
		return false
	}
	return true
}

//检查双5同
func CheckTwoWuTong(chiCards []*ChiCard, cards []AICard) bool {
	cardNum := make(map[int32]int32)
	cards = MergeChiCard(chiCards, cards, true)
	for _, v := range cards {
		cardNum[v.Card] += v.Num
	}
	num := 0
	for _, v := range cardNum {
		if v > 4 {
			num++
		}
	}
	return num >= 2
}

func GetTongNum(chiCards []*ChiCard, cards []AICard) int32 {
	cardNum := make(map[int32]int32)
	for _, v := range chiCards {
		cardNum[v.CardId] += 4
	}
	for _, v := range cards {
		cardNum[v.Card] += v.Num
	}
	maxNum := int32(0)
	for _, v := range cardNum {
		if v > maxNum {
			maxNum = v
		}
	}
	return maxNum
}

func GetZhiNum(chiCards []*ChiCard, cards []AICard) int32 {
	colorNum := make(map[int32]int32)
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_GANG || v.CardType == MJ_CHI_GANG_WAN || v.CardType == MJ_CHI_GANG_AN {
			colorNum[v.CardId/MAHJONG_MASK] += 4
		} else {
			colorNum[v.CardId/MAHJONG_MASK] += 3
		}
	}
	for _, v := range cards {
		colorNum[v.Card/MAHJONG_MASK] += v.Num
	}
	maxNum := int32(0)
	for _, v := range colorNum {
		if v > maxNum {
			maxNum = v
		}
	}
	return maxNum
}

func Check258Jiang2(cards []AICard, lzNum int32) bool {
	if ok, groups := CheckCommonHu(cards, lzNum, true); ok {
		for _, v := range groups {
			for _, vv := range v {
				if vv[2] == 0 {
					if c := vv[0] % MAHJONG_MASK; vv[0] == MAHJONG_LZ || c == MAHJONG_2 || c == MAHJONG_5 || c == MAHJONG_8 {
						return true
					}
				}
			}
		}
	}
	return false
}

func GetYiBanGaoNum2(chiCards []*ChiCard, cards []AICard, lzNum int32) int32 {
	num := int32(0)
	if ok, groups := CheckCommonHu(cards, lzNum, true); ok {
		for _, v := range groups {
			if n := GetYiBanGaoNum(chiCards, v, false); n > num {
				num = n
			}
		}
	}
	return num
}

// ***********************以下为切割牌型的判定*************************

func CheckDuiDuiHuForLZ(chiCards []*ChiCard, cards [][3]int32) bool {
	for _, v := range chiCards {
		if v.CardType != MJ_CHI_GANG &&
			v.CardType != MJ_CHI_GANG_WAN &&
			v.CardType != MJ_CHI_GANG_AN &&
			v.CardType != MJ_CHI_PENG {
			return false
		}
	}
	duiNum := 0
	for i := 0; i < len(cards); i++ {
		if cards[i][2] == 0 {
			duiNum++
			if duiNum > 1 {
				return false
			}
		}
		if cards[i][0] != cards[i][1] && cards[i][1] != MAHJONG_LZ {
			return false
		}
	}
	return true
}

func CheckCaiShenNiu(cards [][3]int32) bool {
	for i := 0; i < len(cards); i++ {
		if cards[i][2] == 0 {
			if cards[i][0] == MAHJONG_LZ {
				return true
			}
			break
		}
	}
	return false
}

func CheckCaiShenTou(cards [][3]int32, huCard int32) bool {
	for i := 0; i < len(cards); i++ {
		if cards[i][2] == 0 {
			if cards[i][1] == MAHJONG_LZ && cards[i][0] == huCard {
				return true
			}
			break
		}
	}
	return false
}

//检查大四喜
func CheckDaSiXi(chiCards []*ChiCard, cards [][3]int32) bool {
	num := 0
	for _, v := range chiCards {
		if c := v.CardId % MAHJONG_MASK; c >= MAHJONG_DONG && c <= MAHJONG_BEI {
			num++
		}
	}
	for _, v := range cards {
		if c := v[0] % MAHJONG_MASK; (c >= MAHJONG_DONG && c <= MAHJONG_BEI || v[0] == MAHJONG_LZ) && v[0] == v[1] && v[2] != 0 {
			num++
		}
	}
	return num == 4
}

//检查小四喜
func CheckXiaoSiXi(chiCards []*ChiCard, cards [][3]int32) bool {
	num, jiang := 0, false
	for _, v := range chiCards {
		if c := v.CardId % MAHJONG_MASK; c >= MAHJONG_DONG && c <= MAHJONG_BEI {
			num++
		}
	}
	for _, v := range cards {
		if c := v[0] % MAHJONG_MASK; (c >= MAHJONG_DONG && c <= MAHJONG_BEI || v[0] == MAHJONG_LZ) && v[0] == v[1] {
			num++
			if v[2] == 0 {
				jiang = true
			}
		}
	}
	return num == 4 && jiang
}

//检查大三元
func CheckDaSanYuan(chiCards []*ChiCard, cards [][3]int32) bool {
	num := 0
	for _, v := range chiCards {
		if c := v.CardId % MAHJONG_MASK; c >= MAHJONG_HONGZHONG && c <= MAHJONG_BAI {
			num++
		}
	}
	for _, v := range cards {
		if c := v[0] % MAHJONG_MASK; (c >= MAHJONG_HONGZHONG && c <= MAHJONG_BAI || v[0] == MAHJONG_LZ) && v[0] == v[1] && v[2] != 0 {
			num++
		}
	}
	return num == 3
}

//检查小三元
func CheckXiaoSanYuan(chiCards []*ChiCard, cards [][3]int32) bool {
	num, jiang := 0, false
	for _, v := range chiCards {
		if c := v.CardId % MAHJONG_MASK; c >= MAHJONG_HONGZHONG && c <= MAHJONG_BAI {
			num++
		}
	}
	for _, v := range cards {
		if c := v[0] % MAHJONG_MASK; (c >= MAHJONG_HONGZHONG && c <= MAHJONG_BAI || v[0] == MAHJONG_LZ) && v[0] == v[1] {
			num++
			if v[2] == 0 {
				jiang = true
			}
		}
	}
	return num == 3 && jiang
}

//检查连七对
func CheckLian7Dui(cards [][3]int32) bool {
	if len(cards) < 7 {
		return false
	}
	min, max := int32(10000), int32(0)
	for _, v := range cards {
		if v[2] != 0 {
			return false
		}
		if v[0] > max {
			max = v[0]
		}
		if v[0] < min {
			min = v[0]
		}
	}
	if max-min != 6 {
		return false
	}
	return true
}

//检查清幺九
func CheckQingYaoJiu(chiCards []*ChiCard, cards [][3]int32) bool {
	for _, v := range chiCards {
		if c := v.CardId % MAHJONG_MASK; v.CardType == MJ_CHI_CHI || c != MAHJONG_1 && c != MAHJONG_9 {
			return false
		}
	}
	for _, v := range cards {
		for _, vv := range v {
			if c := vv % MAHJONG_MASK; c != MAHJONG_1 && c != MAHJONG_9 && vv != MAHJONG_LZ && vv != 0 {
				return false
			}
		}
	}
	return true
}

//检查混幺九
func CheckHunYaoJiu(chiCards []*ChiCard, cards [][3]int32) bool {
	ziyise := true
	for _, v := range chiCards {
		if c := v.CardId % MAHJONG_MASK; v.CardType == MJ_CHI_CHI || c > MAHJONG_1 && c < MAHJONG_9 {
			return false
		} else if c < MAHJONG_DONG || c > MAHJONG_BAI {
			ziyise = false
		}
	}
	for _, v := range cards {
		for _, vv := range v {
			if vv != 0 {
				if c := vv % MAHJONG_MASK; c > MAHJONG_1 && c < MAHJONG_9 {
					return false
				} else if c < MAHJONG_DONG || c > MAHJONG_BAI {
					ziyise = false
				}
			}
		}
	}
	return !ziyise
}

func getLongMask(chiCards []*ChiCard, cards [][3]int32) (uint32, bool) {
	mask := uint32(0)
	for _, v := range chiCards {
		if c := v.CardId % MAHJONG_MASK; v.CardType == MJ_CHI_CHI && c >= MAHJONG_1 && c <= MAHJONG_9 {
			card := v.CardId
			if v.ChiPosBit&(0x01<<1) != 0 {
				card = v.CardId - 2
			} else if v.ChiPosBit&(0x01<<2) != 0 {
				card = v.CardId - 1
			}
			mask |= uint32(0x01) << cardId2Bit(card)
		}
	}
	lzFree := false
	for _, v := range cards {
		if v[2] != 0 {
			if v[0] != v[1] {
				if v[1] == MAHJONG_LZ {
					mask |= uint32(0x01) << cardId2Bit(v[0]-(v[0]%MAHJONG_MASK-1)%3)
				} else if v[2] == MAHJONG_LZ {
					if v[0]+1 == v[1] && (v[0]%MAHJONG_MASK-1)%3 != 0 {
						mask |= uint32(0x01) << cardId2Bit(v[0]-1)
					} else {
						mask |= uint32(0x01) << cardId2Bit(v[0])
					}
				} else {
					mask |= uint32(0x01) << cardId2Bit(v[0])
				}
			} else if v[0] == MAHJONG_LZ {
				lzFree = true
			}
		}
	}
	return mask, lzFree
}

//检查清龙
func CheckQingLong(chiCards []*ChiCard, cards [][3]int32) bool {
	mask, lzFree := getLongMask(chiCards, cards)

	if mask&0x49 == 0x49 || mask>>9&0x49 == 0x49 || mask>>18&0x49 == 0x49 { //0x49=>001 001 001
		return true
	}
	if lzFree { //0x09=>000 001 001, 0x41=>001 000 001, 0x48=>001 001 000
		if mask&0x09 == 0x09 || mask>>9&0x09 == 0x09 || mask>>18&0x09 == 0x09 ||
			mask&0x41 == 0x41 || mask>>9&0x41 == 0x41 || mask>>18&0x41 == 0x41 ||
			mask&0x48 == 0x48 || mask>>9&0x48 == 0x48 || mask>>18&0x48 == 0x48 {
			return true
		}
	}
	return false
}

//检查花龙
func CheckHuaLong(chiCards []*ChiCard, cards [][3]int32) bool {
	// 001000000 000001000 000000001
	// 001000000 000000001 000001000
	// 000001000 001000000 000000001
	// 000001000 000000001 001000000
	// 000000001 001000000 000001000
	// 000000001 000001000 001000000
	group := []uint32{16781313, 16777736, 2129921, 2097728, 294920, 266304}
	mask, lzFree := getLongMask(chiCards, cards)
	for _, v := range group {
		if mask&v == v {
			return true
		}
	}
	if lzFree {
		for _, v := range group {
			count := 0
			res := mask&v ^ v
			for i := uint32(0); i < 27; i++ {
				if res&(0x01<<i) != 0 {
					count++
					if count > 1 {
						break
					}
				}
			}
			if count == 1 {
				return true
			}
		}
	}
	return false
}

//获取暗刻数量，暂不考虑带癞子情况
func GetAnKeNum(chiCards []*ChiCard, cards [][3]int32, huCard int32, zimo bool) int32 {
	existOther, keNum := false, int32(0)
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_GANG_AN { //暗杠充当暗刻
			keNum++
		}
	}
	for _, v := range cards {
		if v[2] == 0 {
			if v[0] == huCard {
				existOther = true
			}
		} else if v[0] != v[1] {
			if v[0] == huCard || v[1] == huCard || v[2] == huCard {
				existOther = true
			}
		} else {
			keNum++
		}
	}
	if zimo || existOther {
		return keNum
	}
	return keNum - 1
}

func getSameSeq(chiCards []*ChiCard, cards [][3]int32) map[int32]int32 {
	mseq := make(map[int32]int32)
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_CHI {
			card := v.CardId
			if v.ChiPosBit&(0x01<<1) != 0 {
				card = v.CardId - 2
			} else if v.ChiPosBit&(0x01<<2) != 0 {
				card = v.CardId - 1
			}
			mseq[card]++
		}
	}
	for _, v := range cards {
		if v[0] != v[1] && v[1] != 0 {
			mseq[v[0]]++
		}
	}
	return mseq
}

//检查一色四同顺
func CheckYiSe4SameSeq(chiCards []*ChiCard, cards [][3]int32) bool {
	mseq := getSameSeq(chiCards, cards)
	for _, v := range mseq {
		if v == 4 {
			return true
		}
	}
	return false
}

//检查一色三同顺
func CheckYiSe3SameSeq(chiCards []*ChiCard, cards [][3]int32) bool {
	mseq := getSameSeq(chiCards, cards)
	for _, v := range mseq {
		if v == 3 {
			return true
		}
	}
	return false
}

//获取一般高数量，等价一色两同顺
func GetYiBanGaoNum(chiCards []*ChiCard, cards [][3]int32, dai19 bool) int32 {
	num := int32(0)
	mseq := getSameSeq(chiCards, cards)
	for k, v := range mseq {
		if !dai19 || k%MAHJONG_MASK == MAHJONG_1 || k%MAHJONG_MASK == MAHJONG_9 {
			num += v / 2
		}
	}
	return num
}

//检查三色三同顺
func Check3Se3SameSeq(chiCards []*ChiCard, cards [][3]int32) bool {
	mc := make(map[int32]int32)
	mseq := getSameSeq(chiCards, cards)
	for k, _ := range mseq {
		mc[k%MAHJONG_MASK]++
	}
	for _, v := range mc {
		if v == 3 {
			return true
		}
	}
	return false
}

//获取喜相逢数量，等价两色两同顺
func GetXiXiangFengNum(chiCards []*ChiCard, cards [][3]int32) int32 {
	mc := make(map[int32]int32)
	mseq := getSameSeq(chiCards, cards)
	for k, _ := range mseq {
		mc[k%MAHJONG_MASK]++
	}
	num := int32(0)
	for _, v := range mc {
		if v == 2 {
			num++
		}
	}
	return num
}

func getkeGangMask(chiCards []*ChiCard, cards [][3]int32) uint32 {
	mask := uint32(0)
	for _, v := range chiCards {
		if v.CardType != MJ_CHI_CHI {
			mask |= uint32(0x01) << cardId2Bit(v.CardId)
		}
	}
	for _, v := range cards {
		if v[2] != 0 && v[0] == v[1] {
			mask |= uint32(0x01) << cardId2Bit(v[0])
		}
	}
	return mask
}

//判断三种颜色是否3连续
func check3Se3Seq(mask uint32) bool {
	summask := mask&0x01ff | (mask>>9)&0x01ff | (mask>>18)&0x01ff

	indexs := []uint32{}
	for i := uint32(0); i < 7; i++ {
		if (summask>>i)&0x7 == 0x7 { //连续三位依次递增
			indexs = append(indexs, i)
		}
	}
	for _, v := range indexs {
		if mask&(uint32(0x01)<<v) != 0 && mask&(uint32(0x01)<<(9+v+1)) != 0 && mask&(uint32(0x01)<<(18+v+2)) != 0 ||
			mask&(uint32(0x01)<<v) != 0 && mask&(uint32(0x01)<<(18+v+1)) != 0 && mask&(uint32(0x01)<<(9+v+2)) != 0 ||
			mask&(uint32(0x01)<<(9+v)) != 0 && mask&(uint32(0x01)<<(v+1)) != 0 && mask&(uint32(0x01)<<(18+v+2)) != 0 ||
			mask&(uint32(0x01)<<(9+v)) != 0 && mask&(uint32(0x01)<<(18+v+1)) != 0 && mask&(uint32(0x01)<<(v+2)) != 0 ||
			mask&(uint32(0x01)<<(18+v)) != 0 && mask&(uint32(0x01)<<(v+1)) != 0 && mask&(uint32(0x01)<<(9+v+2)) != 0 ||
			mask&(uint32(0x01)<<(18+v)) != 0 && mask&(uint32(0x01)<<(9+v+1)) != 0 && mask&(uint32(0x01)<<(v+2)) != 0 {
			return true
		}
	}
	return false
}

//检查一色四节高
func CheckYiSe4JieGao(chiCards []*ChiCard, cards [][3]int32) bool {
	mask := getkeGangMask(chiCards, cards)
	for i := uint32(0); i < 3; i++ {
		for j := uint32(0); j < 6; j++ {
			if (mask>>(i*9+j))&0xf == 0xf { //连续四位依次递增
				return true
			}
		}
	}
	return false
}

//检查一色三节高
func CheckYiSe3JieGao(chiCards []*ChiCard, cards [][3]int32) bool {
	mask := getkeGangMask(chiCards, cards)
	for i := uint32(0); i < 3; i++ {
		for j := uint32(0); j < 7; j++ {
			if (mask>>(i*9+j))&0x7 == 0x7 { //连续三位依次递增
				return true
			}
		}
	}
	return false
}

//检查三色三节高
func Check3Se3JieGao(chiCards []*ChiCard, cards [][3]int32) bool {
	mask := getkeGangMask(chiCards, cards)
	return check3Se3Seq(mask)
}

func getSeqMask(chiCards []*ChiCard, cards [][3]int32) uint32 {
	mask := uint32(0)
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_CHI {
			card := v.CardId
			if v.ChiPosBit&(0x01<<1) != 0 {
				card = v.CardId - 2
			} else if v.ChiPosBit&(0x01<<2) != 0 {
				card = v.CardId - 1
			}
			mask |= uint32(0x01) << cardId2Bit(card)
		}
	}
	for _, v := range cards {
		if v[0] != v[1] {
			mask |= uint32(0x01) << cardId2Bit(v[0])
		}
	}
	return mask
}

//检查一色四步高，每次递增1或2的连续4个同色顺子
func CheckYiSe4BuGao(chiCards []*ChiCard, cards [][3]int32) bool {
	mask := getSeqMask(chiCards, cards)
	for i := uint32(0); i < 3; i++ {
		if (mask>>(i*9))&0x55 == 0x55 { //递增2
			return true
		}
		for j := uint32(0); j < 6; j++ {
			if (mask>>(i*9+j))&0xf == 0xf { //递增1
				return true
			}
		}
	}
	return false
}

//检查一色三步高，每次递增1或2的连续3个同色顺子
func CheckYiSe3BuGao(chiCards []*ChiCard, cards [][3]int32) bool {
	mask := getSeqMask(chiCards, cards)
	for i := uint32(0); i < 3; i++ {
		if (mask>>(i*9))&0x15 == 0x15 || (mask>>(i*9+1))&0x15 == 0x15 || (mask>>(i*9+2))&0x15 == 0x15 { //递增2
			return true
		}
		for j := uint32(0); j < 7; j++ {
			if (mask>>(i*9+j))&0x7 == 0x7 { //递增1
				return true
			}
		}
	}
	return false
}

//检查三色三步高，每次递增1的连续3个不同花色顺子
func Check3Se3BuGao(chiCards []*ChiCard, cards [][3]int32) bool {
	mask := getSeqMask(chiCards, cards)
	return check3Se3Seq(mask)
}

//检查带19（每组牌均带19牌或字牌)
func CheckDai19(chiCards []*ChiCard, cards [][3]int32) bool {
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_CHI {
			card := v.CardId
			if v.ChiPosBit&(0x01<<1) != 0 {
				card = v.CardId - 2
			} else if v.ChiPosBit&(0x01<<2) != 0 {
				card = v.CardId - 1
			}
			if c := card % MAHJONG_MASK; c != MAHJONG_1 && c != MAHJONG_7 {
				return false
			}
		} else if c := v.CardId % MAHJONG_MASK; c > MAHJONG_1 && c < MAHJONG_9 {
			return false
		}
	}

	for _, v := range cards {
		if v[0] == v[1] {
			if c := v[0] % MAHJONG_MASK; c > MAHJONG_1 && c < MAHJONG_9 {
				return false
			}
		} else {
			if c := v[0] % MAHJONG_MASK; c != MAHJONG_1 && c != MAHJONG_7 {
				return false
			}
		}
	}
	return true
}

//检查全带5
func CheckQuanDai5(chiCards []*ChiCard, cards [][3]int32) bool {
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_CHI {
			card := v.CardId
			if v.ChiPosBit&(0x01<<1) != 0 {
				card = v.CardId - 2
			} else if v.ChiPosBit&(0x01<<2) != 0 {
				card = v.CardId - 1
			}
			if c := card % MAHJONG_MASK; c < MAHJONG_3 || c > MAHJONG_5 {
				return false
			}
		} else if c := v.CardId % MAHJONG_MASK; c != MAHJONG_5 {
			return false
		}
	}

	for _, v := range cards {
		if v[0] == v[1] {
			if c := v[0] % MAHJONG_MASK; c != MAHJONG_5 {
				return false
			}
		} else {
			if c := v[0] % MAHJONG_MASK; c < MAHJONG_3 || c > MAHJONG_5 {
				return false
			}
		}
	}
	return true
}

//检查平胡
func CheckPingHu(chiCards []*ChiCard, cards [][3]int32) bool {
	for _, v := range chiCards {
		if v.CardType != MJ_CHI_CHI {
			return false
		}
	}
	for _, v := range cards {
		if v[2] == 0 {
			if v[0]%MAHJONG_MASK >= MAHJONG_DONG {
				return false
			}
		} else if v[0] == v[1] {
			return false
		}
	}
	return true
}

//检查三同刻杠
func Check3SameKeGang(chiCards []*ChiCard, cards [][3]int32) bool {
	mke := make(map[int32]int32)
	for _, v := range chiCards {
		if v.CardType != MJ_CHI_CHI {
			mke[v.CardId%MAHJONG_MASK]++
		}
	}
	for _, v := range cards {
		if v[0] == v[1] && v[1] == v[2] {
			mke[v[0]%MAHJONG_MASK]++
		}
	}
	for _, v := range mke {
		if v == 3 {
			return true
		}
	}
	return false
}

//获取双同刻杠数量，暂不包含癞子牌
func GetTwoSameKeNum(chiCards []*ChiCard, cards [][3]int32) int32 {
	mke := make(map[int32]int32)
	for _, v := range chiCards {
		if v.CardType != MJ_CHI_CHI {
			mke[v.CardId%MAHJONG_MASK]++
		}
	}
	for _, v := range cards {
		if v[0] == v[1] && v[1] == v[2] {
			mke[v[0]%MAHJONG_MASK]++
		}
	}
	num := int32(0)
	for _, v := range mke {
		if v == 2 {
			num++
		}
	}
	return num
}

//获取连六数量，暂不包含癞子牌
func GetLian6Num(chiCards []*ChiCard, cards [][3]int32, dai19 bool) int32 {
	num := int32(0)
	mask := getSeqMask(chiCards, cards)
	for i := uint32(0); i < 3; i++ {
		for j := uint32(0); j < 4; j++ {
			if (mask>>(i*9+j))&0x09 == 0x09 {
				num++
			}
			if dai19 {
				j += 3
			}
		}
	}
	if num == 1 && GetYiBanGaoNum(chiCards, cards, false) == 2 {
		num++
	}
	return num
}

//获取老少副数量, 暂不包含癞子牌
func GetLaoShaoFuNum(chiCards []*ChiCard, cards [][3]int32) int32 {
	num := int32(0)
	mask := getSeqMask(chiCards, cards)
	for i := uint32(0); i < 3; i++ {
		if (mask>>(i*9))&0x41 == 0x41 {
			num++
		}
	}
	if num == 1 && GetYiBanGaoNum(chiCards, cards, false) == 2 {
		num++
	}
	return num
}

//获取幺九刻杠数量，暂不包含癞子
func GetZi19KeGangNum(chiCards []*ChiCard, cards [][3]int32) int32 {
	num := int32(0)
	for _, v := range chiCards {
		if c := v.CardId % MAHJONG_MASK; v.CardType != MJ_CHI_CHI && c == MAHJONG_1 || c == MAHJONG_9 || c >= MAHJONG_DONG {
			num++
		}
	}
	for _, v := range cards {
		if c := v[0] % MAHJONG_MASK; v[0] == v[1] && v[1] == v[2] && (c == MAHJONG_1 || c == MAHJONG_9 || c >= MAHJONG_DONG) {
			num++
		}
	}
	return num
}

//获取四核数量，不含癞子
func GetSiHeNum(cards [][3]int32) int32 {
	cardNum := make(map[int32]int32)
	for _, v := range cards {
		if v[0] == v[1] && v[1] == v[2] {
			cardNum[v[0]] += 3
		}
		if v[0] != v[1] {
			for _, vv := range v {
				if cardNum[vv] == 0 || cardNum[vv] == 3 {
					cardNum[vv]++
				}
			}
		}
	}
	num := int32(0)
	for _, v := range cardNum {
		if v == 4 {
			num++
		}
	}
	return num
}

// 将将胡,全258,可以不是成牌牌形
func CheckJiangJiangHu(chiCards []*ChiCard, cards []AICard) bool {
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_CHI {
			return false
		} else if c := v.CardId % MAHJONG_MASK; c != MAHJONG_2 && c != MAHJONG_5 && c != MAHJONG_8 {
			return false
		}
	}

	for _, v := range cards {
		if c := v.Card % MAHJONG_MASK; c != MAHJONG_2 && c != MAHJONG_5 && c != MAHJONG_8 {
			return false
		}
	}
	return true
}

func CheckJiangJiangTing(chiCards []*ChiCard, cards []AICard) []int32 {
	tings := []int32{}
	if CheckJiangJiangHu(chiCards, cards) && len(cards) > 0 {
		for _, v := range cards {
			if v.Num < 4 {
				tings = append(tings, v.Card)
			}
		}
	}
	return tings
}

func checkLuanSanFeng(masks []uint64, groups *[][][3]int32, sliceGroups [][3]int32, cards []AICard, all bool) bool {
	bHu := bool(false)
	cg := uint64(0)
	for _, v := range cards {
		cg |= uint64(0x01) << cardId2Bit(v.Card)
	}

	for _, v := range masks {
		if v&cg != v {
			continue
		}

		count := int32(0)
		restCards := make([]AICard, len(cards))
		copy(restCards, cards)
		for i := uint32(27); i < 34; i++ {
			if v&(0x01<<i) != 0 {
				count++
				card := cardBit2Id(i)
				for index := 0; index < len(restCards); index++ {
					if restCards[index].Card == card {
						restCards[index].Num--
						if restCards[index].Num == 0 {
							restCards = append(restCards[:index], restCards[index+1:]...)
							index--
						}
					}
				}
			}
		}
		sliceGroupsBK := [][3]int32{}
		copy(sliceGroupsBK, sliceGroups)

		group := make([][3]int32, count/3)
		if all {
			num := 0
			for i := uint32(27); i < 34; i++ {
				if v&(uint64(0x01)<<i) != 0 {
					group[num/3][num%3] = cardBit2Id(i)
					num++
				}
			}
		}
		sliceGroupsBK = append(sliceGroups, group...)

		ok, restGroups := CheckCommonHu(restCards, 0, true)
		if ok {
			bHu = true
			if all {
				for _, v := range restGroups {
					*groups = append(*groups, sliceGroupsBK)
					(*groups)[len(*groups)-1] = append((*groups)[len(*groups)-1], v...)
				}
			}
		} else {
			bHu = checkLuanSanFeng(masks, groups, sliceGroupsBK, restCards, all)
		}
	}

	return bHu
}

//检查不含癞子乱三风
func CheckLuanSanFengHu(cards []AICard, all bool) (bool, [][][3]int32) {
	// 0001110 000000000 000000000 000000000
	// 0001101 000000000 000000000 000000000
	// 0001011 000000000 000000000 000000000
	// 0000111 000000000 000000000 000000000
	// 1110000 000000000 000000000 000000000
	// 1111110 000000000 000000000 000000000
	// 1111101 000000000 000000000 000000000
	// 1111011 000000000 000000000 000000000
	// 1110111 000000000 000000000 000000000
	masks := []uint64{1879048192, 1744830464, 1476395008, 939524096, 15032385536, 16911433728, 16777216000, 16508780544, 15971909632}
	groups := [][][3]int32{}
	sliceGroups := [][3]int32{}

	if ok := checkLuanSanFeng(masks, &groups, sliceGroups, cards, all); ok {
		return ok, groups
	}

	return false, groups
}

// 检查258做将
func Check258Jiang(cards [][3]int32) bool {
	for _, v := range cards {
		if v[2] != 0 {
			continue
		}

		if c := v[0] % MAHJONG_MASK; c == MAHJONG_2 || c == MAHJONG_5 || c == MAHJONG_8 || v[0] == MAHJONG_LZ {
			return true
		}
	}
	return false
}

func CheckBanBanHu(chiCards []*ChiCard, cards []AICard) bool {
	if len(chiCards) > 0 {
		return false
	}

	for _, v := range cards {
		if c := v.Card % MAHJONG_MASK; c == MAHJONG_2 || c == MAHJONG_5 || c == MAHJONG_8 {
			return false
		}
	}
	return true
}

func CheckLiuLiuShun(chiCards []*ChiCard, cards []AICard) bool {
	if len(chiCards) > 0 {
		return false
	}

	var keNum int
	for _, v := range cards {
		if v.Num == 3 {
			keNum++
		}
		if keNum == 2 {
			return true
		}
	}
	return false
}

func CheckMingSiGuiYi(chiCards []*ChiCard, huCard int32) bool {
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_PENG && v.CardId == huCard {
			return true
		}
	}
	return false
}

func CheckAnSiGuiYi(cards []AICard, huCard int32) bool {
	for i := 0; i < len(cards); i++ {
		if cards[i].Card == huCard && cards[i].Num == 4 {
			return true
		}
	}
	return false
}
