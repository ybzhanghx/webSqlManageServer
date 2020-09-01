package utils

import (
	"fmt"
	"time"
)

type TimeStamp int64

var (
	//dateTimeFormats = []string{RFC3339Micro, RFC3339Millis, time.RFC3339, time.RFC3339Nano, ISO8601LocalTime, ISO8601TimeWithReducedPrecision, ISO8601TimeWithReducedPrecisionLocaltime}
	//rxDateTime      = regexp.MustCompile(DateTimePattern)
	//// MarshalFormat sets the time resolution format used for marshaling time (set to milliseconds)
	//MarshalFormat = RFC3339Millis
	//loc, _ := time.LoadLocation("Local")
	TimeFormat1 = "2006-01-02 15:04:05"
	LocZone, _  = time.LoadLocation("Asia/Shanghai")
)

type NullTimeStamp struct {
	TimeStamp int64

	Valid bool // Valid is true if Time is not NULL
}

func (nt *NullTimeStamp) Scan(value interface{}) (err error) {

	if value == nil {
		nt.TimeStamp, nt.Valid = 0, false
		return
	}

	switch v := value.(type) {
	case time.Time:
		nt.TimeStamp, nt.Valid = v.Unix(), true
		return
	case []byte:
		var tmpTime time.Time
		tmpTime, err = time.ParseInLocation(TimeFormat1, string(v), LocZone)
		nt.Valid = err == nil
		if err != nil {
			nt.TimeStamp = 0

		} else {
			nt.TimeStamp = tmpTime.Unix()
		}
		return
	case string:
		var tmpTime time.Time
		tmpTime, err = time.ParseInLocation(TimeFormat1, v, LocZone)
		nt.Valid = err == nil
		if err != nil {
			nt.TimeStamp = 0
		} else {
			nt.TimeStamp = tmpTime.Unix()
		}
		return
	}

	nt.Valid = false
	return fmt.Errorf("Can't convert %T to time.Time", value)
}
