package logic

import (
	"fmt"
	"testing"
)

var TestMap map[string]int32 = map[string]int32{
	"rule_mj":    1,
	"rule_mj_b":  1,
	"rule_niu_b": 0,
	"rule_dan":   0,
	"room_mj":    0,
	"room_mj_b":  0,
}

func TestCheckHuType(t *testing.T) {
	if TestMap["rule_mj"] == 0 {
		return
	}
	r := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	r.lzCards = []int32{415}
	holdCards := []int32{101, 101, 101, 201, 203, 204, 205, 205, 206, 307, 308, 415, 415}
	huType, groups := r.CheckHuType(1, []*ChiCard{}, []int32{}, holdCards, 309, MJ_CHECK_HU_FLAG_GROUP)
	fmt.Println("TestCheckHuType, huType: ", huType, " groups: ", groups)
}

func TestCheckTing(t *testing.T) {
	if TestMap["rule_mj"] == 0 {
		return
	}
	r := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	r.lzCards = []int32{415}
	holdCards := []int32{101, 101, 101, 201, 203, 204, 205, 205, 206, 307, 308, 415, 415}
	fmt.Println("TestCheckTing, ", r.CheckTing([]*ChiCard{}, holdCards, true))
}

func TestCheckOneTing(t *testing.T) {
	if TestMap["rule_mj"] == 0 {
		return
	}
	r := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	r.lzCards = []int32{415}
	holdCards := []int32{101, 101, 201, 203, 204, 205, 205, 206, 208, 307, 308, 415, 415}
	realN, tings := r.CheckNTing([]*ChiCard{}, holdCards, 1, true)
	fmt.Println("TestCheckOneTing, ", realN, tings)
}

func TestCheckTwoTing(t *testing.T) {
	if TestMap["rule_mj"] == 0 {
		return
	}
	r := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	r.lzCards = []int32{415}
	holdCards := []int32{101, 101, 201, 203, 204, 205, 205, 206, 208, 305, 307, 308, 415}
	realN, tings := r.CheckNTing([]*ChiCard{}, holdCards, 2, false)
	fmt.Println("TestCheckTwoTing, ", realN, tings)
}

func TestCheckThreeTing(t *testing.T) {
	if TestMap["rule_mj"] == 0 {
		return
	}
	r := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	r.lzCards = []int32{415}
	holdCards := []int32{101, 101, 201, 203, 204, 206, 206, 207, 209, 304, 305, 307, 308}
	realN, tings := r.CheckNTing([]*ChiCard{}, holdCards, 3, true)
	fmt.Println("TestCheckThreeTing, ", realN, tings)
}
