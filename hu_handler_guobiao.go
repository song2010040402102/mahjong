package logic

import (
	"github.com/astaxie/beego/logs"
	"util"
)

type GUOBIAO_MJ_HuHandler struct {
}

func New_GUOBIAO_MJ_HuHandler() *GUOBIAO_MJ_HuHandler {
	c := &GUOBIAO_MJ_HuHandler{}
	return c
}

func (c *GUOBIAO_MJ_HuHandler) GetHuInfo(r *RoomMahjong, index int32, huType uint64, card int32, fromIndex int32) *HuInfo {
	if huType == MJ_HU_TYPE_NONE {
		logs.Error("[GUOBIAO] GetHuInfo huType none!")
		return nil
	}

	huInfo := MakeHuInfo(r, index, huType, card, fromIndex)
	huInfo.MaxHuNum, huInfo.HuNumInfo, huInfo.CardGroups = c.getMaxHuGroup(r, huInfo)
	return huInfo
}

func (c *GUOBIAO_MJ_HuHandler) FilterHu(r *RoomMahjong, player *MahjongPlayer, huInfo *HuInfo) bool {
	if huInfo == nil {
		logs.Error("[GUOBIAO]FilterHu, huInfo nil!")
		return false
	}
	return true
}

func (c *GUOBIAO_MJ_HuHandler) CalcBalance(r *RoomMahjong) {
	firstHu := r.GetFirstHuIndex()
	if firstHu <= 0 || firstHu > r.GetRule().GetPlayerLimit() {
		return
	}
	if r.players[firstHu].HuInfos[0].FromIndex == firstHu {
		r.players[firstHu].Integal += (r.GetRule().GetPlayerLimit() - 1) * r.players[firstHu].HuInfos[0].MaxHuNum
		for k, v := range r.players {
			if k != firstHu {
				v.Integal -= r.players[firstHu].HuInfos[0].MaxHuNum
			}
		}
	} else {
		r.players[firstHu].Integal += r.players[firstHu].HuInfos[0].MaxHuNum
		r.players[r.players[firstHu].HuInfos[0].FromIndex].Integal -= r.players[firstHu].HuInfos[0].MaxHuNum
	}
}

