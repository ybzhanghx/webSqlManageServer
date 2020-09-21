package utils

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type TimeStamp int64

var (
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

func (nt *NullTimeStamp) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}

	return time.Unix(nt.TimeStamp, 0), nil

}
