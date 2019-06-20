package common

import (
	"fmt"
	"github.com/astaxie/beego/context"
	"math"
	"math/rand"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const (
	DaySeconds = 86400
)

//客户端访问类型
const (
	ClientTypeUnkown = iota
	ClientTypeWeb
	ClientTypeApp
)

//random source string
var RANDOMSTR_BYTES = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Return random string of specific length.
func RandomStr(length int) string {
	result := []byte{}
	for i := 0; i < length; i++ {
		result = append(result, RANDOMSTR_BYTES[rand.Intn(len(RANDOMSTR_BYTES))])
	}
	return string(result)
}

// Return random numbers of specific length.
func RandomNum(length int) string {
	format, randMax := "%0"+fmt.Sprint(length)+"d", math.Pow(10, float64(length))
	return fmt.Sprintf(format, rand.Intn(int(randMax)))
}

type RunFunc func()

type TaskFunc func() error

//防崩溃包装
func SafeRun(runFunc RunFunc) RunFunc {
	return func() {
		defer func() {
			if err, ok := recover().(error); ok {
				LogFuncCritical("%v \n%v", err, string(debug.Stack()))
			}
		}()
		runFunc()
	}
}

//启动一个goroutine，并且进行防崩溃包装
func GoSafeRun(runFunc RunFunc) {
	go SafeRun(runFunc)()
}

//将 "YYYY-MM-DD"格式的日期字符串转换为时间戳
func TimeStrToUint32(str string) (bool, uint32) {
	local, _ := time.LoadLocation("Local")
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str+" 00:00:00", local)
	if err != nil {
		LogFuncError("转换%s为时间戳失败 : %v", str, err)
		return false, 0
	}

	return true, uint32(t.Unix())
}

// @Description get current timestamp. (second)
func NowUint32() uint32 {
	return uint32(time.Now().Unix())
}

// @Description get current timestamp.(millisecond)
func NowInt64MS() int64 {
	return time.Now().UnixNano() / 1000000
}

// @Description get today timestamp range.
// timestamp - base timestamp
// differentDays -  yesterday -1, tomorrow 1
func TheOtherDayTimeRange(differentDays int64) (begin, end int64, beginStr, endStr string) {
	var t time.Time
	if differentDays == 0 {
		t = time.Now()
	} else {
		t = time.Now().Add(time.Duration(differentDays) * time.Hour * 24)
	}
	begin = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	end = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).Unix()
	beginStr = fmt.Sprintf("%.4d%.2d%.2d%.2d%.2d%.2d", t.Year(), t.Month(), t.Day(), 0, 0, 0)
	endStr = fmt.Sprintf("%.4d%.2d%.2d%.2d%.2d%.2d", t.Year(), t.Month(), t.Day(), 23, 59, 59)
	return
}

// @Description begin and end string of timestamp
func TimestampToBeginAndEndString(timestamp int64) (beginStr, endStr string) {
	t := time.Unix(timestamp, 0)
	beginStr = fmt.Sprintf("%.4d%.2d%.2d%.2d%.2d%.2d", t.Year(), t.Month(), t.Day(), 0, 0, 0)
	endStr = fmt.Sprintf("%.4d%.2d%.2d%.2d%.2d%.2d", t.Year(), t.Month(), t.Day(), 23, 59, 59)
	return
}

// @Description get today timestamp range.
func TodayTimeRange() (begin, end int64) {
	begin, end, _, _ = TheOtherDayTimeRange(0)
	return
}

const TimeLayout = "2006-01-02 15:04:05"

var loc, _ = time.LoadLocation("Local")

//// get day time range by month string ,eg : "2019-03"
//func DayTimeRangeByMonthStr(time string) (begin, end int64) {
//	// generate the whole string
//	beginTime, endTime := time + "-01 00:00:00"
//	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
//	sr := theTime.Unix()
//	return
//}

// @Description check if timestamp is in today's time range.
func InTodayTimeRange(timestamp uint32) bool {
	begin, end := TodayTimeRange()

	if int64(timestamp) < end && int64(timestamp) >= begin {
		return true
	}

	return false
}

// @Description check if timestamp is before today.
func BeforeToday(timestamp uint32) bool {
	begin, _ := TodayTimeRange()

	if int64(timestamp) < begin {
		return true
	}

	return false
}

// @Description current time string <YYYYmmDDHHMMSS>.
func NowString() string {
	var now = time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

// @Description current time string mill second <YYYYmmDDHHMMSSmmm>.
func NowStringMillS() string {
	var now = time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d%03d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(),
		time.Duration(now.UnixNano())/time.Millisecond-time.Duration(now.Unix())*(time.Second/time.Millisecond))
}

// @Description get current time string micro second <YYYYmmDDHHMMSSmmmmmm>
func NowStringMicroS() string {
	var now = time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d%06d",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(),
		time.Duration(now.UnixNano())/time.Microsecond)
}

// 格式化传入的时间
func GetTimeFormat(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// 获取timestamp的零点时间
func GetZeroTime(timestamp int64) int64 {
	local, _ := time.LoadLocation("Local")
	timeStr := time.Unix(timestamp, 0).Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr+" 00:00:00", local)
	return t.Unix()
}

// @Description begin and end string of timestamp
func GetBeginAndEndString(timestamp int64) (beginStr, endStr string) {
	begin := GetZeroTime(timestamp)
	end := begin + DaySeconds

	return GetTimeFormat(begin), GetTimeFormat(end)
}