func (c *GUOBIAO_MJ_HuHandler) getMaxHuGroup(r *RoomMahjong, huInfo *HuInfo) (int32, map[int32]int32, [][3]int32) {
	maxHu := int32(0)
	var maxGroup [][3]int32
	var maxMapHuType map[uint64]bool
	var maxMapSubType map[uint32]bool
	var maxMapExtraType map[uint64]bool
	var maxMapHuNum map[int32]int32
	mapHuType, mapSubType, mapExtraType, mapHuNum := make(map[uint64]bool), make(map[uint32]bool), make(map[uint64]bool), make(map[int32]int32)
	huCards, _ := ConvertSliceToAICard(r.players[huInfo.Index].GetCardArray(), huInfo.HuCard, r.GetRule().GetLaiziCard())
	tingCards, _ := ConvertSliceToAICard(r.players[huInfo.Index].GetCardArray(), 0, r.GetRule().GetLaiziCard())
	if huInfo.HuType&MJ_HU_TYPE_COMMON != 0 {
		mapHuType[MJ_HU_TYPE_COMMON] = true
	}
	if huInfo.HuType&MJ_HU_TYPE_QIDUI != 0 {
		mapHuType[MJ_HU_TYPE_QIDUI] = true
	}
	if huInfo.HuType&MJ_HU_TYPE_QUANBUKAO != 0 {
		mapHuType[MJ_HU_TYPE_QUANBUKAO] = true
	}
	if huInfo.HuType&MJ_HU_TYPE_7STARBUKAO != 0 {
		mapHuType[MJ_HU_TYPE_7STARBUKAO] = true
	}
	if huInfo.HuType&MJ_HU_TYPE_13YAO != 0 {
		mapHuType[MJ_HU_TYPE_13YAO] = true
	}
	if huInfo.HuType&MJ_HU_TYPE_ZUHELONG != 0 {
		mapHuType[MJ_HU_TYPE_ZUHELONG] = true
	}
	if huInfo.ExtraType&MJ_HU_EXTRA_TYPE_ZIMO != 0 {
		mapExtraType[MJ_HU_EXTRA_TYPE_ZIMO] = true
	}
	if huInfo.ExtraType&MJ_HU_EXTRA_TYPE_QIANGGANG != 0 {
		mapExtraType[MJ_HU_EXTRA_TYPE_QIANGGANG] = true
	}
	if huInfo.ExtraType&MJ_HU_EXTRA_TYPE_GANGSHANGHUA != 0 {
		mapExtraType[MJ_HU_EXTRA_TYPE_GANGSHANGHUA] = true
	}
	if huInfo.ExtraType&MJ_HU_EXTRA_TYPE_HAIDILAOYUE != 0 {
		mapExtraType[MJ_HU_EXTRA_TYPE_HAIDILAOYUE] = true
	}
	if huInfo.ExtraType&MJ_HU_EXTRA_TYPE_HAIDIPAO != 0 {
		mapExtraType[MJ_HU_EXTRA_TYPE_HAIDIPAO] = true
	}
	if CheckQingYiSe(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_QINGYISE] = true
	} else if CheckHunYiSe(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_HUNYISE] = true
	} else if CheckZiYiSe(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_ZIYISE] = true
	}
	if CheckLvYiSe(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_LVYISE] = true
	}
	if Check9LianBaoDeng(huCards) {
		mapSubType[MJ_HU_SUB_TYPE_JIULIANBAODENG] = true
	}
	if t1, t2 := c.getGangType(r.players[huInfo.Index].ChiCards); t1 > 0 {
		mapSubType[t1] = true
		if t2 > 0 {
			mapSubType[t2] = true
		}
	}
	if CheckYiSeSameLong(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_YISETWOLONG] = true
	} else if Check3SeSameLong(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_SANSETWOLONG] = true
	}
	if CheckQuanPairKe(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_QUANPAIRKE] = true
	}
	if CheckZhongZhang(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_ZHONGZHANG] = true
	}
	if CheckQuanDa(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_QUANDA] = true
	} else if CheckQuanZhong(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_QUANZHONG] = true
	} else if CheckQuanXiao(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_QUANXIAO] = true
	}
	if CheckDaYu5(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_DAYUWU] = true
	} else if CheckXiaoYu5(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_XIAOYUWU] = true
	}
	if Check3FengKeGang(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_SANFENGKE] = true
	}
	if CheckFengKeGang(r.players[huInfo.Index].ChiCards, huCards, COLOR_OTHER*MAHJONG_MASK+r.GetRule().GetPlayerDoor(huInfo.Index)) {
		mapSubType[MJ_HU_SUB_TYPE_MENFENGKE] = true
	}
	if CheckFengKeGang(r.players[huInfo.Index].ChiCards, huCards, COLOR_OTHER*MAHJONG_MASK+r.GetRule().GetCircleWind()+MAHJONG_DONG-1) {
		mapSubType[MJ_HU_SUB_TYPE_QUANFENGKE] = true
	}
	if n := GetJianKeGangNum(r.players[huInfo.Index].ChiCards, huCards); n == 2 {
		mapSubType[MJ_HU_SUB_TYPE_SHUANGJIANKE] = true
	} else if n == 1 {
		mapSubType[MJ_HU_SUB_TYPE_JIANKE] = true
	}
	if CheckTuiBuDao(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_TUIBUDAO] = true
	}
	if Check5MenQi(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_WUMENQI] = true
	} else if CheckQueYiMen(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_QUEYIMEN] = true
	}
	if CheckQuanQiuRen(r.players[huInfo.Index].ChiCards, huCards, huInfo.Index == huInfo.FromIndex) {
		mapSubType[MJ_HU_SUB_TYPE_QUANQIUREN] = true
	} else if CheckBuQiuRen(r.players[huInfo.Index].ChiCards, huInfo.Index == huInfo.FromIndex) {
		mapSubType[MJ_HU_SUB_TYPE_BUQIUREN] = true
	}
	if r.CheckJueZhang(huInfo.HuCard) {
		mapSubType[MJ_HU_SUB_TYPE_HEJUEZHANG] = true
	}
	if CheckMenQing(r.players[huInfo.Index].ChiCards) {
		mapSubType[MJ_HU_SUB_TYPE_MENQING] = true
	}
	if Check4Gui1(huCards) {
		mapSubType[MJ_HU_SUB_TYPE_SIGUIYI] = true
	}
	if CheckWuZi(r.players[huInfo.Index].ChiCards, huCards) {
		mapSubType[MJ_HU_SUB_TYPE_WUZI] = true
	}
	if CheckTingTypeForBianZhang(tingCards) {
		mapSubType[MJ_HU_SUB_TYPE_BIANZHANG] = true
	} else if CheckTingTypeForQiaZi(tingCards) {
		mapSubType[MJ_HU_SUB_TYPE_QIAZI] = true
	} else if CheckTingTypeForDanDiao(tingCards) {
		mapSubType[MJ_HU_SUB_TYPE_DANDIAO] = true
	}
	_, cardGroups := r.GetRule().CheckHuType(huInfo.Index, r.players[huInfo.Index].ChiCards, r.players[huInfo.Index].HuaCards, r.players[huInfo.Index].GetCardArray(), huInfo.HuCard, MJ_CHECK_HU_FLAG_GROUP)
	for _, group := range cardGroups {
		mapTmpHuType := util.CopyMap(mapHuType).(map[uint64]bool)
		mapTmpSubType := util.CopyMap(mapSubType).(map[uint32]bool)
		mapTmpExtraType := util.CopyMap(mapExtraType).(map[uint64]bool)
		mapTmpHuNum := util.CopyMap(mapHuNum).(map[int32]int32)

		if CheckDaSiXi(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_DASIXI] = true
		} else if CheckDaSanYuan(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_DASANYUAN] = true
		} else if CheckXiaoSiXi(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_XIAOSIXI] = true
		} else if CheckXiaoSanYuan(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_XIAOSANYUAN] = true
		}
		if CheckLian7Dui(group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_LIANQIDUI] = true
		}
		if CheckQingYaoJiu(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_QINGYAOJIU] = true
		} else if CheckHunYaoJiu(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_HUNYAOJIU] = true
		}
		if n := GetAnKeNum(r.players[huInfo.Index].ChiCards, group, huInfo.HuCard, huInfo.Index == huInfo.FromIndex); n == 4 {
			mapTmpSubType[MJ_HU_SUB_TYPE_SIANKE] = true
		} else if n == 3 {
			mapTmpSubType[MJ_HU_SUB_TYPE_SANANKE] = true
		} else if n == 2 {
			mapTmpSubType[MJ_HU_SUB_TYPE_TWOANKE] = true
		}
		if CheckYiSe4SameSeq(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_YISESITONG] = true
		} else if CheckYiSe3SameSeq(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_YISESANTONG] = true
		} else if Check3Se3SameSeq(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_SANSESANTONG] = true
		} else if CheckYiSe4JieGao(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_YISESIJIE] = true
		} else if CheckYiSe3JieGao(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_YISESANJIE] = true
		} else if Check3Se3JieGao(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_SANSESANJIE] = true
		} else if CheckYiSe4BuGao(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_YISESIBU] = true
		} else if CheckYiSe3BuGao(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_YISESANBU] = true
		} else if Check3Se3BuGao(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_SANSESANBU] = true
		}
		if CheckQingLong(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_QINGLONG] = true
		} else if CheckHuaLong(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_HUALONG] = true
		}
		if CheckQuanDai5(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_QUANDAIWU] = true
		}
		if CheckDuiDuiHuForLZ(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_DUIDUIHU] = true
		}
		if CheckDai19(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_DAI19] = true
		}
		if CheckPingHu(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_PINGHU] = true
		}
		if Check3SameKeGang(r.players[huInfo.Index].ChiCards, group) {
			mapTmpSubType[MJ_HU_SUB_TYPE_SANTONGKE] = true
		}
		if n := GetTwoSameKeNum(r.players[huInfo.Index].ChiCards, group); n > 0 {
			mapTmpHuNum[MJ_HU_NUM_TYPE_SAME_KE] = n
		}
		if n := GetYiBanGaoNum(r.players[huInfo.Index].ChiCards, group, false); n > 0 {
			mapTmpHuNum[MJ_HU_NUM_TYPE_YIBANGAO] = n
		}
		if n := GetXiXiangFengNum(r.players[huInfo.Index].ChiCards, group); n > 0 {
			mapTmpHuNum[MJ_HU_NUM_TYPE_XIXIANGFENG] = n
		}
		if n := GetLian6Num(r.players[huInfo.Index].ChiCards, group, false); n > 0 {
			mapTmpHuNum[MJ_HU_NUM_TYPE_LIANLIU] = n
		}
		if n := GetLaoShaoFuNum(r.players[huInfo.Index].ChiCards, group); n > 0 {
			mapTmpHuNum[MJ_HU_NUM_TYPE_LAOSHAOFU] = n
		}
		if n := GetZi19KeGangNum(r.players[huInfo.Index].ChiCards, group); n > 0 {
			mapTmpHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19] = n
		}
		if len(mapTmpHuType) == 1 && mapTmpHuType[MJ_HU_TYPE_COMMON] && len(mapTmpSubType) == 0 && len(mapTmpExtraType) == 0 && len(mapTmpHuNum) == 0 {
			mapTmpSubType[MJ_HU_SUB_TYPE_WUFAN] = true
		}
		c.removeContainType(r, huInfo, mapTmpHuType, mapTmpSubType, mapTmpExtraType, mapTmpHuNum)
		if hu := c.calcFan(mapTmpHuType, mapTmpSubType, mapTmpExtraType, mapTmpHuNum, int32(len(r.players[huInfo.Index].HuaCards))); hu > maxHu {
			maxHu, maxGroup = hu, group
			maxMapHuType, maxMapSubType, maxMapExtraType, maxMapHuNum = mapTmpHuType, mapTmpSubType, mapTmpExtraType, mapTmpHuNum
		}
	}
	c.convertHuInfo(huInfo, maxMapHuType, maxMapSubType, maxMapExtraType)
	return maxHu, maxMapHuNum, maxGroup
}

