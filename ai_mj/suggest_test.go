package ai_mj

import (
	"fmt"
	"testing"
)

func TestSuggestCard(t *testing.T) {
	aiCards := []AICard{
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_1, 3},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_2, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_3, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_4, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_5, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_6, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_7, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_8, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_9, 3},
	}	
	fmt.Println("SuggestCard: ", SuggestCard(aiCards))
}
