package ai_mj

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

type AICard struct {
	Card int32
	Num  int32
}

func getCardNum(aiCards []AICard) int32 {
	sl := int32(0)
	for _, v := range aiCards {
		sl += v.Num
	}
	return sl	
}

func CheckHu(cards []AICard) bool {
	if getCardNum(cards)%3 != 2 {
		return false
	}
	lenCards := len(cards)
	aiCards := make([]AICard, lenCards)
	aiCardsBk := make([]AICard, lenCards)
	copy(aiCards, cards)
	sort.Slice(aiCards, func (i, j int) bool {return aiCards[i].Card < aiCards[j].Card})
	copy(aiCardsBk, aiCards)
	for i := 0; i < lenCards; i++ {
		if aiCards[i].Num >= 2 {
			if aiCards[i].Num > 2 && aiCards[i].Card%MAHJONG_MASK >= MAHJONG_DONG {
				continue
			}
			aiCards[i].Num -= 2			
			for j := 0;; {
				for {
					if j >= lenCards {						
						return true
					}
					if aiCards[j].Num > 0 {
						break
					}
					j++
				}
				if aiCards[j].Num >= 3 {
					aiCards[j].Num -= 3
				} else {
					if j > lenCards-3 || aiCards[j].Card%MAHJONG_MASK >= MAHJONG_DONG ||
						aiCards[j+1].Num == 0 || aiCards[j+2].Num == 0 ||
						aiCards[j].Card != aiCards[j+1].Card-1 || aiCards[j+1].Card+1 != aiCards[j+2].Card {
						break
					}
					aiCards[j].Num--
					aiCards[j+1].Num--
					aiCards[j+2].Num--
				}
			}
			copy(aiCards, aiCardsBk)
		}
	}
	return false
}