func (c *GUOBIAO_MJ_HuHandler) removeContainType(r *RoomMahjong, huInfo *HuInfo, mapHuType map[uint64]bool, mapSubType map[uint32]bool, mapExtraType map[uint64]bool, mapHuNum map[int32]int32) {
	if mapSubType[MJ_HU_SUB_TYPE_DASIXI] {
		delete(mapSubType, MJ_HU_SUB_TYPE_QUANFENGKE)
		delete(mapSubType, MJ_HU_SUB_TYPE_MENFENGKE)
		delete(mapSubType, MJ_HU_SUB_TYPE_SANFENGKE)
		delete(mapSubType, MJ_HU_SUB_TYPE_DUIDUIHU)
		delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
	} else if mapSubType[MJ_HU_SUB_TYPE_DASANYUAN] {
		delete(mapSubType, MJ_HU_SUB_TYPE_SHUANGJIANKE)
	} else if mapSubType[MJ_HU_SUB_TYPE_XIAOSIXI] {
		delete(mapSubType, MJ_HU_SUB_TYPE_SANFENGKE)
		if mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19] > 3 {
			mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19] -= 3
		} else {
			delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
		}
	} else if mapSubType[MJ_HU_SUB_TYPE_XIAOSANYUAN] {
		delete(mapSubType, MJ_HU_SUB_TYPE_SHUANGJIANKE)
	}
	if mapSubType[MJ_HU_SUB_TYPE_LVYISE] && mapHuType[MJ_HU_TYPE_QIDUI] {
		delete(mapSubType, MJ_HU_SUB_TYPE_SIGUIYI)
	}
	if mapSubType[MJ_HU_SUB_TYPE_JIULIANBAODENG] {
		delete(mapSubType, MJ_HU_SUB_TYPE_QINGYISE)
		delete(mapSubType, MJ_HU_SUB_TYPE_MENQING)
		delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
	}
	if mapSubType[MJ_HU_SUB_TYPE_SIGANG] {
		delete(mapSubType, MJ_HU_SUB_TYPE_DUIDUIHU)
		delete(mapSubType, MJ_HU_SUB_TYPE_DANDIAO)
	}
	if mapSubType[MJ_HU_SUB_TYPE_LIANQIDUI] {
		delete(mapSubType, MJ_HU_SUB_TYPE_QINGYISE)
		delete(mapSubType, MJ_HU_SUB_TYPE_MENQING)
		delete(mapSubType, MJ_HU_SUB_TYPE_DANDIAO)
		delete(mapHuType, MJ_HU_TYPE_QIDUI)
	}
	if mapHuType[MJ_HU_TYPE_13YAO] {
		delete(mapSubType, MJ_HU_SUB_TYPE_WUMENQI)
		delete(mapSubType, MJ_HU_SUB_TYPE_MENQING)
		delete(mapSubType, MJ_HU_SUB_TYPE_DANDIAO)
		delete(mapSubType, MJ_HU_SUB_TYPE_HUNYAOJIU)
	}
	if mapSubType[MJ_HU_SUB_TYPE_QINGYAOJIU] {
		delete(mapSubType, MJ_HU_SUB_TYPE_DUIDUIHU)
		delete(mapSubType, MJ_HU_SUB_TYPE_DAI19)
		delete(mapSubType, MJ_HU_SUB_TYPE_WUZI)
		delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
		if mapHuType[MJ_HU_TYPE_QIDUI] {
			delete(mapSubType, MJ_HU_SUB_TYPE_SIGUIYI)
		}
	} else if mapSubType[MJ_HU_SUB_TYPE_HUNYAOJIU] {
		delete(mapSubType, MJ_HU_SUB_TYPE_DUIDUIHU)
		delete(mapSubType, MJ_HU_SUB_TYPE_DAI19)
		delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
	}
	if mapSubType[MJ_HU_SUB_TYPE_ZIYISE] {
		delete(mapSubType, MJ_HU_SUB_TYPE_DUIDUIHU)
		delete(mapSubType, MJ_HU_SUB_TYPE_DAI19)
		delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
	}
	if mapSubType[MJ_HU_SUB_TYPE_YISETWOLONG] {
		delete(mapSubType, MJ_HU_SUB_TYPE_QINGYISE)
		delete(mapSubType, MJ_HU_SUB_TYPE_PINGHU)
		delete(mapSubType, MJ_HU_SUB_TYPE_WUZI)
		delete(mapHuNum, MJ_HU_NUM_TYPE_YIBANGAO)
		delete(mapHuNum, MJ_HU_NUM_TYPE_LAOSHAOFU)
	} else if mapSubType[MJ_HU_SUB_TYPE_SANSETWOLONG] {
		delete(mapSubType, MJ_HU_SUB_TYPE_PINGHU)
		delete(mapSubType, MJ_HU_SUB_TYPE_WUZI)
		delete(mapHuNum, MJ_HU_NUM_TYPE_XIXIANGFENG)
		delete(mapHuNum, MJ_HU_NUM_TYPE_LAOSHAOFU)
	}
	if mapSubType[MJ_HU_SUB_TYPE_YISESITONG] {
		delete(mapSubType, MJ_HU_SUB_TYPE_SIGUIYI)
		delete(mapHuNum, MJ_HU_NUM_TYPE_YIBANGAO)
	} else if mapSubType[MJ_HU_SUB_TYPE_YISESANTONG] {
		delete(mapHuNum, MJ_HU_NUM_TYPE_YIBANGAO)
	} else if mapSubType[MJ_HU_SUB_TYPE_YISESIJIE] {
		delete(mapSubType, MJ_HU_SUB_TYPE_DUIDUIHU)
	} else if mapSubType[MJ_HU_SUB_TYPE_YISESIBU] {
		delete(mapHuNum, MJ_HU_NUM_TYPE_LIANLIU)
		delete(mapHuNum, MJ_HU_NUM_TYPE_LAOSHAOFU)
	}
	if mapHuType[MJ_HU_TYPE_QIDUI] {
		delete(mapSubType, MJ_HU_SUB_TYPE_MENQING)
		delete(mapSubType, MJ_HU_SUB_TYPE_DANDIAO)
	}
	if mapHuType[MJ_HU_TYPE_7STARBUKAO] || mapHuType[MJ_HU_TYPE_QUANBUKAO] {
		delete(mapSubType, MJ_HU_SUB_TYPE_WUMENQI)
		delete(mapSubType, MJ_HU_SUB_TYPE_MENQING)
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANPAIRKE] {
		delete(mapSubType, MJ_HU_SUB_TYPE_DUIDUIHU)
		delete(mapSubType, MJ_HU_SUB_TYPE_ZHONGZHANG)
	}
	if mapSubType[MJ_HU_SUB_TYPE_QINGYISE] {
		delete(mapSubType, MJ_HU_SUB_TYPE_WUZI)
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANDA] || mapSubType[MJ_HU_SUB_TYPE_QUANZHONG] || mapSubType[MJ_HU_SUB_TYPE_QUANXIAO] || mapSubType[MJ_HU_SUB_TYPE_DAYUWU] || mapSubType[MJ_HU_SUB_TYPE_XIAOYUWU] {
		delete(mapSubType, MJ_HU_SUB_TYPE_WUZI)
		if mapSubType[MJ_HU_SUB_TYPE_QUANZHONG] {
			delete(mapSubType, MJ_HU_SUB_TYPE_ZHONGZHANG)
		}
	}
	if mapSubType[MJ_HU_SUB_TYPE_QINGLONG] {
		if mapHuNum[MJ_HU_NUM_TYPE_LIANLIU] > 2 {
			mapHuNum[MJ_HU_NUM_TYPE_LIANLIU] -= 2
		} else {
			delete(mapHuNum, MJ_HU_NUM_TYPE_LIANLIU)
		}
		if mapHuNum[MJ_HU_NUM_TYPE_LAOSHAOFU] > 1 {
			mapHuNum[MJ_HU_NUM_TYPE_LAOSHAOFU]--
		} else {
			delete(mapHuNum, MJ_HU_NUM_TYPE_LAOSHAOFU)
		}
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANDAIWU] {
		delete(mapSubType, MJ_HU_SUB_TYPE_ZHONGZHANG)
	}
	if mapSubType[MJ_HU_SUB_TYPE_SANFENGKE] {
		mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19] -= 3
		if mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19] <= 0 {
			delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
		}
	}
	if mapSubType[MJ_HU_SUB_TYPE_TUIBUDAO] {
		delete(mapSubType, MJ_HU_SUB_TYPE_QUEYIMEN)
	}
	if mapExtraType[MJ_HU_EXTRA_TYPE_HAIDILAOYUE] || mapExtraType[MJ_HU_EXTRA_TYPE_GANGSHANGHUA] {
		delete(mapExtraType, MJ_HU_EXTRA_TYPE_ZIMO)
	}
	if mapExtraType[MJ_HU_EXTRA_TYPE_QIANGGANG] {
		delete(mapSubType, MJ_HU_SUB_TYPE_HEJUEZHANG)
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANQIUREN] {
		delete(mapSubType, MJ_HU_SUB_TYPE_DANDIAO)
	}
	if mapSubType[MJ_HU_SUB_TYPE_BUQIUREN] {
		delete(mapSubType, MJ_HU_SUB_TYPE_MENQING)
		delete(mapExtraType, MJ_HU_EXTRA_TYPE_ZIMO)
	}
	if mapSubType[MJ_HU_SUB_TYPE_SHUANGJIANKE] {
		mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19] -= 2
		if mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19] <= 0 {
			delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
		}
	} else if mapSubType[MJ_HU_SUB_TYPE_JIANKE] {
		mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19]--
		if mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19] <= 0 {
			delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
		}
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANFENGKE] {
		mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19]--
		if mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19] <= 0 {
			delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
		}
	}
	if mapSubType[MJ_HU_SUB_TYPE_MENFENGKE] {
		mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19]--
		if mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19] <= 0 {
			delete(mapHuNum, MJ_HU_NUM_TYPE_KE_GANG_ZI19)
		}
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANFENGKE] && mapSubType[MJ_HU_SUB_TYPE_MENFENGKE] && r.GetRule().GetPlayerDoor(huInfo.Index) == r.GetRule().GetCircleWind()+MAHJONG_DONG-1 {
		mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19]++
	}
	if mapSubType[MJ_HU_SUB_TYPE_PINGHU] {
		delete(mapSubType, MJ_HU_SUB_TYPE_WUZI)
	}

	//处理套算一次原则，这个原则着重解决多番型组合必然存在的番型移除问题，所以处理以下特例即可
	if mapSubType[MJ_HU_SUB_TYPE_HUALONG] && mapHuNum[MJ_HU_NUM_TYPE_XIXIANGFENG] == 1 && mapHuNum[MJ_HU_NUM_TYPE_LIANLIU] == 1 {
		//delete(mapHuNum, MJ_HU_NUM_TYPE_LIANLIU) //花龙+喜相逢+连六，按照套算一次原则这里应该去掉喜相逢或连六，但并没有必然关系，所以先去掉
	}
	if mapHuNum[MJ_HU_NUM_TYPE_XIXIANGFENG] == 2 && mapHuNum[MJ_HU_NUM_TYPE_LIANLIU] == 2 {
		mapHuNum[MJ_HU_NUM_TYPE_LIANLIU]-- //2喜相逢+2连六
	}
	if mapHuNum[MJ_HU_NUM_TYPE_XIXIANGFENG] == 2 && mapHuNum[MJ_HU_NUM_TYPE_LAOSHAOFU] == 2 {
		mapHuNum[MJ_HU_NUM_TYPE_LAOSHAOFU]-- //2喜相逢+2老少副
	}
	if mapHuNum[MJ_HU_NUM_TYPE_YIBANGAO] == 2 && mapHuNum[MJ_HU_NUM_TYPE_LIANLIU] == 2 {
		mapHuNum[MJ_HU_NUM_TYPE_LIANLIU]-- //2一般高+2连六
	}
	if mapHuNum[MJ_HU_NUM_TYPE_YIBANGAO] == 2 && mapHuNum[MJ_HU_NUM_TYPE_LAOSHAOFU] == 2 {
		mapHuNum[MJ_HU_NUM_TYPE_LAOSHAOFU]-- //2一般高+2老少副
	}
	if mapHuNum[MJ_HU_NUM_TYPE_YIBANGAO] == 2 && mapHuNum[MJ_HU_NUM_TYPE_XIXIANGFENG] == 2 {
		mapHuNum[MJ_HU_NUM_TYPE_XIXIANGFENG]-- //2一般高+2喜相逢
	}
}

