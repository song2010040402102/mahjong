package logic

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

func TestCheck13BuKao(t *testing.T) {
	cards := []AICard{
		{105, 1}, {109, 1}, {203, 1}, {301, 1}, {306, 1}, {414, 1}, {417, 1},
	}
	ok, groups := Check13BuKaoHu(cards, true)
	fmt.Println("TestCheck13BuKaoHu: ", ok, groups)
	fmt.Println("TestCheck13BuKaoTing: ", Check13BuKaoTing(cards))
}

func TestCheckQuanBuKao(t *testing.T) {
	cards := []AICard{
		{105, 1}, {108, 1}, {203, 1}, {206, 1}, {209, 1}, {301, 1}, {304, 1},
		{307, 1}, {411, 1}, {412, 1}, {413, 1}, {415, 1}, {416, 1}, {417, 1},
	}
	ok, groups := CheckQuanBuKaoHu(cards, true)
	fmt.Println("TestCheckQuanBuKaoHu: ", ok, groups)
	fmt.Println("TestCheckQuanBuKaoTing: ", CheckQuanBuKaoTing(cards[1:]))
}

func TestCheck7StarBuKao(t *testing.T) {
	cards := []AICard{
		{105, 1}, {108, 1}, {203, 1}, {206, 1}, {209, 1}, {301, 1}, {304, 1},
		{411, 1}, {412, 1}, {413, 1}, {414, 1}, {415, 1}, {416, 1}, {417, 1},
	}
	ok, groups := Check7StarBuKaoHu(cards, true)
	fmt.Println("TestCheck7StarBuKaoHu: ", ok, groups)
	fmt.Println("TestCheck7StarBuKaoTing: ", Check7StarBuKaoTing(cards[1:]))
}

func TestCheck13Yao(t *testing.T) {
	cards := []AICard{
		{101, 1}, {109, 1}, {201, 1}, {209, 1}, {301, 1}, {309, 1},
		{411, 1}, {412, 1}, {413, 1}, {414, 1}, {415, 1}, {416, 1}, {417, 2},
	}
	ok, groups := Check13YaoHu(cards, true)
	fmt.Println("TestCheck13YaoHu: ", ok, groups)
	fmt.Println("TestCheck13YaoTing: ", Check13YaoTing(cards[1:]))
}

func TestCheckZuHeLong(t *testing.T) {
	cards := []AICard{
		{101, 1}, {104, 1}, {107, 1}, {202, 1}, {205, 1}, {208, 1}, {303, 1}, {306, 1}, {309, 1},
		{411, 1}, {412, 1}, {415, 1}, {416, 1}, {417, 1},
	}
	ok, groups := CheckZuHeLongHu([]*ChiCard{}, cards, true)
	fmt.Println("TestCheckZuHeLongHu: ", ok, groups)
	fmt.Println("TestCheckZuHeLongTing: ", CheckZuHeLongTing([]*ChiCard{}, cards[1:]))
}

func TestCheckQingLong(t *testing.T) {
	cards := [][3]int32{
		{101, 102, 103}, {104, 105, 106}, {107, 108, 109},
	}
	fmt.Println("TestCheckQingLong: ", CheckQingLong([]*ChiCard{}, cards))
}

func TestCheckHuaLong(t *testing.T) {
	cards := [][3]int32{
		{101, 102, 103}, {104, 105, 106}, {307, 308, 309},
	}
	fmt.Println("TestCheckHuaLong: ", CheckHuaLong([]*ChiCard{}, cards))
}

func TestCheckLuanSanFeng(r *testing.T) {
	cards := []AICard{
		//{411, 1}, {412, 2}, {413, 2}, {414, 1}, {101, 2},
		//{411, 3}, {412, 3}, {413, 3}, {414, 3}, {101, 2},
		//{411, 1}, {412, 1}, {414, 1}, {415, 1}, {416, 1}, {417, 1}, {101, 2},
		//{411, 2}, {412, 2}, {413, 2}, {101, 2},
		{101, 3}, {411, 2}, {412, 3}, {413, 2}, {414, 1},
	}
	ok, groups := CheckLuanSanFengHu(cards, true)
	fmt.Println("TestCheckLuanSanFeng: ", ok, groups)
}
