package logic

import (
	"sort"
	"util"
	//"github.com/astaxie/beego/logs"
)

type AICard struct {
	Card int32
	Num  int32
}

func CheckCommonHu(cards []AICard, all bool) (bool, [][][3]int32) {
	bHu := false
	var groups [][][3]int32
	if all {
		groups = make([][][3]int32, 0, len(cards))
	}
	aiCards := make([]AICard, len(cards))
	aiCardsBk := make([]AICard, len(aiCards))
	copy(aiCards, cards)
	sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })
	copy(aiCardsBk, aiCards)
	for i := 0; i < len(aiCards); i++ {
		if aiCards[i].Num >= 2 {
			if aiCards[i].Num > 2 && aiCards[i].Card%MAHJONG_MASK >= MAHJONG_DONG {
				continue
			}
			aiCards[i].Num -= 2
			var subGroup [][3]int32
			if all {
				subGroup = make([][3]int32, 0, len(cards))
				subGroup = append(subGroup, [3]int32{aiCards[i].Card, aiCards[i].Card, 0})
			}
			for j := 0; ; {
				for {
					if j >= len(aiCards) {
						if all {
							groups = append(groups, subGroup)
							bHu = true
							break
						} else {
							return true, groups
						}
					}
					if aiCards[j].Num > 0 {
						break
					}
					j++
				}
				if j >= len(aiCards) {
					break
				}
				if aiCards[j].Num >= 3 {
					aiCards[j].Num -= 3
					if all {
						subGroup = append(subGroup, [3]int32{aiCards[j].Card, aiCards[j].Card, aiCards[j].Card})
					}
				} else {
					if j > len(aiCards)-3 || aiCards[j].Card%MAHJONG_MASK >= MAHJONG_DONG ||
						aiCards[j+1].Num == 0 || aiCards[j+2].Num == 0 ||
						aiCards[j].Card != aiCards[j+1].Card-1 || aiCards[j+1].Card+1 != aiCards[j+2].Card {
						break
					}
					aiCards[j].Num--
					aiCards[j+1].Num--
					aiCards[j+2].Num--
					if all {
						subGroup = append(subGroup, [3]int32{aiCards[j].Card, aiCards[j+1].Card, aiCards[j+2].Card})
					}
				}
			}
			copy(aiCards, aiCardsBk)
		}
	}
	return bHu, groups
}

