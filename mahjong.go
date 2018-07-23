package mahjong

import (
//	"fmt"
	"sort"
)

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
	MAHJONG_DONG      int32 = 11
	MAHJONG_NAN       int32 = 12
	MAHJONG_XI        int32 = 13
	MAHJONG_BEI       int32 = 14
	MAHJONG_HONGZHONG int32 = 15
	MAHJONG_LVFA      int32 = 16
	MAHJONG_BAI       int32 = 17	
	MAHJONG_ANY	 int32 = 50
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

func CheckHu(cards []aiCard) bool {
	if getCardNum(cards)%3 != 2 {
		return false
	}
	lenCards := len(cards)
	aiCards := make([]aiCard, lenCards)
	aiCardsBk := make([]aiCard, lenCards)
	copy(aiCards, cards)
	sort.Slice(aiCards, func (i, j int) bool {return aiCards[i].card < aiCards[j].card})
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

func CheckHuForLZ(cards []aiCard, lzCard int32) bool {
	if getCardNum(cards)%3 != 2 {
		return false
	}
	lzNum := int32(0)
	lenCards := len(cards)
	aiCards := make([]aiCard, lenCards)
	copy(aiCards, cards)
	sort.Slice(aiCards, func (i, j int) bool {return aiCards[i].card < aiCards[j].card})
	for i := 0; i < lenCards; i++ {
		if aiCards[i].card == lzCard {
			lzNum = aiCards[i].num
			if i == lenCards - 1 {
				aiCards = aiCards[:i]	
			} else {
				aiCards = append(aiCards[:i], aiCards[i+1:]...)
			}			
			lenCards--
			break
		}
	}
	if lzNum <= 0 {
		return CheckHu(cards)
	}

	aiCardsBk := make([]aiCard, lenCards)
	copy(aiCardsBk, aiCards)
	lzNumBk := lzNum

	for i := 0; i < lenCards; i++ {
		if aiCards[i].num >= 2 {
			aiCards[i].num -= 2
		} else if aiCards[i].num == 1 {
			aiCards[i].num--
			lzNum--
		} else {
			continue
		}		
		for j := 0;; {
			for {
				if j >= lenCards {					
					return true
				}
				if aiCards[j].num > 0 {
					break
				}
				j++
			}

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

func CheckTing(cards []aiCard) []int32 {	
	tingInfo := make([]int32, 0, 34)
	if getCardNum(cards)%3 != 1 {
		return tingInfo
	}
	lenCards := len(cards)
	aiCards := make([]aiCard, lenCards)
	aiCardsBk := make([]aiCard, lenCards)
	copy(aiCards, cards)
	sort.Slice(aiCards, func (i, j int) bool {return aiCards[i].card < aiCards[j].card})
	copy(aiCardsBk, aiCards)
	for i := 0; i < lenCards; i++ {
		subTing := make([]int32, 0, 5) //每一轮听牌个数最多三张
		group := make(map[int32]int32)
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
				break
			}			
			if aiCards[j].card%MAHJONG_MASK >= MAHJONG_DONG && (aiCards[j].num == 1 || aiCards[j].num == 4) {
				subTing = subTing[:0]
				break
			} else if aiCards[j].num >= 3 {
				aiCards[j].num -= 3
				group[aiCards[j].card] = 0
			} else if j <= lenCards - 3 && 
					  aiCards[j].num >= 1 && aiCards[j+1].num >= 1 && aiCards[j+2].num >= 1 &&
					  aiCards[j].card == aiCards[j+1].card-1 && aiCards[j+1].card+1 == aiCards[j+2].card {
				aiCards[j].num--
				aiCards[j+1].num--
				aiCards[j+2].num--
				if _, ok := group[aiCards[j].card]; !ok {
					group[aiCards[j].card] = 1
				} else {
					group[aiCards[j].card]++
				}
			} else if aiCards[j].num == 2 {
				aiCards[j].num -= 2
				if len(subTing) > 0 {
					subTing = subTing[:0]
					break
				} else {
					subTing = append(subTing, aiCards[j].card)
				}				
			} else {
				if len(subTing) > 0 || j > lenCards-2 || aiCards[j+1].card-aiCards[j].card > 2 ||
				   (aiCards[j+1].num == 0 && (j > lenCards-3 || aiCards[j+2].num == 0 || aiCards[j+2].card-aiCards[j].card > 2)) {
				   	subTing = subTing[:0]
					break
				} 
				if aiCards[j+1].num > 0 {
					if aiCards[j].card == aiCards[j+1].card-1 {
						if aiCards[j].card%MAHJONG_MASK == MAHJONG_1 {
							subTing = append(subTing, aiCards[j+1].card+1)
						} else if aiCards[j+1].card%MAHJONG_MASK == MAHJONG_9 {
							subTing = append(subTing, aiCards[j].card-1)
						} else {
							subTing = append(subTing, aiCards[j].card-1, aiCards[j+1].card+1)
						}	
					} else {
						subTing = append(subTing, aiCards[j].card+1)
					}
					aiCards[j+1].num--
				} else {
					subTing = append(subTing, aiCards[j].card+1)
					aiCards[j+2].num--
				}
				aiCards[j].num--
			}
		}
		if len(subTing) > 0 {
			if subTing[0]%MAHJONG_MASK >= MAHJONG_4 && subTing[0]%MAHJONG_MASK <= MAHJONG_9 && group[subTing[0]-2] > 0 {
				subTing = append(subTing, subTing[0]-3) //若胡4、5、6、7、8、9、47、58、69, 检测是否胡14、25、36、47、58、69、147、258、369
			}  
			if subTing[0]%MAHJONG_MASK >= MAHJONG_7 && subTing[0]%MAHJONG_MASK <= MAHJONG_9 && group[subTing[0]-5] > 0 {
				subTing = append(subTing, subTing[0]-6)	//继续检测是否胡147、258、369
			}
			tingInfo = append(tingInfo, subTing...)
		}
		copy(aiCards, aiCardsBk)
	}
	if len(tingInfo) > 0 {
		sort.Slice(tingInfo, func(i,j int) bool {return tingInfo[i] < tingInfo[j]})
	}
	return tingInfo
}

func CheckTingForLZ(cards []aiCard, lzCard int32) []int32 {
	tingInfo := make([]int32, 0, 34)
	if getCardNum(cards)%3 != 1 {
		return tingInfo
	}

	lzNum := int32(0)
	lenCards := len(cards)
	aiCards := make([]aiCard, lenCards)
	copy(aiCards, cards)
	sort.Slice(aiCards, func (i, j int) bool {return aiCards[i].card < aiCards[j].card})
	for i := 0; i < lenCards; i++ {
		if aiCards[i].card == lzCard {
			lzNum = aiCards[i].num
			if i == lenCards-1 {
				aiCards = aiCards[:i]
			} else {
				aiCards = append(aiCards[:i], aiCards[i+1:]...)
			}			
			lenCards--
			break
		}
	}
	if lzNum <= 0 {
		return CheckTing(cards)
	}

	aiCardsBk := make([]aiCard, lenCards)
	copy(aiCardsBk, aiCards)
	lzNumBk := lzNum	

	for i := 0; i < lenCards; i++ {		
		subTing := make([]int32, 0, 10) //每一轮最多听8张		
		group := make(map[int32]int32)
		if aiCards[i].num >= 2 {			
			aiCards[i].num -= 2						
		} else if aiCards[i].num == 1 {
			aiCards[i].num--
			lzNum--		
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
				break
			}				
			if aiCards[j].num >= 3 {
				aiCards[j].num -= 3
				group[aiCards[j].card] = 0
			} else if j <= lenCards - 3 && 
					  aiCards[j].num >= 1 && aiCards[j+1].num >= 1 && aiCards[j+2].num >= 1 &&
					  aiCards[j].card == aiCards[j+1].card-1 && aiCards[j+1].card+1 == aiCards[j+2].card {
				aiCards[j].num--
				aiCards[j+1].num--
				aiCards[j+2].num--
				if _, ok := group[aiCards[j].card]; !ok {
					group[aiCards[j].card] = 1
				} else {
					group[aiCards[j].card]++
				}
			} else {				
				if lzNum < 0 {
					subTing = subTing[:0]
					break
				} 
				if aiCards[j].num == 2 {								
					if lzNum >= 0 {
						lzNum--						
					} 
					if group[aiCards[j].card-2] > 0 {
						subTing = append(subTing, aiCards[j].card-3) //检测是否胡前一张
					}
					subTing = append(subTing, aiCards[j].card)
					aiCards[j].num -= 2					
				} else {								
					if j < lenCards-1 && aiCards[j+1].num > 0 && aiCards[j+1].card-aiCards[j].card < 3 {						
						if lzNum >= 0 {
							lzNum--						
						} 
						if aiCards[j+1].card-aiCards[j].card == 1 {
							if aiCards[j].card%MAHJONG_MASK == MAHJONG_1 {								
								subTing = append(subTing, aiCards[j+1].card+1)
							} else if aiCards[j+1].card%MAHJONG_MASK == MAHJONG_9 {
								if group[aiCards[j].card-3] > 0 {
									subTing = append(subTing, aiCards[j].card-4)	//若胡7，检测是否胡47
								}
								if group[aiCards[j].card-6] > 0 {
									subTing = append(subTing, aiCards[j].card-7)	//继续检测是否胡147
								}
								subTing = append(subTing, aiCards[j].card-1)
							} else {
								if aiCards[j].card%MAHJONG_MASK > MAHJONG_4 && aiCards[j].card%MAHJONG_MASK < MAHJONG_8 && group[aiCards[j].card-3] > 0 {
									subTing = append(subTing, aiCards[j].card-4) //若胡47、58、69，检测是否胡147、、258、369
								}								
								subTing = append(subTing, aiCards[j].card-1, aiCards[j+1].card+1)
							}
						} else {
							subTing = append(subTing, aiCards[j].card+1)
						}
						aiCards[j+1].num--
					} else if j < lenCards - 2 && aiCards[j+1].num == 0 && aiCards[j+2].num != 0 && aiCards[j+2].card-aiCards[j].card < 3 {						
						if lzNum >= 0 {
							lzNum--						
						}
						subTing = append(subTing, aiCards[j].card+1) 
						aiCards[j+2].num--
					} else {						
						if lzNum <= 0 {
							subTing = subTing[:0]
							break
						}
						lzNum -= 2
						if aiCards[j].card%MAHJONG_MASK == MAHJONG_1 {
							subTing = append(subTing, aiCards[j].card, aiCards[j].card+1, aiCards[j].card+2)
						} else if aiCards[j].card%MAHJONG_MASK == MAHJONG_2 {
							subTing = append(subTing, aiCards[j].card-1, aiCards[j].card, aiCards[j].card+1, aiCards[j].card+2)
						} else if aiCards[j].card%MAHJONG_MASK == MAHJONG_8 {
							subTing = append(subTing, aiCards[j].card-2, aiCards[j].card-1, aiCards[j].card, aiCards[j].card+1)
						} else if aiCards[j].card%MAHJONG_MASK == MAHJONG_9 {
							subTing = append(subTing, aiCards[j].card-2, aiCards[j].card-1, aiCards[j].card)								
						} else {
							subTing = append(subTing, aiCards[j].card-2, aiCards[j].card-1, aiCards[j].card, aiCards[j].card+1, aiCards[j].card+2)
						}
					}
					aiCards[j].num--
				}
			}
		}
		if len(subTing) > 0 {			
			tingInfo = append(tingInfo, subTing...)
		}
		copy(aiCards, aiCardsBk)
		lzNum = lzNumBk
	} 

	//检查是否听所有牌
	huCards := make([]aiCard, len(cards)+1)
	copy(huCards, cards)
	if CheckHuForLZ(append(huCards, aiCard{MAHJONG_ANY, 1}), lzCard) == true {		
		tingInfo = append(tingInfo, MAHJONG_ANY)
		return tingInfo[len(tingInfo)-1:]
	}

	if len(tingInfo) > 0 {
		sort.Slice(tingInfo, func(i,j int) bool {return tingInfo[i] < tingInfo[j]})
		tmpTing := make([]int32, 0, len(tingInfo))
		tmpTing = append(tmpTing, tingInfo[0])
		for i := 1; i < len(tingInfo); i++ {
			if tingInfo[i] != tmpTing[len(tmpTing)-1] {
				tmpTing = append(tmpTing, tingInfo[i])
			}
		}
		tingInfo = tmpTing
	}
	return tingInfo	
}