package mahjong

const (
	RULE_MAHJONG_GUOBIAO int32 = 1001
	RULE_SC_MAHJONG_XUELIU int32 = 1002
	RULE_SC_MAHJONG_XUEZHAN int32 = 1003
	RULE_HN_MAHJONG_HONGZHONG int32 = 1004
)

const (
	MJ_CHECK_HU_FLAG_GROUP     uint32 = uint32(0x01) << 0
)

const (
	COLOR_WAN   int32 = 1
	COLOR_TONG  int32 = 2
	COLOR_TIAO  int32 = 3
	COLOR_OTHER int32 = 4
)

const (
	MAHJONG_1 int32 = 1
	MAHJONG_2 int32 = 2
	MAHJONG_3 int32 = 3
	MAHJONG_4 int32 = 4
	MAHJONG_5 int32 = 5
	MAHJONG_6 int32 = 6
	MAHJONG_7 int32 = 7
	MAHJONG_8 int32 = 8
	MAHJONG_9 int32 = 9

	MAHJONG_DONG      int32 = 11
	MAHJONG_NAN       int32 = 12
	MAHJONG_XI        int32 = 13
	MAHJONG_BEI       int32 = 14
	MAHJONG_HONGZHONG int32 = 15
	MAHJONG_LVFA      int32 = 16
	MAHJONG_BAI       int32 = 17

	MAHJONG_SPRING int32 = 20
	MAHJONG_SUMMER int32 = 21
	MAHJONG_AUTUMN int32 = 22
	MAHJONG_WINTER int32 = 23
	MAHJONG_MEI    int32 = 24
	MAHJONG_LAN    int32 = 25
	MAHJONG_ZHU    int32 = 26
	MAHJONG_JU     int32 = 27	

	MAHJONG_ANY int32 = 50
)

const (
	MAHJONG_MASK int32 = 100
	MAHJONG_LZ   int32 = 1000
)

const (
	MJ_HU_TYPE_NONE         uint64 = 0
	MJ_HU_TYPE_COMMON       uint64 = uint64(0x01) << 0
	MJ_HU_TYPE_QIDUI        uint64 = uint64(0x01) << 1
)


const (	
	ROBOT_LEVEL_AMATEUR int32 = 1 //业余: 按权重打，按胡杠碰吃优先级选择，与普通玩家水平相当
	ROBOT_LEVEL_MAJOR   int32 = 2 //专业: 按向听数打，智能选择胡杠碰吃过，与麻将比赛选手水平相当	
)

//过胡杠碰吃
const (
	MJ_CHI_PASS        int32 = 0
	MJ_CHI_HU          int32 = 1
	MJ_CHI_GANG        int32 = 2
	MJ_CHI_GANG_WAN    int32 = 3
	MJ_CHI_GANG_AN     int32 = 4
	MJ_CHI_PENG        int32 = 5
	MJ_CHI_CHI         int32 = 6	
)