package mahjong

const MAHJONG_MASK int32 = 100

const (
	COLOR_WAN  int32 = 1
	COLOR_TONG int32 = 2
	COLOR_TIAO int32 = 3
)

const (
	MAHJONG_1    int32 = 1
	MAHJONG_2    int32 = 2
	MAHJONG_3    int32 = 3
	MAHJONG_4    int32 = 4
	MAHJONG_5    int32 = 5
	MAHJONG_6    int32 = 6
	MAHJONG_7    int32 = 7
	MAHJONG_8    int32 = 8
	MAHJONG_9    int32 = 9
	MAHJONG_DONG int32 = 11
)

type aiCard struct {
	card int32
	num  int32
}

func getCardNum(aiCards []aiCard) int32 {
	sl := int32(0)
	for _, v := range aiCards {
		sl += v.num
	}
	return sl	
}

func CheckHu(aiCards []aiCard) bool {
	if getCardNum(aiCards)%3 != 2 {
		return false
	}
	lenCards := len(aiCards)
	aiCardsBk := make([]aiCard, lenCards)
	copy(aiCardsBk, aiCards)
	for i := 0; i < lenCards; i++ {
		if aiCards[i].num >= 2 {
			if aiCards[i].num > 2 && aiCards[i].card%MAHJONG_MASK >= MAHJONG_DONG {
				continue
			}
			aiCards[i].num -= 2			
			for j := 0;; {
				for {
					if j >= lenCards {
						copy(aiCards, aiCardsBk)
						return true
					}
					if aiCards[j].num > 0 {
						break
					}
					j++
				}
				if aiCards[j].num >= 3 {
					aiCards[j].num -= 3
				} else {
					if j > lenCards-3 || aiCards[j].card%MAHJONG_MASK >= MAHJONG_DONG ||
						aiCards[j+1].num == 0 || aiCards[j+2].num == 0 ||
						aiCards[j].card != aiCards[j+1].card-1 || aiCards[j+1].card+1 != aiCards[j+2].card {
						break
					}
					aiCards[j].num--
					aiCards[j+1].num--
					aiCards[j+2].num--
				}
			}
			copy(aiCards, aiCardsBk)
		}
	}
	return false
}

func CheckHuForLZ(aiCards []aiCard, lzCard int32) bool {
	if getCardNum(aiCards)%3 != 2 {
		return false
	}
	lzNum := int32(0)
	lenCards := len(aiCards)
	for i := 0; i < lenCards; i++ {
		if aiCards[i].card == lzCard {
			lzNum = aiCards[i].num
			aiCards = append(aiCards[:i], aiCards[i+1:]...)
			lenCards--
			break
		}
	}
	if lzNum <= 0 {
		return CheckHu(aiCards)
	}

	aiCardsBk := make([]aiCard, lenCards)
	copy(aiCardsBk, aiCards)
	lzNumBk := lzNum

	for i := 0; i < lenCards; i++ {
		if aiCards[i].num == 1 || aiCards[i].num == 4 {
			aiCards[i].num-- //这里之所以把4张相同的牌拿出一个和癞子组合成将，主要让下面的循环尽早返回
			lzNum--
		} else if aiCards[i].num == 2 || aiCards[i].num == 3 {
			aiCards[i].num -= 2
		} else {
			continue
		}		
		for j := 0;; {
			for {
				if j >= lenCards {
					copy(aiCards, aiCardsBk)
					return true
				}
				if aiCards[j].num > 0 {
					break
				}
				j++
			}

			//以下的取牌规则相对复杂，主要为了减少最外层的循环次数，当然也可以按照"刻>顺>对>连>单"的顺序来取牌
			if aiCards[j].num >= 3 {
				aiCards[j].num -= 3
			} else if aiCards[j].num == 2 {
				if aiCards[j].card%MAHJONG_MASK >= MAHJONG_DONG || j > lenCards-2 || aiCards[j+1].num == 0 || aiCards[j+1].card-aiCards[j].card != 1 ||
					(aiCards[j+1].num != 2 && (j > lenCards-3 || aiCards[j+2].num == 0 || aiCards[j+2].card-aiCards[j+1].card != 1 || aiCards[j+2].num == 1)) {
					aiCards[j].num -= 2
					if lzNum <= 0 {
						break
					} else {
						lzNum--
					}
				} else {
					if aiCards[j+1].num == 2 && (j > lenCards-3 || aiCards[j+2].num == 0 || aiCards[j+2].card-aiCards[j+1].card != 1) {
						if lzNum <= 0 {
							break
						} else {
							lzNum--
						}
					} else {
						aiCards[j+2].num--
					}
					aiCards[j].num--
					aiCards[j+1].num--
				}
			} else {
				if j > lenCards-2 || aiCards[j].card%MAHJONG_MASK >= MAHJONG_DONG ||
					aiCards[j+1].card-aiCards[j].card > 2 ||
					(aiCards[j+1].num == 0 && (j > lenCards-3 || aiCards[j+2].num == 0 || aiCards[j+2].card-aiCards[j].card > 2)) {
					if lzNum < 2 {
						break
					} else {
						lzNum -= 2
					}
				} else if j > lenCards-3 || aiCards[j+1].num == 0 || aiCards[j+2].num == 0 ||
					aiCards[j+2].card-aiCards[j+1].card > 1 || aiCards[j+1].card-aiCards[j].card > 1 {
					if lzNum < 1 {
						break
					} else {
						lzNum--
						if aiCards[j+1].num == 0 {
							aiCards[j+2].num--
						} else {
							aiCards[j+1].num--
						}
					}
				} else {
					aiCards[j+1].num--
					aiCards[j+2].num--
				}
				aiCards[j].num--
			}
		}
		copy(aiCards, aiCardsBk)
		lzNum = lzNumBk
	}
	return false
}

