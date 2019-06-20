package controllers

/*
import (
	"common"
	"errors"
	"fmt"
	"sort"
	"sync"
	admindao "utils/admin/dao"
	"utils/admin/models"
	adminmodels "utils/admin/models"
)

type AgentCommissionController struct {
	BaseController
}

// chip 2 commission config
type Chip2CommissionConfigItem struct {
	Commission int32
	Precision  int
	Id         int64
}

type Chip2CommissionConfig struct {
	items    map[int32]Chip2CommissionConfigItem
	sortKeys []int32
	lock     sync.RWMutex
}

func NewChip2CommissionConfig() *Chip2CommissionConfig {
	return &Chip2CommissionConfig{
		items: make(map[int32]Chip2CommissionConfigItem),
	}
}

var Chip2CommissionConfigMgr = NewChip2CommissionConfig()

func (c *Chip2CommissionConfig) Load() (err error) {
	var list []models.Commissionrates
	list, err = admindao.CommissionRateDaoEntity.All()
	if err != nil {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.items = make(map[int32]Chip2CommissionConfigItem)
	c.sortKeys = []int32{}

	for _, config := range list {
		c.items[int32(config.Min)] = Chip2CommissionConfigItem{
			Commission: int32(config.Commission),
			Precision:  int(config.Precision),
			Id:         config.Id,
		}
		c.sortKeys = append(c.sortKeys, int32(config.Min))
	}

	// sort it
	sort.Slice(c.sortKeys, func(i, j int) bool {
		return c.sortKeys[i] < c.sortKeys[j]
	})

	return
}

//
func (c *Chip2CommissionConfig) Commission(chips int32) (commission int32, precision int, found bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	//
	for _, key := range c.sortKeys {
		if key > chips {
			break
		}

		commission = c.items[key].Commission
		precision = c.items[key].Precision
		found = true
	}

	return
}

func (c *Chip2CommissionConfig) CommissionById(id int64) (commission int32, precision int) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, item := range c.items {
		if item.Id == id {
			commission = item.Commission
			precision = item.Precision
			break
		}
	}

	return
}

func (c *Chip2CommissionConfig) Print() {
	c.lock.RLock()
	defer c.lock.RUnlock()

	s := ""

	for _, k := range c.sortKeys {
		if item, ok := c.items[k]; ok {
			s += fmt.Sprintf("%+v\n", item)
		}
	}

	common.LogFuncDebug(s)

	return
}

const (
	CALC_COMMISSION_PER_PAGE = 50
)

type WhiteListCommission struct {
	Commission uint32
	Precision  int
}

type WhiteListCommissionConfig struct {
	items map[uint32]WhiteListCommission
	lock  sync.RWMutex
}

func NewWhiteListCommissionConfig() *WhiteListCommissionConfig {
	return &WhiteListCommissionConfig{
		items: make(map[uint32]WhiteListCommission),
	}
}

var WhiteListCommissionConfigMgr = NewWhiteListCommissionConfig()

func (c *WhiteListCommissionConfig) Load() (err error) {
	var list []adminmodels.AgentWhiteList
	list, err = admindao.AgentWhiteListDaoEntity.All()
	if err != nil {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	c.items = make(map[uint32]WhiteListCommission)

	for _, l := range list {
		c.items[l.Id] = WhiteListCommission{
			Commission: uint32(l.Commission),
			Precision:  int(l.Precision),
		}
	}

	return
}

func (c *WhiteListCommissionConfig) Get(id uint32) (commission int32, precision int, ok bool) {
	var v WhiteListCommission
	c.lock.RLock()
	defer c.lock.RUnlock()

	if v, ok = c.items[id]; ok {
		commission = int32(v.Commission)
		precision = v.Precision
		return
	}

	return
}


// calc users' commission
func calcCommission(timestamp int64) (err error) {
	start, end := common.TimestampToBeginAndEndString(timestamp)
	//timestamp, _, start, end := common.TheOtherDayTimeRange(-1)
	now := int64(common.NowUint32())

	// record log
	id, err := otcdao.CommissionCalcDaoEntity.Add(start, end)
	if err != nil {
		return
	}

	// sync user daily
	err = syncGameUserDaily(timestamp, now, start, end)
	if err != nil {
		otcdao.CommissionCalcDaoEntity.Update(id, otcdao.CommissionStatusFailed, true, fmt.Sprintf("syncGameUserDaily failed : %v", err))
		return
	}

	//calc all users' commission
	var (
		total                   int
		pages                   int
		integer, decimals       int32
		integerTmp, decimalsTmp int32
	)

	total, err = agentdao.AgentPathDaoEntity.QueryTotal()
	if err != nil {
		otcdao.CommissionCalcDaoEntity.Update(id, otcdao.CommissionStatusFailed, true, fmt.Sprintf("AgentPathDaoEntity.QueryTotal failed : %v", err))
		return
	}

	pages = total / CALC_COMMISSION_PER_PAGE
	if total%CALC_COMMISSION_PER_PAGE > 0 {
		pages += 1
	}

	//loop pages
	for page := 1; page <= pages+1; page++ {
		//calc users' commission of one page
		integerTmp, decimalsTmp, err = calcSinglePage(agentdao.CHANNEL_ID_FIRST, timestamp, now, page, CALC_COMMISSION_PER_PAGE)
		if err != nil {
			otcdao.CommissionCalcDaoEntity.Update(id, otcdao.CommissionStatusFailed, true, fmt.Sprintf("calcSinglePage failed : %v", err))
			return
		}

		// add the commission currency
		integer, decimals, err = common.AddCurrency(integer, decimals, integerTmp, decimalsTmp, otccommon.Cursvr.EusdPrecision)
		if err != nil {
			otcdao.CommissionCalcDaoEntity.Update(id, otcdao.CommissionStatusFailed, true, fmt.Sprintf("AddCurrency failed : %v", err))
			return
		}
	}

	// sync channel daily
	var integerDaily, decimalsDaily int32
	integerDaily, decimalsDaily, err = syncGameDaily(timestamp, now, start, end)
	if err != nil {
		common.LogFuncError("%v", err)
		otcdao.CommissionCalcDaoEntity.Update(id, otcdao.CommissionStatusFailed, true, fmt.Sprintf("syncGameDaily failed : %v", err))
		return
	}

	// transform  coin to eusd
	integerDaily, decimalsDaily = common.DivideCurrency(integerDaily, decimalsDaily, otccommon.Cursvr.EusdExchangeRate)

	// game channel
	integerChannel, decimalsChannel := common.MultiplyCurrency(integerDaily, decimalsDaily, otccommon.Cursvr.GameChannelPercentage)

	// for agents and profit
	var (
		integerSum, decimalsSum int32
	)
	integerSum, decimalsSum = common.SubCurrency2(integerDaily, decimalsDaily, integerChannel, decimalsChannel)

	// record profit
	integerProfit, decimalsProfit := common.SubCurrency2(integerSum, decimalsSum, integer, decimals)

	// record commission statistic
	err = otcdao.CommissionStatDaoEntity.Update(timestamp, integerChannel, decimalsChannel, integer, decimals, integerProfit, decimalsProfit)
	if err != nil {
		otcdao.CommissionCalcDaoEntity.Update(id, otcdao.CommissionStatusFailed, true, fmt.Sprintf("save commission stat failed : %v", err))
		return
	}

	otcdao.CommissionCalcDaoEntity.Update(id, otcdao.CommissionStatusDone, true, fmt.Sprintf("daily %d.%d\nchannel %d.%d\nsum %d.%d\ncommission %d.%d\nprofit %d.%d",
		integerDaily, decimalsDaily,
		integerChannel, decimalsChannel,
		integerSum, decimalsSum,
		integer, decimals,
		integerProfit, decimalsProfit))

	return
}

// calc commissions of users in one page
func calcSinglePage(channelId uint32, timestamp, now int64, page, limit int) (integer, decimals int32, err error) {
	// get users' agent path data
	var paths []agentmodels.AgentPath
	paths, err = agentdao.AgentPathDaoEntity.Query(page, limit)
	if err != nil {
		return
	}

	var integerTmp, decimalsTmp int32
	for _, path := range paths {
		integerTmp, decimalsTmp, err = calcSingeCommission(path.Uid, channelId, timestamp, now)
		if err != nil {
			return
		}

		integer, decimals, err = common.AddCurrency(integer, decimals, integerTmp, decimalsTmp, otccommon.Cursvr.EusdPrecision)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
	}

	return
}

// calc singe user's commission
func calcSingeCommission(uid uint64, channelId uint32, timestamp, now int64) (balanceInteger, balanceDecimals int32, err error) {
	// get all low levels
	var paths []agentmodels.AgentPath
	paths, err = agentdao.AgentPathDaoEntity.GetAllLowLevels(uid)
	if err != nil {
		return
	}

	var (
		uids       []uint64
		parentPath agentmodels.AgentPath
	)
	for _, p := range paths {
		uids = append(uids, p.Uid)
		if p.Uid == uid {
			parentPath = p
		}
	}

	//get chips
	var (
		chipsInteger, chipsDecimals int32
		ok                          bool
	)

	chipsInteger, chipsDecimals, err = gamedao.GameUserDailyDaoEntity.Chips(agentdao.CHANNEL_ID_FIRST, timestamp, uids)
	if err != nil {
		return
	}

	common.LogFuncDebug("%d %d.%d sub agent %v", uid, chipsInteger, chipsDecimals, uids)

	// try get balance by whitelist id
	if parentPath.WhiteListId > 0 {
		balanceInteger, balanceDecimals, ok = getCommissionByWhiteListId(parentPath.WhiteListId, chipsInteger, chipsDecimals)
	}

	common.LogFuncDebug("balance %d.%d", balanceInteger, balanceDecimals)

	// try get balance by chips while haven't got by white list id
	if !ok {
		balanceInteger, balanceDecimals, ok = getCommissionByChips(chipsInteger, chipsDecimals)
		if !ok {
			err = fmt.Errorf("getCommissionByChips %d.%d failed !", chipsInteger, chipsDecimals)
			return
		}
	}
	common.LogFuncDebug("balance %d.%d EusdExchangeRate %f", balanceInteger, balanceDecimals, otccommon.Cursvr.EusdExchangeRate)
	// change coin to eusd
	balanceInteger, balanceDecimals = common.DivideCurrency(balanceInteger, balanceDecimals, otccommon.Cursvr.EusdExchangeRate)
	common.LogFuncDebug("balance %d.%d", balanceInteger, balanceDecimals)
	// add commission , no immediately !
	//err = otcdao.AgentDaoEntity.AddBalance(uid, balanceInteger, balanceDecimals)
	//if err != nil {
	//	return
	//}

	common.LogFuncDebug("%d chips %d.%d commission %d.%d whitelistid %d",
		uid, chipsInteger, chipsDecimals, balanceInteger, balanceDecimals, parentPath.WhiteListId)

	// add agent channel commission
	err = agentdao.AgentChannelCommissionDaoEntity.InsertOrUpdate(uid, channelId, balanceInteger, balanceDecimals, timestamp, now)
	if err != nil {
		return
	}

	return
}

func getCommissionByChips(chipsInteger, chipsDecimals int32) (balanceInteger, balanceDecimals int32, found bool) {
	if chipsInteger <= 0 {
		found = true
		return
	}

	// get by chips
	var commission int32
	var precision int
	commission, precision, found = Chip2CommissionConfigMgr.Commission(chipsInteger)
	balanceInteger = subGetCommission(chipsInteger, commission, precision)
	return
}

//
func getCommissionByWhiteListId(whiteListId uint32, chipsInteger, chipsDecimals int32) (balanceInteger, balanceDecimals int32, found bool) {
	if chipsInteger <= 0 {
		found = true
		return
	}
	var (
		commission int32
		precision  int
	)
	commission, precision, found = WhiteListCommissionConfigMgr.Get(whiteListId)
	if found {
		balanceInteger = subGetCommission(chipsInteger, commission, precision)
	}

	return
}

func subGetCommission(chipsInteger int32, commission int32, precision int) int32 {
	return int32((float64(chipsInteger) / float64(precision)) * float64(commission))
}

var (
	errZeroProfit = errors.New("no platform profit")
)

func distributeCommission(timestamp int64) (err error) {

	start, end := common.TimestampToBeginAndEndString(timestamp)

	// record log
	id, err := otcdao.CommissionDistributeDaoEntity.Add(start, end)
	if err != nil {
		return
	}

	// check profit
	var stat otcmodels.CommissionStat
	stat, err = otcdao.CommissionStatDaoEntity.Query(timestamp)
	if err != nil {
		return
	}

	if stat.ProfitInteger <= 0 {
		// to warn
		err = errZeroProfit
		otcdao.CommissionDistributeDaoEntity.Update(id, otcdao.CommissionStatusFailed, true, fmt.Sprintf("%v", err))
		return
	}

	// query total
	var total int
	total, err = agentdao.AgentChannelCommissionDaoEntity.QueryTotal(timestamp)
	if err != nil {
		otcdao.CommissionDistributeDaoEntity.Update(id, otcdao.CommissionStatusFailed, true, fmt.Sprintf("agent channel commission query failed : %v", err))
		return
	}

	//
	const PerPage = 40
	pages := total / PerPage
	if total%PerPage > 0 {
		pages += 1
	}

	// loop pages
	for page := 1; page <= pages; page++ {
		// page set to 1 cause of the rest set of no distribute ones changes every time.
		// we process "the first page" all the time.
		err = distributeSinglePageCommission(timestamp, 1, PerPage)
		if err != nil {
			otcdao.CommissionDistributeDaoEntity.Update(id, otcdao.CommissionStatusFailed, true, fmt.Sprintf("distribute single page commission failed : %v", err))
			return
		}
	}

	// set stat's status to distributed
	err = otcdao.CommissionStatDaoEntity.SetDistributed(timestamp)
	if err != nil {
		otcdao.CommissionDistributeDaoEntity.Update(id, otcdao.CommissionStatusFailed, true, fmt.Sprintf("set commission stat failed : %v", err))
		return
	}

	otcdao.CommissionDistributeDaoEntity.Update(id, otcdao.CommissionStatusDone, true, "")

	return
}

func distributeSinglePageCommission(timestamp int64, page, limit int) (err error) {
	// get page
	var list []agentdao.AgentCommission
	list, err = agentdao.AgentChannelCommissionDaoEntity.Query(timestamp, page, limit)
	if err != nil {
		return
	}

	//
	for _, v := range list {
		err = distributeSingleCommission(v, timestamp)
		if err != nil {
			continue
		}
	}

	return
}

func distributeSingleCommission(commission agentdao.AgentCommission, timestamp int64) (err error) {
	// set agent balance
	err = agentdao.AgentDaoEntity.AddBalance(commission.Uid, commission.Integer, commission.Decimals)
	if err != nil {
		return
	}

	// set commission status to sent
	err = agentdao.AgentChannelCommissionDaoEntity.UpdateStatusToSent(commission.Uid, timestamp)
	if err != nil {
		return
	}

	return
}*/
