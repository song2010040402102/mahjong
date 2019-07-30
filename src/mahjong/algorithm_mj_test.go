package mahjong

import (
	"fmt"
	"testing"
)

func TestCheckCommonHu(t *testing.T) {
	cards := []AICard{
		{104, 2}, {105, 1}, {106, 2}, {107, 2}, {108, 1}, {204, 2}, {205, 1}, {206, 2}, {207, 2}, {208, 1},
	}
	ok, groups := CheckCommonHu(cards, 1, true)
	fmt.Println("TestCheckCommonHu: ", ok, groups)
}

func TestCheckCommonTing(t *testing.T) {
	cards := []AICard{
		{101, 3}, {201, 1}, {203, 1}, {204, 1}, {205, 2}, {206, 1}, {307, 1}, {308, 1},
	}
	fmt.Println("TestCheckCommonTing: ", CheckCommonTing(cards, 2, true))
}

func TestCheckQiDuiHu(t *testing.T) {
	cards := []AICard{
		{105, 2}, {203, 2}, {301, 2}, {306, 1}, {307, 2}, {414, 2}, {417, 2},
	}
	ok, groups := CheckQiDuiHu(cards, 1, true)
	fmt.Println("TestCheckQiDuiHu: ", ok, groups)
}

func TestCheckQiDuiTing(t *testing.T) {
	cards := []AICard{
		{105, 2}, {203, 1}, {301, 1}, {306, 2}, {307, 1}, {414, 1}, {417, 2},
	}
	fmt.Println("TestCheckQiDuiTing: ", CheckQiDuiTing(cards, 3))
}