func CheckCommonHuForLZ(cards []AICard, lzNum int32, all bool) (bool, [][][3]int32) {
	if lzNum <= 0 {
		return CheckCommonHu(cards, all)
	}

	bHu := false
	var groups [][][3]int32
	if all {
		groups = make([][][3]int32, 0, len(cards))
	}

	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })

	aiCardsBk := make([]AICard, len(aiCards))
	copy(aiCardsBk, aiCards)
	lzNumBk := lzNum

	for i := -1; i < len(aiCards); i++ {
		var subGroup [][3]int32
		if i == -1 && lzNum < 2 {
			continue
		}
		if all {
			subGroup = make([][3]int32, 0, len(cards))
			if i == -1 {
				subGroup = append(subGroup, [3]int32{MAHJONG_LZ, MAHJONG_LZ, 0})
			} else if aiCards[i].Num == 1 {
				subGroup = append(subGroup, [3]int32{aiCards[i].Card, MAHJONG_LZ, 0})
			} else {
				subGroup = append(subGroup, [3]int32{aiCards[i].Card, aiCards[i].Card, 0})
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
		for j := 0; ; {
			for {
				if j >= len(aiCards) {
					if all {
						if lzNum == 3 { //当牌型非常好的时候会出现这种极端情况
							subGroup = append(subGroup, [3]int32{MAHJONG_LZ, MAHJONG_LZ, MAHJONG_LZ})
						}
						groups = append(groups, subGroup)
						bHu = true
						break
					} else {
						return true, groups
					}
				}
				if aiCards[j].Num > 0 {
					break
				}
				j++
			}
			if j >= len(aiCards) {
				break
			}

			if aiCards[j].Num >= 3 {
				aiCards[j].Num -= 3
				if all {
					subGroup = append(subGroup, [3]int32{aiCards[j].Card, aiCards[j].Card, aiCards[j].Card})
				}
			} else if c := aiCards[j].Card % MAHJONG_MASK; c >= MAHJONG_1 && c <= MAHJONG_9 && j <= len(aiCards)-3 &&
				aiCards[j].Num >= 1 && aiCards[j+1].Num >= 1 && aiCards[j+2].Num >= 1 &&
				aiCards[j].Card == aiCards[j+1].Card-1 && aiCards[j+1].Card+1 == aiCards[j+2].Card &&
				(aiCards[j].Num == 1 || aiCards[j+1].Num == 2 || aiCards[j+2].Num > 1) {
				aiCards[j].Num--
				aiCards[j+1].Num--
				aiCards[j+2].Num--
				if all {
					subGroup = append(subGroup, [3]int32{aiCards[j].Card, aiCards[j+1].Card, aiCards[j+2].Card})
				}
			} else if aiCards[j].Num == 2 {
				aiCards[j].Num -= 2
				if lzNum <= 0 {
					break
				} else {
					lzNum--
				}
				if all {
					subGroup = append(subGroup, [3]int32{aiCards[j].Card, aiCards[j].Card, MAHJONG_LZ})
				}
			} else {
				if j > len(aiCards)-2 || aiCards[j].Card%MAHJONG_MASK >= MAHJONG_DONG ||
					aiCards[j+1].Card-aiCards[j].Card > 2 ||
					(aiCards[j+1].Num == 0 && (j > len(aiCards)-3 || aiCards[j+2].Num == 0 || aiCards[j+2].Card-aiCards[j].Card > 2)) {
					if lzNum < 2 {
						break
					} else {
						lzNum -= 2
					}
					if all {
						subGroup = append(subGroup, [3]int32{aiCards[j].Card, MAHJONG_LZ, MAHJONG_LZ})
					}
				} else {
					if lzNum < 1 {
						break
					} else {
						lzNum--
						if aiCards[j+1].Num == 0 {
							aiCards[j+2].Num--
							if all {
								subGroup = append(subGroup, [3]int32{aiCards[j].Card, aiCards[j+2].Card, MAHJONG_LZ})
							}
						} else {
							aiCards[j+1].Num--
							if all {
								subGroup = append(subGroup, [3]int32{aiCards[j].Card, aiCards[j+1].Card, MAHJONG_LZ})
							}
						}
					}
				}
				aiCards[j].Num--
			}
		}
		copy(aiCards, aiCardsBk)
		lzNum = lzNumBk
	}
	return bHu, groups
}

func CheckCommonTing(cards []AICard, all bool) []int32 {
	tingInfo := make([]int32, 0, 34)
	aiCards := make([]AICard, len(cards))
	aiCardsBk := make([]AICard, len(aiCards))
	copy(aiCards, cards)
	sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })
	copy(aiCardsBk, aiCards)
	for i := 0; i < len(aiCards); i++ {
		subTing := make([]int32, 0, 5) //每一轮听牌个数最多三张
		seq := make(map[int32]int32)
		if aiCards[i].Num >= 2 {
			if aiCards[i].Num > 2 && aiCards[i].Card%MAHJONG_MASK >= MAHJONG_DONG {
				continue
			}
			aiCards[i].Num -= 2
		} else if aiCards[i].Num == 1 {
			aiCards[i].Num--
			subTing = append(subTing, aiCards[i].Card)
		} else {
			continue
		}
		for j := 0; ; {
			for {
				if j >= len(aiCards) || aiCards[j].Num > 0 {
					break
				}
				j++
			}
			if j >= len(aiCards) {
				break
			}
			if aiCards[j].Card%MAHJONG_MASK >= MAHJONG_DONG && (aiCards[j].Num == 1 || aiCards[j].Num == 4) {
				subTing = subTing[:0]
				break
			} else if aiCards[j].Num >= 3 {
				aiCards[j].Num -= 3
			} else if c := aiCards[j].Card % MAHJONG_MASK; c >= MAHJONG_1 && c <= MAHJONG_9 && j <= len(aiCards)-3 &&
				aiCards[j].Num >= 1 && aiCards[j+1].Num >= 1 && aiCards[j+2].Num >= 1 &&
				aiCards[j].Card == aiCards[j+1].Card-1 && aiCards[j+1].Card+1 == aiCards[j+2].Card &&
				(aiCards[j].Num == 1 || aiCards[j+1].Num == 2 || aiCards[j+2].Num > 1) {
				aiCards[j].Num--
				aiCards[j+1].Num--
				aiCards[j+2].Num--
				if _, ok := seq[aiCards[j].Card]; !ok {
					seq[aiCards[j].Card] = 1
				} else {
					seq[aiCards[j].Card]++
				}
			} else if aiCards[j].Num == 2 {
				aiCards[j].Num -= 2
				if len(subTing) > 0 {
					subTing = subTing[:0]
					break
				} else {
					subTing = append(subTing, aiCards[j].Card)
				}
			} else {
				if len(subTing) > 0 || j > len(aiCards)-2 || aiCards[j+1].Card-aiCards[j].Card > 2 ||
					(aiCards[j+1].Num == 0 && (j > len(aiCards)-3 || aiCards[j+2].Num == 0 || aiCards[j+2].Card-aiCards[j].Card > 2)) {
					subTing = subTing[:0]
					break
				}
				if aiCards[j+1].Num > 0 {
					if aiCards[j].Card == aiCards[j+1].Card-1 {
						if aiCards[j].Card%MAHJONG_MASK == MAHJONG_1 {
							subTing = append(subTing, aiCards[j+1].Card+1)
						} else if aiCards[j+1].Card%MAHJONG_MASK == MAHJONG_9 {
							subTing = append(subTing, aiCards[j].Card-1)
						} else {
							subTing = append(subTing, aiCards[j].Card-1, aiCards[j+1].Card+1)
						}
					} else {
						subTing = append(subTing, aiCards[j].Card+1)
					}
					aiCards[j+1].Num--
				} else {
					subTing = append(subTing, aiCards[j].Card+1)
					aiCards[j+2].Num--
				}
				aiCards[j].Num--
			}
		}
		if len(subTing) > 0 {
			if subTing[0]%MAHJONG_MASK >= MAHJONG_4 && subTing[0]%MAHJONG_MASK <= MAHJONG_9 && seq[subTing[0]-2] > 0 {
				subTing = append(subTing, subTing[0]-3) //若胡4、5、6、7、8、9、47、58、69, 检测是否胡14、25、36、47、58、69、147、258、369
				if subTing[0]%MAHJONG_MASK >= MAHJONG_7 && subTing[0]%MAHJONG_MASK <= MAHJONG_9 && seq[subTing[0]-5] > 0 {
					subTing = append(subTing, subTing[0]-6) //继续检测是否胡147、258、369
				}
			}			
			tingInfo = append(tingInfo, subTing...)			
			if !all {
				break
			}
		}
		copy(aiCards, aiCardsBk)
	}
	if len(tingInfo) > 0 {
		//对结果进行排序除重
		sort.Slice(tingInfo, func(i, j int) bool { return tingInfo[i] < tingInfo[j] })
		tingInfo = util.UniqueSlice(tingInfo)
	}
	return tingInfo
}

