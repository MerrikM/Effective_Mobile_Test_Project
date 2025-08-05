package model

import (
	"strings"
	"time"
)

type SubscriptionDetails struct {
	ID          int          `db:"id" json:"id"`
	ServiceName string       `db:"service_name" json:"service_name"`
	Price       int          `db:"price" json:"price"`
	UserID      string       `db:"user_id" json:"user_id"`
	StartDate   DayMonthYear `db:"start_date" json:"start_date"`
	EndDate     DayMonthYear `db:"end_date" json:"end_date"`
}

type DayMonthYear time.Time

func (date *DayMonthYear) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	t, err := time.Parse("02-01-2006", str)
	if err != nil {
		return err
	}
	*date = DayMonthYear(t)
	return nil
}

func (date DayMonthYear) MarshalJSON() ([]byte, error) {
	t := time.Time(date)
	return []byte(`"` + t.Format("02-01-2006") + `"`), nil
}

func (date DayMonthYear) ToTime() time.Time {
	return time.Time(date)
}
