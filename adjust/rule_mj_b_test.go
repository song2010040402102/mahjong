package logic

import (
	"testing"
)

func BenchmarkCheckHuType(b *testing.B) {
	if TestMap["rule_mj_b"] == 0 {
		return
	}
	r := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	r.LaiziCard = 415
	holdCards := []int32{101, 101, 101, 201, 203, 204, 205, 205, 206, 307, 308, 415, 415}
	for i := 0; i < b.N; i++ {
		r.CheckHuType([]*ChiCard{}, holdCards, 309, false)
	}
}

func BenchmarkCheckTing(b *testing.B) {
	if TestMap["rule_mj_b"] == 0 {
		return
	}
	r := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	r.LaiziCard = 415
	holdCards := []int32{101, 101, 101, 201, 203, 204, 205, 205, 206, 307, 308, 415, 415}
	for i := 0; i < b.N; i++ {
		r.CheckTing([]*ChiCard{}, holdCards, true)
	}
}

func BenchmarkCheckOneTing(b *testing.B) {
	if TestMap["rule_mj_b"] == 0 {
		return
	}
	r := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	r.LaiziCard = 415
	holdCards := []int32{101, 101, 201, 203, 204, 205, 205, 206, 208, 307, 308, 415, 415}
	for i := 0; i < b.N; i++ {
		r.CheckNTing([]*ChiCard{}, holdCards, 1, true)
	}
}

func BenchmarkCheckTwoTing(b *testing.B) {
	if TestMap["rule_mj_b"] == 0 {
		return
	}
	r := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	r.LaiziCard = 415
	holdCards := []int32{101, 101, 201, 203, 204, 205, 205, 206, 208, 305, 307, 308, 415}
	for i := 0; i < b.N; i++ {
		r.CheckNTing([]*ChiCard{}, holdCards, 2, true)
	}
}

func BenchmarkCheckThreeTing(b *testing.B) {
	if TestMap["rule_mj_b"] == 0 {
		return
	}
	r := NewRuleMahjong(RULE_ZJ_MAHJONG_TAIZHOU_HY)
	r.LaiziCard = 415
	holdCards := []int32{101, 101, 201, 203, 204, 206, 206, 207, 209, 304, 305, 307, 308}
	for i := 0; i < b.N; i++ {
		r.CheckNTing([]*ChiCard{}, holdCards, 3, true)
	}
}
