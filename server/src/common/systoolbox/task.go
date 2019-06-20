package systoolbox

import (
	"common"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/streadway/amqp"
	"log"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// bounds provides a range of acceptable values (plus a map of name to value).
type bounds struct {
	min, max uint
	names    map[string]uint
}

// The bounds for each field.
var (
	seconds = bounds{0, 59, nil}
	minutes = bounds{0, 59, nil}
	hours   = bounds{0, 23, nil}
	days    = bounds{1, 31, nil}
	months  = bounds{1, 12, map[string]uint{
		"jan": 1,
		"feb": 2,
		"mar": 3,
		"apr": 4,
		"may": 5,
		"jun": 6,
		"jul": 7,
		"aug": 8,
		"sep": 9,
		"oct": 10,
		"nov": 11,
		"dec": 12,
	}}
	weeks = bounds{0, 6, map[string]uint{
		"sun": 0,
		"mon": 1,
		"tue": 2,
		"wed": 3,
		"thu": 4,
		"fri": 5,
		"sat": 6,
	}}
)

const (
	// Set the top bit if a star was included in the expression.
	starBit = 1 << 63
)

// Schedule time taks schedule
type Schedule struct {
	Second uint64
	Minute uint64
	Hour   uint64
	Day    uint64
	Month  uint64
	Week   uint64
}

// TaskFunc task func type
type TaskFunc func() error

//type StoreResultFunc func(name string, code int, detail string)

const (
	ResultSuccessful = iota
	ResultSkip
	ResultFailure
)

//easyjson:json
type TaskMsg struct {
	Id       uint32 `json:"id"`
	Name     string `json:"name"`
	Spec     string `json:"spec"`
	FuncName string `json:"func_name"`
	Switch   uint8  `json:"switch"`
}

//easyjson:json
type TaskMsgList struct {
	List []TaskMsg `json:"list"`
}

//easyjson:json
type TaskResult struct {
	AppName   string `json:"app_name"`
	RegionId  int64  `json:"region_id"`
	ServerId  int64  `json:"server_id"`
	Name      string `json:"name"`
	Code      int32  `json:"code"`
	Detail    string `json:"detail"`
	BeginTime uint32 `json:"begin_time"`
	EndTime   uint32 `json:"end_time"`
	Ctime     uint32 `json:"ctime"`
}

//easyjson:json
type TaskDetailItem struct {
	Name   string `json:"name"`
	Status int    `json:"status"`
	Switch int    `json:"switch"`
	Spec   string `json:"spec"`
	Prev   int64  `json:"prev"`
	Next   int64  `json:"next"`
}

//easyjson:json
type TaskDetail struct {
	AppName  string           `json:"app_name"`
	RegionId int64            `json:"region_id"`
	ServerId int64            `json:"server_id"`
	Items    []TaskDetailItem `json:"list"`
}

var TaskMgr = NewTaskManager()

const (
	Unchanged int32 = iota //tasks unchanged
	Changed                //tasks changed
)

const (
	Unstarted int32 = iota //taskmgr no started
	Started                //taskmgr started
)

// task manager
type TaskManager struct {
	AdminTaskList map[string]*Task
	stop          chan bool
	changed       chan bool
	lock          sync.RWMutex

	Changed, Started int32

	appName            string
	regionId, serverId int64
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		AdminTaskList: make(map[string]*Task),
		stop:          make(chan bool),
		changed:       make(chan bool),
	}
}

// StartTask start all tasks
func (m *TaskManager) Start(appName string, regionId, serverId int64) {
	m.appName = appName
	m.regionId = regionId
	m.serverId = serverId

	//m.lock.Lock()
	//defer m.lock.Unlock()

	if m.isStarted() {
		//If already started， no need to start another goroutine.
		return
	}

	m.setStarted(true)

	//check changed
	go func() {
		for {
			if m.isChanged() { //if changed, send to main routine
				m.changed <- true
				m.setChanged(false) //reset
			}
			time.Sleep(time.Second * 10)
		}
	}()

	// main routine
	go common.SafeRun(m.run)()
}