func CheckCommonTingForLZ(cards []AICard, lzNum int32, all bool) []int32 {
	if lzNum <= 0 {
		return CheckCommonTing(cards, all)
	}

	tingInfo := make([]int32, 0, 34)
	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })

	//检查是否听所有牌
	huCards := make([]AICard, len(cards)+1)
	copy(huCards, cards)
	if ok, _ := CheckCommonHuForLZ(append(huCards, AICard{MAHJONG_ANY, 1}), lzNum, false); ok {
		tingInfo = append(tingInfo, MAHJONG_ANY)
		return tingInfo
	}

	aiCardsBk := make([]AICard, len(aiCards))
	copy(aiCardsBk, aiCards)
	lzNumBk := lzNum

	for i := 0; i < len(aiCards); i++ {
		subTing := make([]int32, 0, 10) //每一轮最多听8张
		seq := make(map[int32]int32)
		if aiCards[i].Num >= 2 {
			aiCards[i].Num -= 2
		} else if aiCards[i].Num == 1 {
			aiCards[i].Num--
			lzNum--
			subTing = append(subTing, aiCards[i].Card)
		} else {
			continue
		}
		for j := 0; ; {
			for {
				if j >= len(aiCards) || aiCards[j].Num > 0 {
					break
				}
				j++
			}
			if j >= len(aiCards) {
				break
			}
			if aiCards[j].Num >= 3 {
				aiCards[j].Num -= 3
			} else if c := aiCards[j].Card % MAHJONG_MASK; c >= MAHJONG_1 && c <= MAHJONG_9 && j <= len(aiCards)-3 &&
				aiCards[j].Num >= 1 && aiCards[j+1].Num >= 1 && aiCards[j+2].Num >= 1 &&
				aiCards[j].Card == aiCards[j+1].Card-1 && aiCards[j+1].Card+1 == aiCards[j+2].Card &&
				(aiCards[j].Num == 1 || aiCards[j+1].Num == 2 || aiCards[j+2].Num > 1) {
				aiCards[j].Num--
				aiCards[j+1].Num--
				aiCards[j+2].Num--
				if _, ok := seq[aiCards[j].Card]; !ok {
					seq[aiCards[j].Card] = 1
				} else {
					seq[aiCards[j].Card]++
				}
			} else {
				if lzNum < 0 {
					subTing = subTing[:0]
					break
				}
				if aiCards[j].Num == 2 {
					if lzNum >= 0 {
						lzNum--
					}
					if seq[aiCards[j].Card-2] > 0 {
						subTing = append(subTing, aiCards[j].Card-3) //检测是否胡前一张
					}
					subTing = append(subTing, aiCards[j].Card)
					aiCards[j].Num -= 2
				} else {
					if aiCards[j].Card%MAHJONG_MASK >= MAHJONG_DONG {
						if lzNum <= 0 {
							subTing = subTing[:0]
							break
						}
						lzNum -= 2
						subTing = append(subTing, aiCards[j].Card)
					} else if j < len(aiCards)-1 && aiCards[j+1].Num > 0 && aiCards[j+1].Card-aiCards[j].Card < 3 {
						if lzNum >= 0 {
							lzNum--
						}
						if aiCards[j+1].Card-aiCards[j].Card == 1 {
							if aiCards[j].Card%MAHJONG_MASK == MAHJONG_1 {
								subTing = append(subTing, aiCards[j+1].Card+1)
							} else if aiCards[j+1].Card%MAHJONG_MASK == MAHJONG_9 {
								if seq[aiCards[j].Card-3] > 0 {
									subTing = append(subTing, aiCards[j].Card-4) //若胡7，检测是否胡47
									if seq[aiCards[j].Card-6] > 0 {
										subTing = append(subTing, aiCards[j].Card-7) //继续检测是否胡147
									}
								}								
								subTing = append(subTing, aiCards[j].Card-1)
							} else {
								if aiCards[j].Card%MAHJONG_MASK > MAHJONG_4 && aiCards[j].Card%MAHJONG_MASK < MAHJONG_8 && seq[aiCards[j].Card-3] > 0 {
									subTing = append(subTing, aiCards[j].Card-4) //若胡47、58、69，检测是否胡147、、258、369
								}
								subTing = append(subTing, aiCards[j].Card-1, aiCards[j+1].Card+1)
							}
						} else {
							subTing = append(subTing, aiCards[j].Card+1)
						}
						aiCards[j+1].Num--
					} else if j < len(aiCards)-2 && aiCards[j+1].Num == 0 && aiCards[j+2].Num != 0 && aiCards[j+2].Card-aiCards[j].Card < 3 {
						if lzNum >= 0 {
							lzNum--
						}
						subTing = append(subTing, aiCards[j].Card+1)
						aiCards[j+2].Num--
					} else {
						if lzNum <= 0 {
							subTing = subTing[:0]
							break
						}
						lzNum -= 2
						if aiCards[j].Card%MAHJONG_MASK == MAHJONG_1 {
							subTing = append(subTing, aiCards[j].Card, aiCards[j].Card+1, aiCards[j].Card+2)
						} else if aiCards[j].Card%MAHJONG_MASK == MAHJONG_2 {
							subTing = append(subTing, aiCards[j].Card-1, aiCards[j].Card, aiCards[j].Card+1, aiCards[j].Card+2)
						} else if aiCards[j].Card%MAHJONG_MASK == MAHJONG_8 {
							subTing = append(subTing, aiCards[j].Card-2, aiCards[j].Card-1, aiCards[j].Card, aiCards[j].Card+1)
						} else if aiCards[j].Card%MAHJONG_MASK == MAHJONG_9 {
							subTing = append(subTing, aiCards[j].Card-2, aiCards[j].Card-1, aiCards[j].Card)
						} else {
							subTing = append(subTing, aiCards[j].Card-2, aiCards[j].Card-1, aiCards[j].Card, aiCards[j].Card+1, aiCards[j].Card+2)
						}
					}
					aiCards[j].Num--
				}
			}
		}
		if len(subTing) > 0 {
			tingInfo = append(tingInfo, subTing...)
			if !all {
				break
			}
		}
		copy(aiCards, aiCardsBk)
		lzNum = lzNumBk
	}

	if len(tingInfo) > 0 {
		//对结果进行排序除重
		sort.Slice(tingInfo, func(i, j int) bool { return tingInfo[i] < tingInfo[j] })
		tingInfo = util.UniqueSlice(tingInfo)
	}
	return tingInfo
}