// @Description get remote ip address
func ClientIP(ctx *context.Context) string {
	ip := ctx.Request.Header.Get("X-Forwarded-For")
	if strings.Contains(ip, "127.0.0.1") || ip == "" {
		ip = ctx.Request.Header.Get("X-real-ip")
	}

	if ip == "" {
		return "127.0.0.1"
	}

	return ip
}

//分页
type Meta struct {
	Total     int64 `json:"total"`
	Page      int64 `json:"page"`
	Limit     int64 `json:"limit"`
	TotalPage int64 `json:"total_page"`
	Offset    int64 `json:"offset"`
}

func MakeMeta(total, page, limit int64) (m *Meta) {
	m = &Meta{}
	if total < 1 {
		return
	}

	m.Total = total
	if page < 1 {
		page = 1
	}
	m.Page = page
	if limit < 1 {
		limit = 20
	}
	m.Limit = limit
	m.Offset = limit * (page - 1)
	m.TotalPage = int64(math.Ceil(float64(total) / float64(limit)))

	return
}

//string 转 uint
func StrToUint(strNumber string, value interface{}) (err error) {
	var number interface{}
	number, err = strconv.ParseUint(strNumber, 10, 64)
	switch v := number.(type) {
	case uint64:
		switch d := value.(type) {
		case *uint64:
			*d = v
		case *uint:
			*d = uint(v)
		case *uint16:
			*d = uint16(v)
		case *uint32:
			*d = uint32(v)
		case *uint8:
			*d = uint8(v)
		}
	}
	return
}

// @Description convert time string (eg. "2019/04/12") to timestamp
func TimeStrToUint32Plus(s string) (timestamp uint32, err error) {
	var ErrTimeStrNoExpect = fmt.Errorf("time string [%s] no expect.", s)

	s = strings.Replace(s, "/", "-", -1)

	if ok, tmp := TimeStrToUint32(s); !ok {
		err = ErrTimeStrNoExpect
		return
	} else {
		timestamp = tmp
	}

	return
}

////当前时间加 限定时间 转string
func GetNowTimeStr(add int64) string {
	t := int64(time.Now().Unix() + add)
	timeTemplate1 := "2006-01-02 15:04:05"
	s := time.Unix(t, 0).Format(timeTemplate1)
	return s

}

//

// int64数字除以10000 不丢精度
func Init64DivisorToStr(divisor int64, value float64) string {
	return fmt.Sprintf("%.4f", float64(divisor)/value)
}

//返回n天前零点时间截 n -1 即前一天
func GetPreDateStartTimestamp(n int) (ts int64) {
	t := time.Now().AddDate(0, 0, -1*n)
	ts = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
	return
}

//返回n天前最晚时间截 n -1 即前一天
func GetPreDateOverTimestamp(n int) (ts int64) {
	t := time.Now().AddDate(0, 0, -1*n)
	ts = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location()).Unix()
	return
}

//返回当前月份第一天时间截
func GetMonthDateStartTimestamp(t time.Time) (ts int64) {
	//t := time.Now()
	ts = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location()).Unix()
	return
}

//返回当前月份最后一天时间截
func GetMonthDateOverTimestamp(t time.Time) (ts int64) {
	t1 := t.AddDate(0, 1, 0)
	ts = time.Date(t1.Year(), t1.Month(), t1.Day(), 23, 59, 59, 0, t.Location()).Unix()
	return
}

//得到两个时间段内的所有时间， 跨年
func GetAllMonthRangeTime(start, over int64) (res []map[int][]int) {
	//统计单月还是多个月份
	start_tm := time.Unix(start, 0)
	startYear := start_tm.Year()
	startMonth := int(start_tm.Month())

	over_tm := time.Unix(over, 0)
	overMonth := int(over_tm.Month())
	overYear := over_tm.Year()

	muchYear := overYear - startYear
	//fmt.Println(startYear, startMonth, overYear, overMonth)
	if muchYear > 0 {
		var yearMap map[int][]int
		for n := 0; n <= muchYear; n++ {
			if n == 0 {
				yearMap = make(map[int][]int, overMonth)
				month_nums := overMonth
				for nm := 0; nm < month_nums; nm++ {
					yearMap[overYear] = append(yearMap[overYear], overMonth)
					overMonth--
				}
			} else if startYear == overYear {
				yearMap = make(map[int][]int, 12)
				for nm := 12; nm >= startMonth; nm-- {
					yearMap[overYear] = append(yearMap[overYear], nm)
				}
			} else {
				yearMap = make(map[int][]int, 12)
				for nm := 12; nm >= 1; nm-- {
					yearMap[overYear] = append(yearMap[overYear], nm)
				}
			}
			res = append(res, yearMap)
			overYear--
		}

	} else {
		yearMap := make(map[int][]int, 1)
		muchMonth := overMonth - startMonth
		if muchMonth > 0 {
			for m := 0; m < muchMonth; m++ {
				yearMap[overYear] = append(yearMap[overYear], overMonth)
				overMonth--
			}
		}
		res = append(res, yearMap)
	}

	//数据遍历demo
	//for _,v := range res {
	//	for year,month := range v {
	//		fmt.Println(year , month)
	//	}
	//}
	return
}