func (m *TaskManager) isStarted() bool {
	return Started == atomic.LoadInt32(&m.Started)
}

func (m *TaskManager) setStarted(isStarted bool) {
	if isStarted {
		atomic.StoreInt32(&m.Started, Started)
	} else {
		atomic.StoreInt32(&m.Started, Unstarted)
	}
}

func (m *TaskManager) isChanged() bool {
	return Changed == atomic.LoadInt32(&m.Changed)
}

func (m *TaskManager) setChanged(isStarted bool) {
	if isStarted {
		atomic.StoreInt32(&m.Changed, Changed)
	} else {
		atomic.StoreInt32(&m.Changed, Unchanged)
	}
}

func (m *TaskManager) SetAllNext(now time.Time) {
	//m.lock.RLock()
	//defer m.lock.RUnlock()

	// set all tasks' next
	for _, t := range m.AdminTaskList {
		t.SetNext(now)
	}
}

func (m *TaskManager) TaskList() map[string]*Task {
	//m.lock.RLock()
	//defer m.lock.RUnlock()

	return m.AdminTaskList
}

func (m *TaskManager) run() {
	now := time.Now().Local()

	m.SetAllNext(now)

	for {
		common.LogFuncDebug("task poll begin...")
		AdminTaskList := m.TaskList()
		sortList := NewMapSorter(AdminTaskList)
		sortList.Sort()
		var effective time.Time
		common.LogFuncDebug("task poll len(AdminTaskList) %d, sortList.Vals %v", len(AdminTaskList), sortList.Vals)
		if len(AdminTaskList) == 0 || sortList.Vals[0].GetNext().IsZero() {
			// If there are no entries yet, just sleep - it still handles new entries
			// and stop requests.
			effective = now.AddDate(10, 0, 0)
		} else {
			effective = sortList.Vals[0].GetNext()
			common.LogFuncDebug("task poll get new effective %v by task %s", effective.String(), sortList.Vals[0].Taskname)
		}

		common.LogFuncDebug("task poll will awake in %v secs",
			effective.Sub(now).Seconds())

		select {
		case now = <-time.After(effective.Sub(now)):
			// Run every entry whose next time was this effective time.
			for _, e := range sortList.Vals {
				tmp := e
				common.LogFuncDebug("effective %v, task %s next %v", effective.String(), tmp.Taskname, tmp.Next)
				if tmp.Next != effective {
					break
				}

				//common.LogFuncDebug("start run task %s this time %v ", e.Taskname, e.Next.String())
				go common.SafeRun(func() {
					tmp.Run(now)
				})()

				tmp.Prev = tmp.Next
				tmp.Next = tmp.Spec.Next(now)
				common.LogFuncDebug("task %s 's next set to %v", tmp.Taskname, tmp.Next.String())
			}
			common.LogFuncDebug("task poll 1 poll finished, to be continue...")
			continue
		case <-m.changed:
			now = time.Now().Local()
			for _, t := range m.TaskList() {
				// skip disable ones
				if t.GetSwitch() == TaskSwitchDisable {
					continue
				}
				t.Next = t.Spec.Next(now)
			}
			common.LogFuncDebug("task poll changed, poll refresh, to be continue...")
			continue
		case <-m.stop: //stop the whole manager
			//common.LogFuncDebug("task poll stop...")
			return
		}
	}
}

// StopTask stop all tasks
func (m *TaskManager) Stop() {
	//m.lock.Lock()
	//defer m.lock.Unlock()

	if m.isStarted() {
		m.setStarted(false)
		m.stop <- true
	}

}

// PauseTask pause the task, but no stop the running logic
func (m *TaskManager) DisableTask(taskName string) {
	for _, t := range m.TaskList() {
		if t.GetName() == taskName {
			t.SetSwitch(TaskSwitchDisable)
		}
	}
}