func CheckQiDuiHu(cards []AICard) bool {
	for i := 0; i < len(cards); i++ {
		if cards[i].Num%2 != 0 {
			return false
		}
	}
	return true
}

func CheckQiDuiHuForLZ(cards []AICard, lzNum int32) bool {
	if lzNum <= 0 {
		return CheckQiDuiHu(cards)
	}

	for i := 0; i < len(cards); i++ {
		if cards[i].Num%2 != 0 {
			if lzNum > 0 {
				lzNum--
			} else {
				return false
			}
		}
	}
	return true
}

func CheckQiDuiTing(cards []AICard) []int32 {
	tingInfo := make([]int32, 0, 1)
	for i := 0; i < len(cards); i++ {
		if cards[i].Num%2 != 0 {
			if len(tingInfo) > 0 {
				return tingInfo[:0]
			} else {
				tingInfo = append(tingInfo, cards[i].Card)
			}
		}
	}
	return tingInfo
}

func CheckQiDuiTingForLZ(cards []AICard, lzNum int32) []int32 {
	if lzNum <= 0 {
		return CheckQiDuiTing(cards)
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
		sort.Slice(tingInfo, func(i, j int) bool { return tingInfo[i] < tingInfo[j] })
	}
	return tingInfo
}

func CheckQuanBuKaoHu(cards []AICard) bool {
	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })
	for i := 0; i < len(aiCards); i++ {
		if aiCards[i].Num > 1 || (aiCards[i].Card%MAHJONG_MASK < MAHJONG_DONG && i < len(aiCards)-1 && aiCards[i+1].Num > 0 && aiCards[i+1].Card-aiCards[i].Card <= 2) {
			return false
		}
	}
	return true
}

