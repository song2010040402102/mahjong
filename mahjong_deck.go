package logic

import (
	"github.com/astaxie/beego/logs"
	"util"
)

const (
	COLOR_WAN   int32 = 1
	COLOR_TONG  int32 = 2
	COLOR_TIAO  int32 = 3
	COLOR_OTHER int32 = 4
)

const (
	MAHJONG_1 int32 = 1
	MAHJONG_2 int32 = 2
	MAHJONG_3 int32 = 3
	MAHJONG_4 int32 = 4
	MAHJONG_5 int32 = 5
	MAHJONG_6 int32 = 6
	MAHJONG_7 int32 = 7
	MAHJONG_8 int32 = 8
	MAHJONG_9 int32 = 9

	MAHJONG_DONG      int32 = 11
	MAHJONG_NAN       int32 = 12
	MAHJONG_XI        int32 = 13
	MAHJONG_BEI       int32 = 14
	MAHJONG_HONGZHONG int32 = 15
	MAHJONG_LVFA      int32 = 16
	MAHJONG_BAI       int32 = 17

	MAHJONG_SPRING int32 = 20
	MAHJONG_SUMMER int32 = 21
	MAHJONG_AUTUMN int32 = 22
	MAHJONG_WINTER int32 = 23
	MAHJONG_MEI    int32 = 24
	MAHJONG_LAN    int32 = 25
	MAHJONG_ZHU    int32 = 26
	MAHJONG_JU     int32 = 27

	MAHJONG_MOUSE int32 = 31
	MAHJONG_GOD   int32 = 32
	MAHJONG_CAT   int32 = 33
	MAHJONG_POT   int32 = 34
	MAHJONG_DA    int32 = 35

	MAHJONG_ANY int32 = 50
)

//几乎包含可能的麻将牌设置，若存在麻将规则限定牌中没有1万或9条等，需把设置扩大至牌种类全集
const (
	WITH_WAN    uint64 = 0x01 << 0
	WITH_TONG   uint64 = 0x01 << 1
	WITH_TIAO   uint64 = 0x01 << 2
	WITH_OTHER  uint64 = 0x01 << 3
	WITH_1      uint64 = 0x01 << 4
	WITH_2      uint64 = 0x01 << 5
	WITH_3      uint64 = 0x01 << 6
	WITH_4      uint64 = 0x01 << 7
	WITH_5      uint64 = 0x01 << 8
	WITH_6      uint64 = 0x01 << 9
	WITH_7      uint64 = 0x01 << 10
	WITH_8      uint64 = 0x01 << 11
	WITH_9      uint64 = 0x01 << 12
	WITH_FENG   uint64 = 0x01 << 13
	WITH_DONG   uint64 = 0x01 << 14
	WITH_NAN    uint64 = 0x01 << 15
	WITH_XI     uint64 = 0x01 << 16
	WITH_BEI    uint64 = 0x01 << 17
	WITH_JIAN   uint64 = 0x01 << 18
	WITH_ZHONG  uint64 = 0x01 << 19
	WITH_FA     uint64 = 0x01 << 20
	WITH_BAI    uint64 = 0x01 << 21
	WITH_HUA    uint64 = 0x01 << 22
	WITH_SPRING uint64 = 0x01 << 23
	WITH_SUMMER uint64 = 0x01 << 24
	WITH_AUTUMN uint64 = 0x01 << 25
	WITH_WINTER uint64 = 0x01 << 26
	WITH_MEI    uint64 = 0x01 << 27
	WITH_LAN    uint64 = 0x01 << 28
	WITH_ZHU    uint64 = 0x01 << 29
	WITH_JU     uint64 = 0x01 << 30
	WITH_END    uint64 = 0x01 << 31

	WITH_SUZHOU uint64 = 0x01 << 63
)

const (
	MAHJONG_MASK int32 = 100
	MAHJONG_LZ   int32 = 1000 //癞子牌标记
)

func GetMahjongSeq(card int32) int32 {
	return card % MAHJONG_MASK
}

func GetMahjongColor(card int32) int32 {
	return card / MAHJONG_MASK
}

type MahjongDeck struct {
	cards []int32
	flag  uint64
}