func (m *TaskManager) AddTaskList(list []*Task) {
	//m.lock.Lock()
	//defer m.lock.Unlock()

	isChanged := false

	for _, t := range list {
		if tExist, ok := m.AdminTaskList[t.GetName()]; ok {
			if tExist.tryReset(t) {
				isChanged = true
			}
		} else {
			t.Switch = TaskSwitchEnable
			m.AdminTaskList[t.GetName()] = t
			isChanged = true
		}
		common.LogFuncDebug("add task : %v", t)
	}

	if isChanged && m.isStarted() {
		m.setChanged(true)
	}
}

func (m *TaskManager) SwitchTaskList(list []*Task) {
	m.changeSwitchTaskList(list)
}

//func (m *TaskManager) DisableTaskList(list []*Task) {
//	m.changeSwitchTaskList(list, TaskSwitchDisable)
//}

func (m *TaskManager) changeSwitchTaskList(list []*Task) {
	isChanged := false

	for _, t := range list {
		if v, ok := m.AdminTaskList[t.Taskname]; ok {
			v.SetSwitch(t.Switch)
			isChanged = true
		}
		common.LogFuncDebug("switch task[%s] to %d", t.Taskname, t.Switch)
	}

	if isChanged && m.isStarted() {
		m.setChanged(true)
	}
}

// DeleteTask delete task with name, but no stop the running logic
func (m *TaskManager) DeleteTask(taskname string) {
	//m.lock.Lock()
	//defer m.lock.Unlock()

	delete(m.AdminTaskList, taskname)
	if m.isStarted() {
		m.setChanged(true)
	}
}

// RunTaskList run task right away
func (m *TaskManager) RunTaskList(list []*Task) {
	now := time.Now().Local()
	for _, t := range list {
		if task, ok := m.AdminTaskList[t.GetName()]; ok {
			task.Run(now)
		}
	}
}

// DeleteTask delete task with name, but no stop the running logic
func (m *TaskManager) DeleteTaskList(list []*Task) {
	//m.lock.Lock()
	//defer m.lock.Unlock()

	isChanged := false
	for _, t := range list {
		delete(m.AdminTaskList, t.GetName())
		isChanged = true
	}

	if m.isStarted() && isChanged {
		m.setChanged(true)
	}
}

// EditTask edit task, but no stop the running logic.
func (m *TaskManager) EditTask(taskname string, t Task) {
	// delete and add.
	//m.lock.Lock()
	//defer m.lock.Unlock()

	delete(m.AdminTaskList, taskname)
	m.AdminTaskList[t.GetName()] = &t

	if m.isStarted() {
		m.setChanged(true)
	}
}

func (m *TaskManager) Detail() (detail TaskDetail) {
	detail.AppName = m.appName
	detail.RegionId = m.regionId
	detail.ServerId = m.serverId
	for _, t := range m.AdminTaskList {
		detail.Items = append(detail.Items, t.detail())
	}
	return
}

const (
	StatusIdle = iota
	StatusRunning
)

var StatusDetail = map[int]string{
	StatusIdle:    "Idle",
	StatusRunning: "Running",
}

const (
	TaskSwitchDisable = iota
	TaskSwitchEnable
)

// Task task struct
type Task struct {
	Taskname string
	Spec     *Schedule
	SpecStr  string
	Prev     time.Time
	Next     time.Time
	DoFunc   TaskFunc
	SRFunc   StoreResultFunc
	Switch   int32 //enable/disable
	Status   int32 //running/idle
	lock     sync.Mutex
}

// NewTask add new task with name, time and func
func NewTask(tname string, spec string, f TaskFunc, srf StoreResultFunc) *Task {
	if srf == nil {
		srf = taskStoreResult
	}

	task := &Task{
		Taskname: tname,
		DoFunc:   f,
		SpecStr:  spec,
		SRFunc:   srf,
	}
	task.SetCron(spec)
	return task
}

// GetName get task name
func (t *Task) GetName() string {
	//t.lock.RLock()
	//defer t.lock.RUnlock()

	return t.Taskname
}

// GetSpec get spec string
func (t *Task) GetSpec() string {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.SpecStr
}

