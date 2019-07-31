package mahjong

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

func uniqueHuGroup(groups [][][3]int32) [][][3]int32 {
	type GroupInfo struct {
		group [][3]int32
		seq   []uint64
	}
	sGroupInfo := make([]*GroupInfo, 0, len(groups))
	for _, group := range groups {
		pGroupInfo := &GroupInfo{
			group: group,
			seq:   make([]uint64, len(group)-1),
		}
		for i := 0; i < len(group)-1; i++ {
			pGroupInfo.seq[i] = uint64(group[i][0]) | (uint64(group[i][1]) << 20) | (uint64(group[i][2]) << 40)
		}
		sort.Slice(pGroupInfo.seq, func(i, j int) bool { return pGroupInfo.seq[i] < pGroupInfo.seq[j] })
		sGroupInfo = append(sGroupInfo, pGroupInfo)
	}
	sGroupInfo = util.UniqueSlice2(sGroupInfo, func(i, j int) bool {
		seq1, seq2 := sGroupInfo[i].seq, sGroupInfo[j].seq
		if len(seq1) != len(seq2) {
			return false
		}
		for k := 0; k < len(seq1); k++ {
			if seq1[k] != seq2[k] {
				return false
			}
		}
		return true
	}, false).([]*GroupInfo)
	ret := make([][][3]int32, 0, len(sGroupInfo))
	for _, v := range sGroupInfo {
		ret = append(ret, v.group)
	}
	return ret
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
				groups = append(groups, uniqueHuGroup(createHuGroup(treeCard))...)
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
		sort.Slice(tingInfo, func(i, j int) bool { return tingInfo[i] < tingInfo[j] })
		tingInfo = util.UniqueSlice(tingInfo, true).([]int32)
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
