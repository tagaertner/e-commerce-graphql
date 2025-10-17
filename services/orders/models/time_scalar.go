package models

import (
	"fmt"
	"time"
	"io"
	"database/sql/driver"
)
// Todo remove all file when change to simple time like job story
type Time time.Time

func (t Time) MarshalGQL(w io.Writer) {
	str := time.Time(t).Format(time.RFC3339)
	_, _ = w.Write([]byte(fmt.Sprintf("%q", str)))
}

func (t *Time) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("time should be a string")
	}
	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	*t = Time(parsed)
	return nil
}

func (t Time)Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *Time) Scan(value interface{}) error{
	if val, ok := value.(time.Time); ok {
		*t = Time(val)
		return nil
	}
	return fmt.Errorf("cannot scan value into TIme: %v", value)
}

// Now returns the current time as a models.Time

func Now() Time {
    return Time(time.Now().UTC())
}