func NewMahjongDeck(umask uint64) *MahjongDeck {
	deck := &MahjongDeck{}
	for i := uint64(0); i < 64; i++ {
		if flag := uint64(0x01 << i); flag != WITH_END {
			deck.flag |= flag
		} else {
			break
		}
	}
	deck.flag &^= umask
	if deck.flag&WITH_OTHER == 0 {
		deck.flag &^= WITH_FENG
		deck.flag &^= WITH_JIAN
		deck.flag &^= WITH_HUA
	}
	if deck.flag&WITH_FENG == 0 {
		deck.flag &^= WITH_DONG
		deck.flag &^= WITH_NAN
		deck.flag &^= WITH_XI
		deck.flag &^= WITH_BEI
	}
	if deck.flag&WITH_JIAN == 0 {
		deck.flag &^= WITH_ZHONG
		deck.flag &^= WITH_FA
		deck.flag &^= WITH_BAI
	}
	if deck.flag&WITH_HUA == 0 {
		deck.flag &^= WITH_SPRING
		deck.flag &^= WITH_SUMMER
		deck.flag &^= WITH_AUTUMN
		deck.flag &^= WITH_WINTER
		deck.flag &^= WITH_MEI
		deck.flag &^= WITH_LAN
		deck.flag &^= WITH_ZHU
		deck.flag &^= WITH_JU
	}
	if deck.flag&(WITH_DONG|WITH_NAN|WITH_XI|WITH_BEI) == 0 {
		deck.flag &^= WITH_FENG
	}
	if deck.flag&(WITH_ZHONG|WITH_FA|WITH_BAI) == 0 {
		deck.flag &^= WITH_JIAN
	}
	if deck.flag&(WITH_SPRING|WITH_SUMMER|WITH_AUTUMN|WITH_WINTER|WITH_MEI|WITH_LAN|WITH_ZHU|WITH_JU) == 0 {
		deck.flag &^= WITH_HUA
	}
	if deck.flag&(WITH_FENG|WITH_JIAN|WITH_HUA) == 0 {
		deck.flag &^= WITH_OTHER
	}
	if umask&(WITH_SUZHOU) != 0 {
		deck.flag |= WITH_SUZHOU
	}
	deck.cards = make([]int32, 0, 160)
	cards := deck.GetCardKind(true)
	for _, v := range cards {
		if c := v % MAHJONG_MASK; c >= MAHJONG_SPRING && c <= MAHJONG_JU {
			deck.cards = append(deck.cards, v)
		} else if c := v % MAHJONG_MASK; c >= MAHJONG_MOUSE && c <= MAHJONG_POT {
			deck.cards = append(deck.cards, v)
		} else {
			for i := 0; i < 4; i++ {
				deck.cards = append(deck.cards, v)
			}
		}
	}
	return deck
}

func (ck *MahjongDeck) GetCards() []int32 {
	return ck.cards
}

func (ck *MahjongDeck) GetFlag() uint64 {
	return ck.flag
}

func (ck *MahjongDeck) GetColors(hasOther bool) []int32 {
	colors := []int32{}
	if ck.flag&WITH_WAN != 0 {
		colors = append(colors, COLOR_WAN)
	}
	if ck.flag&WITH_TONG != 0 {
		colors = append(colors, COLOR_TONG)
	}
	if ck.flag&WITH_TIAO != 0 {
		colors = append(colors, COLOR_TIAO)
	}
	if hasOther && ck.flag&WITH_OTHER != 0 {
		colors = append(colors, COLOR_OTHER)
	}
	return colors
}