func (t *Task) tryReset(tInput *Task) (ok bool) {
	t.lock.Lock()
	defer t.lock.Unlock()

	if tInput.SpecStr != t.SpecStr {
		t.SpecStr = tInput.SpecStr
		t.Spec = t.parse(t.SpecStr)
		ok = true
	}

	return
}

// GetStatus get current task status
func (t *Task) GetSwitch() int32 {
	return atomic.LoadInt32(&t.Switch)
}

// SetSwitch set task switch
func (t *Task) SetSwitch(s int32) {
	atomic.StoreInt32(&t.Switch, s)
}

// Run run all tasks
func (t *Task) Run(now time.Time) (err error) {
	beginTime := common.NowUint32()
	endTime := uint32(0)

	// if running, do nothing
	if atomic.LoadInt32(&t.Status) == StatusRunning ||
		atomic.LoadInt32(&t.Switch) == TaskSwitchDisable {
		taskStoreResult(t.Taskname, ResultSkip, fmt.Sprint(t), beginTime, endTime)
		common.LogFuncDebug("skip task %v", t.Taskname)
		return
	}

	common.LogFuncDebug("task %v start to run", t.Taskname)

	atomic.StoreInt32(&t.Status, StatusRunning)
	defer func() { atomic.StoreInt32(&t.Status, StatusIdle) }()

	err = t.DoFunc()
	if t.SRFunc != nil {
		endTime = common.NowUint32()
		if err == nil {
			t.SRFunc(t.Taskname, ResultSuccessful, "", beginTime, endTime)
		} else {
			t.SRFunc(t.Taskname, ResultFailure, fmt.Sprint(t, err), beginTime, endTime)
		}
	}

	return
}

// SetNext set next time for this task
func (t *Task) SetNext(now time.Time) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.Next = t.Spec.Next(now)
	common.LogFuncDebug("task[%s] next call time %s", t.Taskname, t.Next.String())
}

// GetNext get the next call time of this task
func (t *Task) GetNext() time.Time {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.Next
}

// SetPrev set prev time of this task
func (t *Task) SetPrev(now time.Time) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.Prev = now
}

// GetPrev get prev time of this task
func (t *Task) GetPrev() time.Time {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.Prev
}

// six columns mean：
//       second：0-59
//       minute：0-59
//       hour：1-23
//       day：1-31
//       month：1-12
//       week：0-6（0 means Sunday）

// SetCron some signals：
//       *： any time
//       ,：　 separate signal
//　　    －：duration
//       /n : do as n times of time duration
/////////////////////////////////////////////////////////
//	0/30 * * * * *                        every 30s
//	0 43 21 * * *                         21:43
//	0 15 05 * * * 　　                     05:15
//	0 0 17 * * *                          17:00
//	0 0 17 * * 1                           17:00 in every Monday
//	0 0,10 17 * * 0,2,3                   17:00 and 17:10 in every Sunday, Tuesday and Wednesday
//	0 0-10 17 1 * *                       17:00 to 17:10 in 1 min duration each time on the first day of month
//	0 0 0 1,15 * 1                        0:00 on the 1st day and 15th day of month
//	0 42 4 1 * * 　 　                     4:42 on the 1st day of month
//	0 0 21 * * 1-6　　                     21:00 from Monday to Saturday
//	0 0,10,20,30,40,50 * * * *　           every 10 min duration
//	0 */10 * * * * 　　　　　　              every 10 min duration
//	0 * 1 * * *　　　　　　　　               1:00 to 1:59 in 1 min duration each time
//	0 0 1 * * *　　　　　　　　               1:00
//	0 0 */1 * * *　　　　　　　               0 min of hour in 1 hour duration
//	0 0 * * * *　　　　　　　　               0 min of hour in 1 hour duration
//	0 2 8-20/3 * * *　　　　　　             8:02, 11:02, 14:02, 17:02, 20:02
//	0 30 5 1,15 * *　　　　　　              5:30 on the 1st day and 15th day of month
func (t *Task) SetCron(spec string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.Spec = t.parse(spec)
}

