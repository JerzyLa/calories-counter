package common

import (
	"encoding/json"
	"strings"
	"time"
)

func String(s string) *string {
	return &s
}

type Date time.Time
type Time time.Time

func (d *Date) String() string {
	if d == nil {
		return ""
	}
	t := time.Time(*d)
	res := strings.Split(t.String(), " ")[0]
	return res
}

func (t *Time) String() string {
	if t == nil {
		return ""
	}
	tt := time.Time(*t)
	res := strings.Split(tt.String(), " ")[1][:8]
	return res
}

func (d *Date) UnmarshalJSON(b []byte) error {
	if len(b) != 12 {
		return &time.ParseError{}
	}
	b = b[1:11]
	t, err := time.Parse("2006-01-02", string(b))
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(time.Time(d))
	if err != nil {
		return nil, err
	}
	s := string(b)
	res := strings.Split(s, "T")[0] + "\""
	return []byte(res), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	if len(b) != 10 {
		return &time.ParseError{}
	}
	b = b[1:9]
	tt, err := time.Parse("15:04:05", string(b))
	if err != nil {
		return err
	}
	*t = Time(tt)
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(time.Time(t))
	if err != nil {
		return nil, err
	}
	s := string(b)
	res := strings.Split(s, "T")[1][:8]
	res = "\"" + res + "\""
	return []byte(res), nil
}