func CheckTing(aiCards []aiCard) []int32 {	
	tingInfo := make([]int32, 0, 34)
	if getCardNum(aiCards)%3 != 1 {
		return tingInfo
	}
	lenCards := len(aiCards)
	aiCardsBk := make([]aiCard, lenCards)
	copy(aiCardsBk, aiCards)
	for i := 0; i < lenCards; i++ {
		subTing := make([]int32, 0, 5) //每一轮听牌个数最多三张
		if aiCards[i].num >= 2 {
			if aiCards[i].num > 2 && aiCards[i].card%MAHJONG_MASK >= MAHJONG_DONG {
				continue
			}
			aiCards[i].num -= 2						
		} else if aiCards[i].num == 1 {
			aiCards[i].num--
			subTing = append(subTing, aiCards[i].card)
		} else {
			continue
		}
		for j := 0;; {
			for {
				if j >= lenCards || aiCards[j].num > 0{					
					break
				}				
				j++
			}
			if j >= lenCards {
				tingInfo = append(tingInfo, subTing...)
				break
			}

			if aiCards[j].card%MAHJONG_MASK >= MAHJONG_DONG && (aiCards[j].num == 1 || aiCards[j].num == 4) {
				break
			} else if aiCards[j].num >= 3 {
				aiCards[j].num -= 3
			} else if j <= lenCards - 3 && 
					  aiCards[j].num >= 1 && aiCards[j+1].num >= 1 && aiCards[j+2].num >= 1 &&
					  aiCards[j].card == aiCards[j+1].card-1 && aiCards[j+1].card+1 == aiCards[j+2].card {
				aiCards[j].num--
				aiCards[j+1].num--
				aiCards[j+2].num--
			} else if aiCards[j].num == 2 {
				aiCards[j].num -= 2
				if len(subTing) > 0 {
					break
				} else {
					subTing = append(subTing, aiCards[j].card)
				}				
			} else {
				if j >= lenCards-1 {
					break
				} 
				if aiCards[j+1].num > 0 && aiCards[j].card == aiCards[j+1].card-1 {
					if len(subTing) > 0 {
						subTing = subTing[:0]
						break
					} else {
						if aiCards[j].card%MAHJONG_MASK == MAHJONG_1 {
							subTing = append(subTing, aiCards[j+1].card+1)
						} else if aiCards[j+1].card%MAHJONG_MASK == MAHJONG_9 {
							subTing = append(subTing, aiCards[j].card-1)
						} else {
							subTing = append(subTing, aiCards[j].card-1, aiCards[j+1].card+1)
						}
					}
				} //todo
			}
		}
		copy(aiCards, aiCardsBk)
	}
	return tingInfo
}