func (t *Task) parse(spec string) *Schedule {
	if len(spec) > 0 && spec[0] == '@' {
		return t.parseSpec(spec)
	}
	// Split on whitespace.  We require 5 or 6 fields.
	// (second) (minute) (hour) (day of month) (month) (day of week, optional)
	fields := strings.Fields(spec)
	if len(fields) != 5 && len(fields) != 6 {
		log.Panicf("Expected 5 or 6 fields, found %d: %s", len(fields), spec)
	}

	// If a sixth field is not provided (DayOfWeek), then it is equivalent to star.
	if len(fields) == 5 {
		fields = append(fields, "*")
	}

	schedule := &Schedule{
		Second: getField(fields[0], seconds),
		Minute: getField(fields[1], minutes),
		Hour:   getField(fields[2], hours),
		Day:    getField(fields[3], days),
		Month:  getField(fields[4], months),
		Week:   getField(fields[5], weeks),
	}

	return schedule
}

func (t *Task) parseSpec(spec string) *Schedule {
	switch spec {
	case "@yearly", "@annually":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    1 << days.min,
			Month:  1 << months.min,
			Week:   all(weeks),
		}

	case "@monthly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    1 << days.min,
			Month:  all(months),
			Week:   all(weeks),
		}

	case "@weekly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    all(days),
			Month:  all(months),
			Week:   1 << weeks.min,
		}

	case "@daily", "@midnight":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   1 << hours.min,
			Day:    all(days),
			Month:  all(months),
			Week:   all(weeks),
		}

	case "@hourly":
		return &Schedule{
			Second: 1 << seconds.min,
			Minute: 1 << minutes.min,
			Hour:   all(hours),
			Day:    all(days),
			Month:  all(months),
			Week:   all(weeks),
		}
	}
	log.Panicf("Unrecognized descriptor: %s", spec)
	return nil
}

func (t *Task) detail() TaskDetailItem {
	//t.lock.Lock()
	//defer t.lock.Unlock()

	return TaskDetailItem{
		Name:   t.Taskname,
		Status: int(atomic.LoadInt32(&t.Status)),
		Switch: int(atomic.LoadInt32(&t.Switch)),
		Spec:   t.SpecStr,
		Prev:   t.Prev.Unix(),
		Next:   t.Next.Unix(),
	}
}

// Next set schedule to next time
func (s *Schedule) Next(t time.Time) time.Time {

	// Start at the earliest possible time (the upcoming second).
	t = t.Add(1*time.Second - time.Duration(t.Nanosecond())*time.Nanosecond)

	// This flag indicates whether a field has been incremented.
	added := false

	// If no time is found within five years, return zero.
	yearLimit := t.Year() + 5

WRAP:
	if t.Year() > yearLimit {
		return time.Time{}
	}

	// Find the first applicable month.
	// If it's this month, then do nothing.
	for 1<<uint(t.Month())&s.Month == 0 {
		// If we have to add a month, reset the other parts to 0.
		if !added {
			added = true
			// Otherwise, set the date at the beginning (since the current time is irrelevant).
			t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		}
		t = t.AddDate(0, 1, 0)

		// Wrapped around.
		if t.Month() == time.January {
			goto WRAP
		}
	}

	// Now get a day in that month.
	for !dayMatches(s, t) {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		}
		t = t.AddDate(0, 0, 1)

		if t.Day() == 1 {
			goto WRAP
		}
	}

	for 1<<uint(t.Hour())&s.Hour == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
		}
		t = t.Add(1 * time.Hour)

		if t.Hour() == 0 {
			goto WRAP
		}
	}

	for 1<<uint(t.Minute())&s.Minute == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, t.Location())
		}
		t = t.Add(1 * time.Minute)

		if t.Minute() == 0 {
			goto WRAP
		}
	}

	for 1<<uint(t.Second())&s.Second == 0 {
		if !added {
			added = true
			t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, t.Location())
		}
		t = t.Add(1 * time.Second)

		if t.Second() == 0 {
			goto WRAP
		}
	}

	return t
}

