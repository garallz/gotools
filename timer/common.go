package timer

import "time"

// parse time format
const (
	TimeFormatString = "2006-01-02 15:04:05"
	TimeFormatLength = 19
)

// GetNowStamp : get time now stamp : millitime
func GetNowStamp() int64 {
	return time.Now().UnixNano() / MilliTimeUnit
}

// GetTimeStamp : get time now stamp : millitime
func GetTimeStamp(t time.Time) int64 {
	return t.UnixNano() / MilliTimeUnit
}

// GetNextStamp :add time duration stamp : millitime
func GetNextStamp(d time.Duration) int64 {
	return time.Now().Add(d).UnixNano() / MilliTimeUnit
}

// ConvDurationStamp : convert time duration to uint64 : millitime
func ConvDurationStamp(d time.Duration) int64 {
	return int64(d) / MilliTimeUnit
}

// AddDurationStamp : add time duration to uint64 : millitime
func AddDurationStamp(t time.Time, d time.Duration) int64 {
	return t.Add(d).UnixNano() / MilliTimeUnit
}

// time unit
const (
	NanoTimeUnit   = 1
	MicroTimeUnit  = 1000 * NanoTimeUnit
	MilliTimeUnit  = 1000 * MicroTimeUnit
	SecondTimeUnit = 1000 * MilliTimeUnit
	MinuteTimeUnit = 60 * SecondTimeUnit
	HourTimeUnit   = 60 * MinuteTimeUnit
	DayTimeUnit    = 24 * HourTimeUnit
)

// ByNext sort by next time
type ByNext []TimerStruct

func (t ByNext) Len() int {
	return len(t)
}

func (t ByNext) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t ByNext) Less(i, j int) bool {
	return t[i].Next() < t[j].Next()
}

// TimeSplit : split next time run function
func TimeSplit(rows []TimerStruct, timestamp int64) ([]TimerStruct, []TimerStruct) {
	var mins, maxs = make([]TimerStruct, 0), make([]TimerStruct, 0)
	for _, row := range rows {
		if row.Next() <= timestamp {
			mins = append(mins, row)
		} else {
			maxs = append(maxs, row)
		}
	}
	return mins, maxs
}