func CheckHuForLZ(cards []AICard, lzCard int32) bool {
	if getCardNum(cards)%3 != 2 {
		return false
	}
	lzNum := int32(0)
	lenCards := len(cards)
	aiCards := make([]AICard, lenCards)
	copy(aiCards, cards)
	sort.Slice(aiCards, func (i, j int) bool {return aiCards[i].Card < aiCards[j].Card})
	for i := 0; i < lenCards; i++ {
		if aiCards[i].Card == lzCard {
			lzNum = aiCards[i].Num
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

	aiCardsBk := make([]AICard, lenCards)
	copy(aiCardsBk, aiCards)
	lzNumBk := lzNum

	for i := 0; i < lenCards; i++ {
		if aiCards[i].Num >= 2 {
			aiCards[i].Num -= 2
		} else if aiCards[i].Num == 1 {
			aiCards[i].Num--
			lzNum--
		} else {
			continue
		}		
		for j := 0;; {
			for {
				if j >= lenCards {					
					return true
				}
				if aiCards[j].Num > 0 {
					break
				}
				j++
			}
			
			if aiCards[j].Num >= 3 {
				aiCards[j].Num -= 3
			} else if j <= lenCards - 3 && 
					  aiCards[j].Num >= 1 && aiCards[j+1].Num >= 1 && aiCards[j+2].Num >= 1 &&
					  aiCards[j].Card == aiCards[j+1].Card-1 && aiCards[j+1].Card+1 == aiCards[j+2].Card &&
					  (aiCards[j].Num == 1 || aiCards[j+1].Num == 2 || aiCards[j+2].Num > 1) {
				aiCards[j].Num--
				aiCards[j+1].Num--
				aiCards[j+2].Num--
			} else if aiCards[j].Num == 2 {
				aiCards[j].Num -= 2
				if lzNum <= 0 {
					break
				} else {
					lzNum--
				}
			} else {
				if j > lenCards-2 || aiCards[j].Card%MAHJONG_MASK >= MAHJONG_DONG ||
					aiCards[j+1].Card-aiCards[j].Card > 2 ||
					(aiCards[j+1].Num == 0 && (j > lenCards-3 || aiCards[j+2].Num == 0 || aiCards[j+2].Card-aiCards[j].Card > 2)) {
					if lzNum < 2 {
						break
					} else {
						lzNum -= 2
					}
				} else {
					if lzNum < 1 {
						break
					} else {
						lzNum--
						if aiCards[j+1].Num == 0 {
							aiCards[j+2].Num--
						} else {
							aiCards[j+1].Num--
						}
					}
				}
				aiCards[j].Num--
			}
		}
		copy(aiCards, aiCardsBk)
		lzNum = lzNumBk
	}
	return false
}

func CheckTing(cards []AICard) []int32 {	
	tingInfo := make([]int32, 0, 34)
	if getCardNum(cards)%3 != 1 {
		return tingInfo
	}
	lenCards := len(cards)
	aiCards := make([]AICard, lenCards)
	aiCardsBk := make([]AICard, lenCards)
	copy(aiCards, cards)
	sort.Slice(aiCards, func (i, j int) bool {return aiCards[i].Card < aiCards[j].Card})
	copy(aiCardsBk, aiCards)
	for i := 0; i < lenCards; i++ {
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
		for j := 0;; {
			for {
				if j >= lenCards || aiCards[j].Num > 0{					
					break
				}				
				j++
			}
			if j >= lenCards {
				break
			}			
			if aiCards[j].Card%MAHJONG_MASK >= MAHJONG_DONG && (aiCards[j].Num == 1 || aiCards[j].Num == 4) {
				subTing = subTing[:0]
				break
			} else if aiCards[j].Num >= 3 {
				aiCards[j].Num -= 3				
			} else if j <= lenCards - 3 && 
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
				if len(subTing) > 0 || j > lenCards-2 || aiCards[j+1].Card-aiCards[j].Card > 2 ||
				   (aiCards[j+1].Num == 0 && (j > lenCards-3 || aiCards[j+2].Num == 0 || aiCards[j+2].Card-aiCards[j].Card > 2)) {
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
			}  
			if subTing[0]%MAHJONG_MASK >= MAHJONG_7 && subTing[0]%MAHJONG_MASK <= MAHJONG_9 && seq[subTing[0]-5] > 0 {
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

func CheckTingForLZ(cards []AICard, lzCard int32) []int32 {
	tingInfo := make([]int32, 0, 34)
	if getCardNum(cards)%3 != 1 {
		return tingInfo
	}

	lzNum := int32(0)
	lenCards := len(cards)
	aiCards := make([]AICard, lenCards)
	copy(aiCards, cards)
	sort.Slice(aiCards, func (i, j int) bool {return aiCards[i].Card < aiCards[j].Card})
	for i := 0; i < lenCards; i++ {
		if aiCards[i].Card == lzCard {
			lzNum = aiCards[i].Num
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

	//检查是否听所有牌
	huCards := make([]AICard, len(cards)+1)
	copy(huCards, cards)
	if CheckHuForLZ(append(huCards, AICard{MAHJONG_ANY, 1}), lzCard) == true {		
		tingInfo = append(tingInfo, MAHJONG_ANY)
		return tingInfo
	}

	aiCardsBk := make([]AICard, lenCards)
	copy(aiCardsBk, aiCards)
	lzNumBk := lzNum	

	for i := 0; i < lenCards; i++ {		
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
		for j := 0;; {
			for {
				if j >= lenCards || aiCards[j].Num > 0{					
					break
				}				
				j++
			}			
			if j >= lenCards {
				break
			}				
			if aiCards[j].Num >= 3 {
				aiCards[j].Num -= 3				
			} else if j <= lenCards - 3 && 
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
					if j < lenCards-1 && aiCards[j+1].Num > 0 && aiCards[j+1].Card-aiCards[j].Card < 3 {						
						if lzNum >= 0 {
							lzNum--						
						} 
						if aiCards[j+1].Card-aiCards[j].Card == 1 {
							if aiCards[j].Card%MAHJONG_MASK == MAHJONG_1 {								
								subTing = append(subTing, aiCards[j+1].Card+1)
							} else if aiCards[j+1].Card%MAHJONG_MASK == MAHJONG_9 {
								if seq[aiCards[j].Card-3] > 0 {
									subTing = append(subTing, aiCards[j].Card-4)	//若胡7，检测是否胡47
								}
								if seq[aiCards[j].Card-6] > 0 {
									subTing = append(subTing, aiCards[j].Card-7)	//继续检测是否胡147
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
					} else if j < lenCards - 2 && aiCards[j+1].Num == 0 && aiCards[j+2].Num != 0 && aiCards[j+2].Card-aiCards[j].Card < 3 {						
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
		}
		copy(aiCards, aiCardsBk)
		lzNum = lzNumBk
	} 	

	if len(tingInfo) > 0 {
		//对结果进行排序除重
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