func CheckQuanBuKaoTing(cards []AICard) []int32 {
	tingInfo := make([]int32, 0, 34)
	aiCards := make([]AICard, len(cards))
	copy(aiCards, cards)
	sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })

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
		for j := lb; j < rb; {
			tingInfo = append(tingInfo, j)
			if j%MAHJONG_MASK == MAHJONG_9 {
				if j/MAHJONG_MASK == COLOR_TIAO {
					j = COLOR_OTHER*MAHJONG_MASK + MAHJONG_DONG
				} else {
					j += MAHJONG_MASK - MAHJONG_9 + MAHJONG_1
				}
			} else {
				j++
			}
		}
	}
	return tingInfo
}

//检查单调，不适用于七对，因为七对的胡牌只能单调
func CheckTingTypeForDanDiao(cards []AICard) bool {
	if len(cards) == 1 && cards[0].Num == 1 {
		return true
	}
	tings := CheckCommonTing(cards, true)
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
	if ok, _ := CheckCommonHu(aiCards, false); ok {
		return true
	}
	return false
}

//检查卡子，例：1万3万胡2万
func CheckTingTypeForQiaZi(cards []AICard) bool {
	tings := CheckCommonTing(cards, true)
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
	if ok, _ := CheckCommonHu(aiCards, false); ok {
		return true
	}
	return false
}

//检查边张，例：1万2万胡3万
func CheckTingTypeForBianZhang(cards []AICard) bool {
	tings := CheckCommonTing(cards, true)
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
	if ok, _ := CheckCommonHu(aiCards, false); ok {
		return true
	}
	return false
}

