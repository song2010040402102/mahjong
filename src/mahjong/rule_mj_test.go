package mahjong

import (
	"fmt"
	"testing"
)

func TestCheckHuType(t *testing.T) {
	r := NewRuleMahjong(RULE_HN_MAHJONG_HONGZHONG)
	r.SetLaiziCard([]int32{415})
	holdCards := []int32{101, 101, 101, 201, 203, 204, 205, 205, 206, 307, 308, 415, 415}
	huType, groups := r.CheckHuType([]*ChiCard{}, holdCards, 309, MJ_CHECK_HU_FLAG_GROUP)
	fmt.Println("TestCheckHuType, huType: ", huType, " groups: ", groups)
}

func TestCheckTing(t *testing.T) {
	r := NewRuleMahjong(RULE_HN_MAHJONG_HONGZHONG)
	r.SetLaiziCard([]int32{415})
	holdCards := []int32{101, 101, 101, 201, 203, 204, 205, 205, 206, 307, 308, 415, 415}
	fmt.Println("TestCheckTing, ", r.CheckTing([]*ChiCard{}, holdCards, true))
}

func TestCheckOneTing(t *testing.T) {
	r := NewRuleMahjong(RULE_HN_MAHJONG_HONGZHONG)
	r.SetLaiziCard([]int32{415})
	holdCards := []int32{101, 101, 201, 203, 204, 205, 205, 206, 208, 307, 308, 415, 415}
	realN, tings := r.CheckNTing([]*ChiCard{}, holdCards, 1, true)
	fmt.Println("TestCheckOneTing, ", realN, tings)
}

func TestCheckTwoTing(t *testing.T) {
	r := NewRuleMahjong(RULE_HN_MAHJONG_HONGZHONG)
	r.SetLaiziCard([]int32{415})
	holdCards := []int32{101, 101, 201, 203, 204, 205, 205, 206, 208, 305, 307, 308, 415}
	realN, tings := r.CheckNTing([]*ChiCard{}, holdCards, 2, true)
	fmt.Println("TestCheckTwoTing, ", realN, tings)
}

func TestCheckThreeTing(t *testing.T) {
	r := NewRuleMahjong(RULE_HN_MAHJONG_HONGZHONG)
	r.SetLaiziCard([]int32{415})
	holdCards := []int32{101, 101, 201, 203, 204, 206, 206, 207, 209, 304, 305, 307, 308}
	realN, tings := r.CheckNTing([]*ChiCard{}, holdCards, 3, true)
	fmt.Println("TestCheckThreeTing, ", realN, tings)
}