func dayMatches(s *Schedule, t time.Time) bool {
	var (
		domMatch = 1<<uint(t.Day())&s.Day > 0
		dowMatch = 1<<uint(t.Weekday())&s.Week > 0
	)

	if s.Day&starBit > 0 || s.Week&starBit > 0 {
		return domMatch && dowMatch
	}
	return domMatch || dowMatch
}

// MapSorter sort map for tasker
type MapSorter struct {
	Keys []string
	Vals []*Task
}

// NewMapSorter create new tasker map
func NewMapSorter(m map[string]*Task) *MapSorter {
	ms := &MapSorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]*Task, 0, len(m)),
	}
	for k, v := range m {
		ms.Keys = append(ms.Keys, k)
		ms.Vals = append(ms.Vals, v)
	}
	return ms
}

// Sort sort tasker map
func (ms *MapSorter) Sort() {
	sort.Sort(ms)
}

func (ms *MapSorter) Len() int { return len(ms.Keys) }
func (ms *MapSorter) Less(i, j int) bool {
	if ms.Vals[i].GetNext().IsZero() {
		return false
	}
	if ms.Vals[j].GetNext().IsZero() {
		return true
	}
	return ms.Vals[i].GetNext().Before(ms.Vals[j].GetNext())
}
func (ms *MapSorter) Swap(i, j int) {
	ms.Vals[i], ms.Vals[j] = ms.Vals[j], ms.Vals[i]
	ms.Keys[i], ms.Keys[j] = ms.Keys[j], ms.Keys[i]
}

func getField(field string, r bounds) uint64 {
	// list = range {"," range}
	var bits uint64
	ranges := strings.FieldsFunc(field, func(r rune) bool { return r == ',' })
	for _, expr := range ranges {
		bits |= getRange(expr, r)
	}
	return bits
}

// getRange returns the bits indicated by the given expression:
//   number | number "-" number [ "/" number ]
func getRange(expr string, r bounds) uint64 {

	var (
		start, end, step uint
		rangeAndStep     = strings.Split(expr, "/")
		lowAndHigh       = strings.Split(rangeAndStep[0], "-")
		singleDigit      = len(lowAndHigh) == 1
	)

	var extrastar uint64
	if lowAndHigh[0] == "*" || lowAndHigh[0] == "?" {
		start = r.min
		end = r.max
		extrastar = starBit
	} else {
		start = parseIntOrName(lowAndHigh[0], r.names)
		switch len(lowAndHigh) {
		case 1:
			end = start
		case 2:
			end = parseIntOrName(lowAndHigh[1], r.names)
		default:
			log.Panicf("Too many hyphens: %s", expr)
		}
	}

	switch len(rangeAndStep) {
	case 1:
		step = 1
	case 2:
		step = mustParseInt(rangeAndStep[1])

		// Special handling: "N/step" means "N-max/step".
		if singleDigit {
			end = r.max
		}
	default:
		log.Panicf("Too many slashes: %s", expr)
	}

	if start < r.min {
		log.Panicf("Beginning of range (%d) below minimum (%d): %s", start, r.min, expr)
	}
	if end > r.max {
		log.Panicf("End of range (%d) above maximum (%d): %s", end, r.max, expr)
	}
	if start > end {
		log.Panicf("Beginning of range (%d) beyond end of range (%d): %s", start, end, expr)
	}

	return getBits(start, end, step) | extrastar
}

// parseIntOrName returns the (possibly-named) integer contained in expr.
func parseIntOrName(expr string, names map[string]uint) uint {
	if names != nil {
		if namedInt, ok := names[strings.ToLower(expr)]; ok {
			return namedInt
		}
	}
	return mustParseInt(expr)
}

// mustParseInt parses the given expression as an int or panics.
func mustParseInt(expr string) uint {
	num, err := strconv.Atoi(expr)
	if err != nil {
		log.Panicf("Failed to parse int from %s: %s", expr, err)
	}
	if num < 0 {
		log.Panicf("Negative number (%d) not allowed: %s", num, expr)
	}

	return uint(num)
}