func CheckQingYiSe(chiCard []*ChiCard, cards []AICard) bool {
	color := int32(0)
	//先判断吃牌花色是否相同
	for _, v := range chiCard {
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

func CheckZiYiSe(chiCard []*ChiCard, cards []AICard) bool {
	color := int32(0)
	//先判断吃牌花色是否相同
	for _, v := range chiCard {
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

func CheckHunYiSe(chiCard []*ChiCard, cards []AICard) bool {
	color1 := int32(0)
	color2 := int32(0)
	//先判断吃牌花色
	for _, v := range chiCard {
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
	return true
}

func CheckDuiDuiHu(chiCard []*ChiCard, cards []AICard) bool {
	for _, v := range chiCard {
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

//用于四川麻将带19，仅用于MJ_HU_TYPE_COMMON
func CheckDai19(chiCard []*ChiCard, cards []AICard) bool {
	for _, v := range chiCard {
		if v.CardType == MJ_CHI_GANG ||
			v.CardType == MJ_CHI_GANG_WAN ||
			v.CardType == MJ_CHI_GANG_AN ||
			v.CardType == MJ_CHI_PENG {
			if c := v.CardId % MAHJONG_MASK; c != MAHJONG_1 && c != MAHJONG_9 {
				return false
			}
		} else if v.CardType == MJ_CHI_CHI {
			if c := v.CardId % MAHJONG_MASK; !((c == MAHJONG_1 || c == MAHJONG_9) ||
				((c == MAHJONG_2 && v.ChiPosBit != 3) || (c == MAHJONG_8 && v.ChiPosBit != 1)) ||
				((c == MAHJONG_3 && v.ChiPosBit == 1) || (c == MAHJONG_7 && v.ChiPosBit == 3))) {
				return false
			}
		}
	}

	aiCards := make([]AICard, len(cards))
	aiCardsBk := make([]AICard, len(aiCards))
	copy(aiCards, cards)
	sort.Slice(aiCards, func(i, j int) bool { return aiCards[i].Card < aiCards[j].Card })
	copy(aiCardsBk, aiCards)
	for i := 0; i < len(aiCards); i++ {
		if aiCards[i].Card%MAHJONG_MASK >= MAHJONG_DONG {
			return false
		}
		if aiCards[i].Num >= 2 {
			if c := aiCards[i].Card % MAHJONG_MASK; c != MAHJONG_1 && c != MAHJONG_9 {
				continue
			}
			aiCards[i].Num -= 2
			for j := 0; ; {
				for {
					if j >= len(aiCards) {
						return true
					}
					if aiCards[j].Num > 0 {
						break
					}
					j++
				}
				c := aiCards[j].Card % MAHJONG_MASK
				if c != MAHJONG_1 && c != MAHJONG_7 && c != MAHJONG_9 {
					break
				}
				if c == MAHJONG_9 {
					aiCards[j].Num = 0
				} else {
					if aiCards[j].Num >= 3 {
						if j >= len(aiCards)-2 || aiCards[j+1].Num == 0 || aiCards[j+2].Num == 0 ||
							aiCards[j+1].Card-aiCards[j].Card != 1 || aiCards[j+2].Card-aiCards[j+1].Card != 1 ||
							aiCards[j].Num != aiCards[j+1].Num || aiCards[j+1].Num != aiCards[j+2].Num {
							if c == MAHJONG_1 {
								aiCards[j].Num -= 3
							} else {
								break
							}
						} else {
							aiCards[j].Num = 0
							aiCards[j+1].Num = 0
							aiCards[j+2].Num = 0
						}
					} else {
						aiCards[j].Num--
						aiCards[j+1].Num--
						aiCards[j+2].Num--
					}
				}
			}
			copy(aiCards, aiCardsBk)
		}
	}
	return false
}

//用于四川麻将的将对（对对胡并且每组都有258）
func CheckDai258(chiCard []*ChiCard, cards []AICard) bool {
	for _, v := range chiCard {
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

//用于特殊七对
func GetLongNum(chiCard []*ChiCard, cards []AICard) int32 {
	if len(chiCard) > 0 {
		return 0
	}
	longNum := int32(0)
	for i := 0; i < len(cards); i++ {
		if cards[i].Num == 4 {
			longNum++
		}
	}
	return longNum
}

//全老头，用于嘉兴麻将
func CheckQuanLaoTou(chiCard []*ChiCard, cards []AICard) bool {
	for _, v := range chiCard {
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
