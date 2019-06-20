package common

import (
	"fmt"
	"testing"
	"time"
)

func TestTheOtherDayTimeRange(t *testing.T) {
	//
	_, _, beginStr, endStr := TheOtherDayTimeRange(0)
	t.Logf("%s, %s\n", beginStr, endStr)
	_, _, beginStr, endStr = TheOtherDayTimeRange(1)
	t.Logf("%s, %s\n", beginStr, endStr)
	_, _, beginStr, endStr = TheOtherDayTimeRange(-1)
	t.Logf("%s, %s\n", beginStr, endStr)
}

func TestTimeStrToUint32Plus(t *testing.T) {
	s := "2019/04/12"
	timestamp, err := TimeStrToUint32Plus(s)
	if err != nil {
		t.Fatalf("failed : %v", err)
	}
	if timestamp != 1554998400 {
		t.Fatalf("%s should be 1554998400", s)
	}

	s = "2019/4/12/0"
	_, err = TimeStrToUint32Plus(s)
	if err == nil {
		t.Fatalf("should return err.")
	}
}

func TestTimeFormat(t *testing.T) {
	curZeroTime := GetZeroTime(time.Now().Unix())
	strTime := GetTimeFormat(curZeroTime)
	fmt.Println(curZeroTime, strTime)
	begin, end := GetBeginAndEndString(time.Now().Unix())
	fmt.Println(begin, end)
}
