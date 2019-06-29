package logic

import (
	"fmt"
	"testing"
)

func TestCheckCommonHu(t *testing.T) {
	cards := []AICard{
		{105, 2}, {203, 3}, {305, 1}, {306, 2}, {307, 2}, {308, 1}, {417, 3},
	}
	ok, groups := CheckCommonHu(cards, true)
	fmt.Println("TestCheckCommonHu: ", ok, groups)
}

func TestCheckCommonHuForLZ(t *testing.T) {
	cards := []AICard{
		{101, 2}, {102, 2}, {103, 2}, {104, 1}, {105, 1}, {106, 3}, {108, 2},
	}
	ok, groups := CheckCommonHuForLZ(cards, 1, true)
	fmt.Println("TestCheckCommonHuForLZ: ", ok, groups)
}

func TestCheckCommonTing(t *testing.T) {
	cards := []AICard{
		{105, 2}, {203, 3}, {305, 1}, {306, 2}, {307, 1}, {308, 1}, {417, 3},
	}
	fmt.Println("TestCheckCommonTing: ", CheckCommonTing(cards, true))
}

func TestCheckCommonTingForLZ(t *testing.T) {
	cards := []AICard{
		{101, 3}, {201, 1}, {203, 1}, {204, 1}, {205, 2}, {206, 1}, {307, 1}, {308, 1},
	}
	fmt.Println("TestCheckCommonTingForLZ: ", CheckCommonTingForLZ(cards, 2, true))
}

func TestCheckQiDuiHu(t *testing.T) {
	cards := []AICard{
		{105, 2}, {203, 2}, {301, 2}, {306, 2}, {307, 2}, {414, 2}, {417, 2},
	}
	fmt.Println("TestCheckQiDuiHu: ", CheckQiDuiHu(cards))
}

func TestCheckQiDuiHuForLZ(t *testing.T) {
	cards := []AICard{
		{105, 2}, {203, 2}, {301, 2}, {306, 1}, {307, 2}, {414, 2}, {417, 2},
	}
	fmt.Println("TestCheckQiDuiHuForLZ: ", CheckQiDuiHuForLZ(cards, 1))
}

func TestCheckQiDuiTing(t *testing.T) {
	cards := []AICard{
		{105, 2}, {203, 2}, {301, 2}, {306, 2}, {307, 1}, {414, 2}, {417, 2},
	}
	fmt.Println("TestCheckQiDuiTing: ", CheckQiDuiTing(cards))
}

func TestCheckQiDuiTingForLZ(t *testing.T) {
	cards := []AICard{
		{105, 2}, {203, 1}, {301, 1}, {306, 2}, {307, 1}, {414, 1}, {417, 2},
	}
	fmt.Println("TestCheckQiDuiTingForLZ: ", CheckQiDuiTingForLZ(cards, 3))
}

func TestCheckQuanBuKao(t *testing.T) {
	cards := []AICard{
		{105, 1}, {109, 1}, {203, 1}, {301, 1}, {306, 1}, {414, 1}, {417, 1},
	}
	fmt.Println("TestCheckQuanBuKaoHu: ", CheckQuanBuKaoHu(cards))
	fmt.Println("TestCheckQuanBuKaoTing: ", CheckQuanBuKaoTing(cards))
}
