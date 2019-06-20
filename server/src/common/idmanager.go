package common

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"math"
	"sync"
	"time"
)

const (
	GO_BEGIN_TIME = "2006-01-02 15:04:05" //go开始时间
	ID_BEGIN_TIME = "2018-01-01 00:00:00" //ID开始时间
)

//@ Title snowflake id divider
type IdDivider struct {
	lastTime int64
	curIndex int64
	lock     *sync.Mutex
	regionId int64 //region id
	serverId int64 //server id
}

func NewIdDivider(regionId, serverId int64) *IdDivider {
	return &IdDivider{
		regionId: regionId,
		serverId: serverId,
		lock:     &sync.Mutex{},
	}
}

const (
	BITS_TIMESTAMP = 42
	BITS_REGION_ID = 5
	BITS_SERVER_ID = 5
	BITS_ID        = 12

	BITS_SHIFT_TIMESTAMP = BITS_REGION_ID + BITS_SERVER_ID + BITS_ID
	BITS_SHIFT_REGION_ID = BITS_SERVER_ID + BITS_ID
	BITS_SHIFT_SERVER_ID = BITS_ID
	BITS_SHIFT_ID        = 0
)

var (
	ID_MASK_TIMESTAMP = uint64(math.Pow(2, BITS_TIMESTAMP)-1) << BITS_SHIFT_TIMESTAMP
	ID_MASK_REGION    = uint64(math.Pow(2, BITS_REGION_ID)-1) << BITS_SHIFT_REGION_ID
	ID_MASK_SERVER    = uint64(math.Pow(2, BITS_SERVER_ID)-1) << BITS_SHIFT_SERVER_ID
	ID_MASK_ID        = uint64(math.Pow(2, BITS_ID)-1) << BITS_SHIFT_ID
)

//@ Title generate unique id
//format:
// |      timestamp  (42bits, can used for 139 years) |  region id (4bits, 16 regions at most)  |  server id(10bits, 1024 servers at most)  |    id index  (8bits, 256 identities at most)     |
// |             相对时间，单位：毫秒        |               区域标识             |             服务器标识                   |                  id 计数器                      |
func (this *IdDivider) GenId(beginTimeMilliSecs int64) (uid uint64, err error) {
	//check region id
	if this.regionId > int64(math.Pow(2, BITS_SHIFT_REGION_ID)-1) ||
		this.regionId < 0 {
		//
		return 0, errors.New(fmt.Sprintf("region id(%d) invalid, valid range  [0, %d]", this.regionId, int64(math.Pow(2, BITS_SHIFT_REGION_ID)-1)))
	}

	//check server id
	if this.serverId > int64(math.Pow(2, BITS_SHIFT_SERVER_ID)-1) ||
		this.serverId < 0 {
		//
		return 0, errors.New(fmt.Sprintf("server id(%d) invalid, valid range  [0, %d]", this.serverId, int64(math.Pow(2, BITS_SHIFT_SERVER_ID)-1)))
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	timestampPart := (time.Now().UnixNano()/int64(time.Millisecond) - beginTimeMilliSecs) << BITS_SHIFT_TIMESTAMP
	regionPart := this.regionId << BITS_SHIFT_REGION_ID
	serverPart := this.serverId << BITS_SHIFT_SERVER_ID

	if this.lastTime != timestampPart {
		if this.lastTime > timestampPart { //check clock move backwards.
			return 0, fmt.Errorf("lastTime %d > timestampPart %d, maybe clock is moving backwards. Rejecting requests until %d", this.lastTime, timestampPart, this.lastTime)
		} else {
			this.curIndex = 0
			this.lastTime = timestampPart
		}
	} else {
		this.curIndex++
		//check if is overflow
		if this.curIndex > int64(math.Pow(2, BITS_ID)-1) {
			return 0, fmt.Errorf("id index %d is overflow.  valid range [0, %d]", this.curIndex, int64(math.Pow(2, BITS_ID)-1))
		}
	}

	return uint64(timestampPart | regionPart | serverPart | this.curIndex), nil
}

// @Description parse timestamp into readable string.
func IdTimestamp(timestamp int64) string {
	tm := time.Unix(int64((timestamp+mgr.beginTimeMilliSecs)/1000), 0)
	return tm.Format("2006-01-02 03:04:05 PM\n")
}

// @Description decode id to parts.
func IdDecode(id uint64) (timestamp, regionId, serverId, index int64) {
	timestamp = int64((id & ID_MASK_TIMESTAMP) >> BITS_SHIFT_TIMESTAMP)
	regionId = int64((id & ID_MASK_REGION) >> BITS_SHIFT_REGION_ID)
	serverId = int64((id & ID_MASK_SERVER) >> BITS_SHIFT_SERVER_ID)
	index = int64((id & ID_MASK_ID) >> BITS_SHIFT_ID)
	return
}

//@ Title id manager
type IdManager struct {
	beginTimeMilliSecs int64
	dividers           []*IdDivider
}

func NewIdManager(num int, regionId, serverId int64) (mgr *IdManager) {
	mgr = &IdManager{
		dividers: make([]*IdDivider, 0, num),
	}

	for i := 0; i < num; i++ {
		mgr.dividers = append(mgr.dividers, NewIdDivider(regionId, serverId))
	}

	return
}

//@ Title init id manager
func (this *IdManager) Init() (err error) {
	beginTime, err := time.Parse(GO_BEGIN_TIME, ID_BEGIN_TIME)
	if err != nil {
		return
	}
	this.beginTimeMilliSecs = beginTime.UnixNano() / int64(time.Millisecond)
	return
}

var mgr *IdManager

func IdManagerInit(num int, regionId, serverId int64) error {
	mgr = NewIdManager(num, regionId, serverId)
	return mgr.Init()
}

//generate id by type
func IdManagerGen(idType int) (id uint64, err error) {
	if idType < 0 || idType >= len(mgr.dividers) {
		return 0, errors.New(fmt.Sprintf("wrong id type %d", idType))
	}

	return mgr.dividers[idType].GenId(mgr.beginTimeMilliSecs)
}

func IdMgrInit(idDivderSize int) (err error) {
	var regionId, serverId int64
	regionId, err = beego.AppConfig.Int64("RegionId")
	if err != nil {
		panic("no specific RegionId")
	}

	serverId, err = beego.AppConfig.Int64("ServerId")
	if err != nil {
		panic("no specific ServerId")
	}

	err = IdManagerInit(idDivderSize, regionId, serverId)

	return
}

// get region id of specific id.
func GetRegionId(id, regionNum uint64) uint64 {
	return id % regionNum
}