func (c *GUOBIAO_MJ_HuHandler) calcFan(mapHuType map[uint64]bool, mapSubType map[uint32]bool, mapExtraType map[uint64]bool, mapHuNum map[int32]int32, huaNum int32) int32 {
	hu := int32(0)
	if mapSubType[MJ_HU_SUB_TYPE_DASIXI] {
		hu += 88
	}
	if mapSubType[MJ_HU_SUB_TYPE_DASANYUAN] {
		hu += 88
	}
	if mapSubType[MJ_HU_SUB_TYPE_LVYISE] {
		hu += 88
	}
	if mapSubType[MJ_HU_SUB_TYPE_JIULIANBAODENG] {
		hu += 88
	}
	if mapSubType[MJ_HU_SUB_TYPE_SIGANG] {
		hu += 88
	}
	if mapSubType[MJ_HU_SUB_TYPE_LIANQIDUI] {
		hu += 88
	}
	if mapHuType[MJ_HU_TYPE_13YAO] {
		hu += 88
	}

	if mapSubType[MJ_HU_SUB_TYPE_QINGYAOJIU] {
		hu += 64
	}
	if mapSubType[MJ_HU_SUB_TYPE_XIAOSIXI] {
		hu += 64
	}
	if mapSubType[MJ_HU_SUB_TYPE_XIAOSANYUAN] {
		hu += 64
	}
	if mapSubType[MJ_HU_SUB_TYPE_ZIYISE] {
		hu += 64
	}
	if mapSubType[MJ_HU_SUB_TYPE_SIANKE] {
		hu += 64
	}
	if mapSubType[MJ_HU_SUB_TYPE_YISETWOLONG] {
		hu += 64
	}

	if mapSubType[MJ_HU_SUB_TYPE_YISESITONG] {
		hu += 48
	}
	if mapSubType[MJ_HU_SUB_TYPE_YISESIJIE] {
		hu += 48
	}

	if mapSubType[MJ_HU_SUB_TYPE_YISESIBU] {
		hu += 32
	}
	if mapSubType[MJ_HU_SUB_TYPE_SANGANG] {
		hu += 32
	}
	if mapSubType[MJ_HU_SUB_TYPE_HUNYAOJIU] {
		hu += 32
	}

	if mapHuType[MJ_HU_TYPE_QIDUI] {
		hu += 24
	}
	if mapHuType[MJ_HU_TYPE_7STARBUKAO] {
		hu += 24
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANPAIRKE] {
		hu += 24
	}
	if mapSubType[MJ_HU_SUB_TYPE_QINGYISE] {
		hu += 24
	}
	if mapSubType[MJ_HU_SUB_TYPE_YISESANTONG] {
		hu += 24
	}
	if mapSubType[MJ_HU_SUB_TYPE_YISESANJIE] {
		hu += 24
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANDA] {
		hu += 24
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANZHONG] {
		hu += 24
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANXIAO] {
		hu += 24
	}

	if mapSubType[MJ_HU_SUB_TYPE_QINGLONG] {
		hu += 16
	}
	if mapSubType[MJ_HU_SUB_TYPE_SANSETWOLONG] {
		hu += 16
	}
	if mapSubType[MJ_HU_SUB_TYPE_YISESANBU] {
		hu += 16
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANDAIWU] {
		hu += 16
	}
	if mapSubType[MJ_HU_SUB_TYPE_SANTONGKE] {
		hu += 16
	}
	if mapSubType[MJ_HU_SUB_TYPE_SANANKE] {
		hu += 16
	}

	if mapHuType[MJ_HU_TYPE_QUANBUKAO] {
		hu += 12
	}
	if mapHuType[MJ_HU_TYPE_ZUHELONG] {
		hu += 12
	}
	if mapSubType[MJ_HU_SUB_TYPE_DAYUWU] {
		hu += 12
	}
	if mapSubType[MJ_HU_SUB_TYPE_XIAOYUWU] {
		hu += 12
	}
	if mapSubType[MJ_HU_SUB_TYPE_SANFENGKE] {
		hu += 12
	}

	if mapSubType[MJ_HU_SUB_TYPE_HUALONG] {
		hu += 8
	}
	if mapSubType[MJ_HU_SUB_TYPE_TUIBUDAO] {
		hu += 8
	}
	if mapSubType[MJ_HU_SUB_TYPE_SANSESANTONG] {
		hu += 8
	}
	if mapSubType[MJ_HU_SUB_TYPE_SANSESANJIE] {
		hu += 8
	}
	if mapSubType[MJ_HU_SUB_TYPE_WUFAN] {
		hu += 8
	}
	if mapExtraType[MJ_HU_EXTRA_TYPE_HAIDILAOYUE] {
		hu += 8
	}
	if mapExtraType[MJ_HU_EXTRA_TYPE_HAIDIPAO] {
		hu += 8
	}
	if mapExtraType[MJ_HU_EXTRA_TYPE_GANGSHANGHUA] {
		hu += 8
	}
	if mapExtraType[MJ_HU_EXTRA_TYPE_QIANGGANG] {
		hu += 8
	}
	if mapSubType[MJ_HU_SUB_TYPE_TWOANGANG] {
		hu += 8
	}

	if mapSubType[MJ_HU_SUB_TYPE_DUIDUIHU] {
		hu += 6
	}
	if mapSubType[MJ_HU_SUB_TYPE_HUNYISE] {
		hu += 6
	}
	if mapSubType[MJ_HU_SUB_TYPE_SANSESANBU] {
		hu += 6
	}
	if mapSubType[MJ_HU_SUB_TYPE_WUMENQI] {
		hu += 6
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANQIUREN] {
		hu += 6
	}
	if mapSubType[MJ_HU_SUB_TYPE_SHUANGJIANKE] {
		hu += 6
	}

	if mapSubType[MJ_HU_SUB_TYPE_DAI19] {
		hu += 4
	}
	if mapSubType[MJ_HU_SUB_TYPE_BUQIUREN] {
		hu += 4
	}
	if mapSubType[MJ_HU_SUB_TYPE_TWOMINGGANG] {
		hu += 4
	}
	if mapSubType[MJ_HU_SUB_TYPE_HEJUEZHANG] {
		hu += 4
	}

	if mapSubType[MJ_HU_SUB_TYPE_JIANKE] {
		hu += 2
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUANFENGKE] {
		hu += 2
	}
	if mapSubType[MJ_HU_SUB_TYPE_MENFENGKE] {
		hu += 2
	}
	if mapSubType[MJ_HU_SUB_TYPE_MENQING] {
		hu += 2
	}
	if mapSubType[MJ_HU_SUB_TYPE_PINGHU] {
		hu += 2
	}
	if mapSubType[MJ_HU_SUB_TYPE_SIGUIYI] {
		hu += 2
	}
	hu += 2 * mapHuNum[MJ_HU_NUM_TYPE_SAME_KE]
	if mapSubType[MJ_HU_SUB_TYPE_TWOANKE] {
		hu += 2
	}
	if mapSubType[MJ_HU_SUB_TYPE_ANGANG] {
		hu += 2
	}
	if mapSubType[MJ_HU_SUB_TYPE_ZHONGZHANG] {
		hu += 2
	}

	hu += mapHuNum[MJ_HU_NUM_TYPE_YIBANGAO]
	hu += mapHuNum[MJ_HU_NUM_TYPE_XIXIANGFENG]
	hu += mapHuNum[MJ_HU_NUM_TYPE_LIANLIU]
	hu += mapHuNum[MJ_HU_NUM_TYPE_LAOSHAOFU]
	hu += mapHuNum[MJ_HU_NUM_TYPE_KE_GANG_ZI19]
	if mapSubType[MJ_HU_SUB_TYPE_MINGGANG] {
		hu += 1
	}
	if mapSubType[MJ_HU_SUB_TYPE_QUEYIMEN] {
		hu += 1
	}
	if mapSubType[MJ_HU_SUB_TYPE_WUZI] {
		hu += 1
	}
	if mapSubType[MJ_HU_SUB_TYPE_BIANZHANG] {
		hu += 1
	}
	if mapSubType[MJ_HU_SUB_TYPE_QIAZI] {
		hu += 1
	}
	if mapSubType[MJ_HU_SUB_TYPE_DANDIAO] {
		hu += 1
	}
	if mapExtraType[MJ_HU_EXTRA_TYPE_ZIMO] {
		hu += 1
	}
	hu += huaNum
	return hu
}

