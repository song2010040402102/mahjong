package ai_mj

import (	
	"math"		
)

const ai_fact_num int32 = 10 //数量因子
const ai_fact_dis int32 = 10 //距离因子
const ai_fact_pos int32 = 1  //位置因子

const ai_rel_edg_val float32 = 5.0 //边界值，小于边界值的牌可认为没有关系的牌

//获取牌的评估值
func getCardValue(cards []AICard) map[int32]float32 {		
	//对每张牌进行价值评估，由位置、距离、数量三个指标构成
	valCards := map[int32]float32{}
	for i := 0; i < len(cards); i++ {
		val := float32(0)
		if cards[i].Num > 1 {
			nv := float32(ai_fact_num * (cards[i].Num - 1)) //数量评估值
			val += nv
		}
		if c := cards[i].Card % MAHJONG_MASK; c >= MAHJONG_1 && c <= MAHJONG_9 {
			if i < len(cards)-1 {
				d := cards[i+1].Card - cards[i].Card
				if d == 1 && (cards[i].Card%MAHJONG_MASK == MAHJONG_1 || cards[i+1].Card%MAHJONG_MASK == MAHJONG_9) {
					d = 2
				}
				if d > 0 && d <= MAHJONG_9-MAHJONG_1 {
					if d > 2 {
						d *= 2 //这里可避免22588牌型和12牌型先打1而不打5
					}
					dv := float32(ai_fact_dis) / float32(d) //距离评估值
					val += dv                               //右距离值
					valCards[cards[i+1].Card] = dv        //左距离值
				}
			}
			pv := float32(float64(ai_fact_pos) / (math.Abs(float64(cards[i].Card%MAHJONG_MASK-MAHJONG_5)) + 1)) //位置评估值
			val += pv
		}		
		_, ok := valCards[cards[i].Card]
		if ok {
			valCards[cards[i].Card] += val
		} else {
			valCards[cards[i].Card] = val
		}
	}
	return valCards
}

func SuggestCard(cards []AICard) int32 {
	
	return 0
}