func (ck *MahjongDeck) GetCardKind(hasHua bool) []int32 {
	color := uint32(0)
	cards := []int32{}
	if ck.flag&WITH_WAN > 0 {
		color |= 0x01 << uint32(COLOR_WAN)
	}
	if ck.flag&WITH_TONG > 0 {
		color |= 0x01 << uint32(COLOR_TONG)
	}
	if ck.flag&WITH_TIAO > 0 {
		color |= 0x01 << uint32(COLOR_TIAO)
	}
	for i := COLOR_WAN; i <= COLOR_TIAO; i++ {
		if color&(0x01<<uint32(i)) != 0 {
			if ck.flag&WITH_1 > 0 {
				cards = append(cards, i*MAHJONG_MASK+MAHJONG_1)
			}
			if ck.flag&WITH_2 > 0 {
				cards = append(cards, i*MAHJONG_MASK+MAHJONG_2)
			}
			if ck.flag&WITH_3 > 0 {
				cards = append(cards, i*MAHJONG_MASK+MAHJONG_3)
			}
			if ck.flag&WITH_4 > 0 {
				cards = append(cards, i*MAHJONG_MASK+MAHJONG_4)
			}
			if ck.flag&WITH_5 > 0 {
				cards = append(cards, i*MAHJONG_MASK+MAHJONG_5)
			}
			if ck.flag&WITH_6 > 0 {
				cards = append(cards, i*MAHJONG_MASK+MAHJONG_6)
			}
			if ck.flag&WITH_7 > 0 {
				cards = append(cards, i*MAHJONG_MASK+MAHJONG_7)
			}
			if ck.flag&WITH_8 > 0 {
				cards = append(cards, i*MAHJONG_MASK+MAHJONG_8)
			}
			if ck.flag&WITH_9 > 0 {
				cards = append(cards, i*MAHJONG_MASK+MAHJONG_9)
			}
		}
	}
	if ck.flag&WITH_DONG > 0 {
		cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_DONG)
	}
	if ck.flag&WITH_NAN > 0 {
		cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_NAN)
	}
	if ck.flag&WITH_XI > 0 {
		cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_XI)
	}
	if ck.flag&WITH_BEI > 0 {
		cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_BEI)
	}
	if ck.flag&WITH_ZHONG > 0 {
		cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_HONGZHONG)
	}
	if ck.flag&WITH_FA > 0 {
		cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_LVFA)
	}
	if ck.flag&WITH_BAI > 0 {
		cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_BAI)
	}
	if hasHua {
		if ck.flag&WITH_SPRING > 0 {
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_SPRING)
		}
		if ck.flag&WITH_SUMMER > 0 {
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_SUMMER)
		}
		if ck.flag&WITH_AUTUMN > 0 {
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_AUTUMN)
		}
		if ck.flag&WITH_WINTER > 0 {
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_WINTER)
		}
		if ck.flag&WITH_MEI > 0 {
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_MEI)
		}
		if ck.flag&WITH_LAN > 0 {
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_LAN)
		}
		if ck.flag&WITH_ZHU > 0 {
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_ZHU)
		}
		if ck.flag&WITH_JU > 0 {
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_JU)
		}
		if ck.flag&WITH_SUZHOU > 0 {
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_MOUSE)
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_GOD)
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_CAT)
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_POT)
			cards = append(cards, COLOR_OTHER*MAHJONG_MASK+MAHJONG_DA)
		}
	}
	return cards
}

func (ck *MahjongDeck) Shuffle(roomId int64, rule IMahjong, holdCards map[int32][]int32) {
	util.RandShuffle(ck.cards, 0, len(ck.cards))
	logs.Info("[MahjongDeck]Shuffle, roomId: ", roomId, "cards: ", ck.cards)
	if debug := GetDebugMJRuleByRoom(roomId); debug != nil {
		ck.setCards(debug.Cards, debug.Restc, rule.GetPlayerLimit(), rule.GetPlayerCardNum())
	} else if len(holdCards) > 0 {
		ck.setCards(holdCards, []int32{}, rule.GetPlayerLimit(), rule.GetPlayerCardNum())
	}
	logs.Info("[MahjongDeck]Shuffle, roomId: ", roomId, "newCards: ", ck.cards)
}

func (ck *MahjongDeck) setCards(holdCards map[int32][]int32, restCards []int32, playerNum int32, cardNum int32) {
	if len(holdCards) == 0 && len(restCards) == 0 {
		return
	}
	newCards := make([]int32, len(ck.cards))
	for n := 0; n < 2; n++ {
		for k, v := range holdCards {
			if k <= 0 {
				continue
			}

			start := int32(0)
			rest := cardNum - int32(len(v))
			if k == 1 {
				rest++
				if int32(len(v)) > cardNum+1 {
					v = v[:cardNum+1]
				}
			} else {
				start = (k-1)*cardNum + 1
				if int32(len(v)) > cardNum {
					v = v[:cardNum]
				}
			}

			//设定指定的手牌
			if n == 0 {
				copy(newCards[start:], v)
				for _, vv := range v {
					ck.cards = util.RemoveSliceElem(ck.cards, vv, false).([]int32)
				}
			}

			//随机其余手牌
			if n == 1 {
				for i := int32(0); i < rest; i++ {
					newCards[start+int32(len(v))+i] = ck.cards[0]
					ck.cards = ck.cards[1:]
				}
			}
		}
	}

	//指定剩余的牌
	if len(restCards) > 0 {
		copy(newCards[playerNum*cardNum+1:], restCards)
		for _, v := range restCards {
			ck.cards = util.RemoveSliceElem(ck.cards, v, false).([]int32)
		}
	}

	//复制桌面上剩余的牌
	copy(newCards[playerNum*cardNum+1+int32(len(restCards)):], ck.cards)

	ck.cards = newCards
}