func (c *GUOBIAO_MJ_HuHandler) convertHuInfo(huInfo *HuInfo, mapHuType map[uint64]bool, mapSubType map[uint32]bool, mapExtraType map[uint64]bool) {
	huInfo.HuType = 0
	huInfo.ClearSubType()
	huInfo.ExtraType = 0
	for k, _ := range mapHuType {
		huInfo.HuType |= k
	}
	for k, _ := range mapSubType {
		huInfo.AddSubType(k)
	}
	for k, _ := range mapExtraType {
		huInfo.ExtraType |= k
	}
}

func (c *GUOBIAO_MJ_HuHandler) getGangType(chiCards []*ChiCard) (uint32, uint32) {
	ming, an := 0, 0
	for _, v := range chiCards {
		if v.CardType == MJ_CHI_GANG || v.CardType == MJ_CHI_GANG_WAN {
			ming++
		} else if v.CardType == MJ_CHI_GANG_AN {
			an++
		}
	}
	t1, t2 := uint32(0), uint32(0)
	if ming+an == 4 {
		t1 = MJ_HU_SUB_TYPE_SIGANG
		if an == 2 {
			t2 = MJ_HU_SUB_TYPE_TWOANGANG
		} else if an == 1 {
			t2 = MJ_HU_SUB_TYPE_ANGANG
		}
	} else if ming+an == 3 {
		t1 = MJ_HU_SUB_TYPE_SANGANG
		if an == 2 {
			t2 = MJ_HU_SUB_TYPE_TWOANGANG
		} else if an == 1 {
			t2 = MJ_HU_SUB_TYPE_ANGANG
		}
	} else if ming+an == 2 {
		if an == 2 {
			t1 = MJ_HU_SUB_TYPE_TWOANGANG
		} else {
			t1 = MJ_HU_SUB_TYPE_TWOMINGGANG
			if an == 1 {
				t2 = MJ_HU_SUB_TYPE_ANGANG
			}
		}
	} else if ming == 1 {
		t1 = MJ_HU_SUB_TYPE_MINGGANG
	} else if an == 1 {
		t1 = MJ_HU_SUB_TYPE_ANGANG
	}
	return t1, t2
}