// getBits sets all bits in the range [min, max], modulo the given step size.
func getBits(min, max, step uint) uint64 {
	var bits uint64

	// If step is 1, use shifts.
	if step == 1 {
		return ^(math.MaxUint64 << (max + 1)) & (math.MaxUint64 << min)
	}

	// Else, use a simple loop.
	for i := min; i <= max; i += step {
		bits |= 1 << i
	}
	return bits
}

// all returns all bits within the given bounds.  (plus the star bit)
func all(r bounds) uint64 {
	return getBits(r.min, r.max, 1) | starBit
}

// task store result function type
type StoreResultFunc func(name string, code int32, detail string, beginTime, endTime uint32)

//
func taskStoreResult(name string, code int32, detail string, beginTime, endTime uint32) {
	result := TaskResult{
		AppName:   TaskMgr.appName,
		RegionId:  TaskMgr.regionId,
		ServerId:  TaskMgr.serverId,
		Name:      name,
		Code:      code,
		Detail:    detail,
		BeginTime: beginTime,
		EndTime:   endTime,
		Ctime:     common.NowUint32(),
	}
	buf, err := easyjson.Marshal(result)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	err = common.RabbitMQPublish(common.RabbitMQBusinessNameTaskResult,
		common.RabbitMQExchangeTaskResult,
		buf)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
}

// TaskBroadcast process task distribute message
func TaskDistribute(delivery amqp.Delivery, fContainer interface{}) (err error) {
	list := TaskMsgList{}
	err = easyjson.Unmarshal(delivery.Body, &list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	//common.LogFuncDebug("TaskDistribute message received : %s %s %v",
	//	delivery.Exchange, delivery.RoutingKey, list)

	if len(list.List) <= 0 {
		return
	}

	switch delivery.RoutingKey {
	case common.RabbitMQRoutingKeyTaskDistribute:
		{
			if taskList, ok := CheckTaskFunc(list, true, fContainer, taskStoreResult); ok {
				if len(taskList) > 0 {
					TaskMgr.AddTaskList(taskList)
				}
			}
		}
	case common.RabbitMQRoutingKeyTaskSwitch:
		{
			if taskList, ok := CheckTaskFunc(list, false, fContainer, taskStoreResult); ok {
				if len(taskList) > 0 {
					TaskMgr.SwitchTaskList(taskList)
				}
			}
		}
	case common.RabbitMQRoutingKeyTaskDelete:
		{
			if taskList, ok := CheckTaskFunc(list, false, fContainer, taskStoreResult); ok {
				if len(taskList) > 0 {
					TaskMgr.DeleteTaskList(taskList)
				}
			}
		}
	case common.RabbitMQRoutingKeyTaskRun:
		{
			if taskList, ok := CheckTaskFunc(list, false, fContainer, taskStoreResult); ok {
				if len(taskList) > 0 {
					TaskMgr.RunTaskList(taskList)
				}
			}
		}
	default:
		common.LogFuncError("unknown routing key %s", delivery.RoutingKey)
		return
	}

	return
}

// check if function exist
func CheckTaskFunc(list TaskMsgList, isCheckFunc bool, fContainer interface{}, srf StoreResultFunc) (taskList []*Task, ok bool) {
	if isCheckFunc {
		fc := reflect.ValueOf(fContainer)
		for _, t := range list.List {
			fvZero := reflect.Value{}
			fv := fc.MethodByName(t.FuncName)
			if fv == fvZero { //while the function no found, just panic to gain attention of the administrator.
				common.LogFuncError("function cron.FunctionContainer.%s no found", t.FuncName)
				return
			}

			f := func() (err error) {
				fv.Call(nil)
				return nil
			}

			taskList = append(taskList, NewTask(t.Name, t.Spec, f, srf))
		}
	} else { //only need function name
		for _, t := range list.List {
			taskList = append(taskList, &Task{
				Taskname: t.Name,
				Switch:   int32(t.Switch),
			})
		}
	}

	ok = true

	return
}
