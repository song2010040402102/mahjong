package ai_mj

import (
	"fmt"
	"testing"
)

func TestCheckHu(t *testing.T) {
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
	for j := 0; j < len(aiCards); j++ {
		aiCards[j].Num++
		fmt.Println(j, CheckHu(aiCards))
		aiCards[j].Num--
	}
}

func TestCheckHuForLZ(t *testing.T) {
	aiCards := []AICard{
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_1, 3},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_2, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_3, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_4, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_5, 2},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_6, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_7, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_8, 1},
		{COLOR_WAN*MAHJONG_MASK + MAHJONG_9, 3},
		{COLOR_TIAO*MAHJONG_MASK + MAHJONG_1, 0},
	}
	aiCards[len(aiCards)-1].Num = 1
	for j1 := 0; j1 < len(aiCards)-1; j1++ {
		aiCards[j1].Num--
		fmt.Println(j1, CheckHuForLZ(aiCards, COLOR_TIAO*MAHJONG_MASK+MAHJONG_1))
		aiCards[j1].Num++
	}

	aiCards[len(aiCards)-1].Num = 2
	for j1 := 0; j1 < len(aiCards)-1; j1++ {
		if aiCards[j1].Num == 0 {
			continue
		}
		aiCards[j1].Num--
		for j2 := j1; j2 < len(aiCards)-1; j2++ {
			if aiCards[j2].Num == 0 {
				continue
			}
			aiCards[j2].Num--
			fmt.Println(j1, j2, CheckHuForLZ(aiCards, COLOR_TIAO*MAHJONG_MASK+MAHJONG_1))
			aiCards[j2].Num++
		}
		aiCards[j1].Num++
	}

	aiCards[len(aiCards)-1].Num = 3
	for j1 := 0; j1 < len(aiCards)-1; j1++ {
		if aiCards[j1].Num == 0 {
			continue
		}
		aiCards[j1].Num--
		for j2 := j1; j2 < len(aiCards)-1; j2++ {
			if aiCards[j2].Num == 0 {
				continue
			}
			aiCards[j2].Num--
			for j3 := j2; j3 < len(aiCards)-1; j3++ {
				if aiCards[j3].Num == 0 {
					continue
				}
				aiCards[j3].Num--
				fmt.Println(j1, j2, j3, CheckHuForLZ(aiCards, COLOR_TIAO*MAHJONG_MASK+MAHJONG_1))
				aiCards[j3].Num++
			}
			aiCards[j2].Num++
		}
		aiCards[j1].Num++
	}

	aiCards[len(aiCards)-1].Num = 4
	for j1 := 0; j1 < len(aiCards)-1; j1++ {
		if aiCards[j1].Num == 0 {
			continue
		}
		aiCards[j1].Num--
		for j2 := j1; j2 < len(aiCards)-1; j2++ {
			if aiCards[j2].Num == 0 {
				continue
			}
			aiCards[j2].Num--
			for j3 := j2; j3 < len(aiCards)-1; j3++ {
				if aiCards[j3].Num == 0 {
					continue
				}
				aiCards[j3].Num--
				for j4 := j3; j4 < len(aiCards)-1; j4++ {
					if aiCards[j4].Num == 0 {
						continue
					}
					aiCards[j4].Num--
					fmt.Println(j1, j2, j3, j4, CheckHuForLZ(aiCards, COLOR_TIAO*MAHJONG_MASK+MAHJONG_1))
					aiCards[j4].Num++
				}
				aiCards[j3].Num++
			}
			aiCards[j2].Num++
		}
		aiCards[j1].Num++
	}
}

func TestCheckTing(t *testing.T) {
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
	fmt.Println("TestCheckTing: ", CheckTing(aiCards))
}

func TestCheckTingForLZ(t *testing.T) {
	aiCards := []AICard{
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_1, 2},
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_2, 1},
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_3, 3},
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_4, 1},
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_5, 1},
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_6, 1},		
		 {COLOR_WAN*MAHJONG_MASK + MAHJONG_9, 1},
		 {COLOR_TIAO*MAHJONG_MASK + MAHJONG_1, 3},
	}
	fmt.Println("TestCheckTingForLZ: ", CheckTingForLZ(aiCards, COLOR_TIAO*MAHJONG_MASK + MAHJONG_1))
}