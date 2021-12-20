package table

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Time 自定义时间
type Time time.Time

const (
	format = "2006-01-02 15:04:05"
	zone   = "Asia/Shanghai"
)

// UnmarshalJSON implements json unmarshal interface.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+format+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

// MarshalJSON implements json marshal interface.
func (t Time) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte{'n', 'u', 'l', 'l'}, nil
	}

	b := make([]byte, 0, len(format)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, format)
	b = append(b, '"')
	return b, nil
}

// String ...
func (t Time) String() string {
	return time.Time(t).Format(format)
}

// local ...
func (t Time) local() time.Time {
	loc, _ := time.LoadLocation(zone)
	return time.Time(t).In(loc)
}

// Value ...
func (t Time) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

// Scan value of time.Time
func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
