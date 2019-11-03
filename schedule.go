package goschedule

import (
	"errors"
	"fmt"
	"time"
)

type SCycle uint8

const (
	Monthly SCycle = 1
	Weekly  SCycle = 1 << 1
	Daily   SCycle = 1 << 2
)

type Recurrence struct {
	Cycle     SCycle
	WeekIndex int
	DayIndex  int
}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

type Schedule struct {
	Period TimeRange
	Rule   Recurrence
}

func InitSchedule(start, end time.Time, c SCycle, w, d int) (res Schedule, err error) {
	if end.Before(start.AddDate(0, 0, 1)) {
		err = errors.New("end date should after start")
		return
	}
	if d < 1 || w < 0 {
		err = errors.New("weekIndex and dayIndex invalid")
	}
	if c == Weekly && d > 6 {
		err = errors.New("weekly schedule should have day index <=6")
	}
	if c == Monthly {
		if d > 31 {
			err = errors.New("monthly schedule should have day index <=31")
		}
	}
	if err != nil {
		return
	}
	res = Schedule{
		Period: TimeRange{
			Start: time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location()),
			End:   time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location()),
		},
		Rule: Recurrence{
			Cycle:     c,
			WeekIndex: w,
			DayIndex:  d,
		},
	}
	return
}

type CheckerInfc interface {
	IsScheduledDay(t time.Time, s Schedule) (res bool)
	LocScheduledDays(s Schedule) []time.Time
}

type SChecker struct{}

func GetChecker() (c SChecker) {
	c = SChecker{}
	return
}

func (c *SChecker) IsScheduledDay(t time.Time, s Schedule) (res bool) {
	res = false
	if t.Before(s.Period.Start.AddDate(0, 0, 1)) || t.After(s.Period.End) {
		return
	}
	switch s.Rule.Cycle {
	case Monthly:
		res = monthlyCheck(t, s.Rule.WeekIndex, s.Rule.DayIndex)
		break
	case Weekly:
		res = weeklyCheck(t, s.Rule.DayIndex)
		break
	case Daily:
		duration := int(t.Sub(s.Period.Start).Hours()) / 24
		res = (duration%s.Rule.DayIndex == 0)
		break
	default:
		res = false
	}
	return
}

func monthlyCheck(t time.Time, weekIndex, dayIndex int) (res bool) {
	if weekIndex < 1 { // If WeekIndex < 1, schedule use the dayIndex-th day of the month
		res = time.Now().In(t.Location()).Day() == dayIndex
		return
	}
	res = (weekIndex == t.Day()/7 && dayIndex == int(t.Weekday()))
	return
}

func weeklyCheck(t time.Time, dayIndex int) (res bool) {
	res = (int(t.Weekday()) == dayIndex)
	return
}

func (c *SChecker) LocScheduledDays(s Schedule) (res []time.Time, err error) {
	res = make([]time.Time, 0)
	switch s.Rule.Cycle {
	case Daily:
		res = dailyAppend(s)
		break
	case Weekly:
		res = weeklyAppend(s)
		break
	case Monthly:
		res = monthlyAppend(s)
		break
	default:
		err = fmt.Errorf("undefined recurrence cycle %d", s.Rule.Cycle)
	}
	return
}

func dailyAppend(s Schedule) (res []time.Time) {
	res = make([]time.Time, 0)
	d := s.Period.Start.AddDate(0, 0, s.Rule.DayIndex)
	for {
		if !d.Before(s.Period.End) {
			break
		}
		res = append(res, d)
		d = d.AddDate(0, 0, s.Rule.DayIndex)
	}
	return
}

func weeklyAppend(s Schedule) (res []time.Time) {
	d := s.Period.Start.AddDate(0, 0,
		(7+s.Rule.DayIndex-int(s.Period.Start.Weekday()))%7)
	for {
		if !d.Before(s.Period.End) {
			break
		}
		res = append(res, d)
		d = d.AddDate(0, 0, 7)
	}
	return
}

func monthlyAppend(s Schedule) (res []time.Time) {
	res = make([]time.Time, 0)
	m := time.Date(s.Period.Start.Year(), s.Period.Start.Month(), 1, 0, 0, 0, 0, s.Period.Start.Location())
	var d time.Time
	if s.Rule.WeekIndex > 0 {
		d = m.AddDate(0, 0,
			7*(s.Rule.WeekIndex-1)+(7+s.Rule.DayIndex-int(s.Period.Start.Weekday()))%7)
	} else {
		d = m.AddDate(0, 0, s.Rule.DayIndex)
	}
	if d.Before(m.AddDate(0, 1, 0)) && d.After(s.Period.Start) && d.Before(s.Period.End) {
		res = append(res, d)
	}
	for {
		m = m.AddDate(0, 1, 0)
		if s.Rule.WeekIndex > 0 {
			d = m.AddDate(0, 0,
				7*(s.Rule.WeekIndex-1)+(7+s.Rule.DayIndex-int(m.Weekday()))%7)
		} else {
			d = m.AddDate(0, 0, s.Rule.DayIndex)
		}
		if !d.Before(s.Period.End) {
			break
		}
		if d.Before(m.AddDate(0, 1, 0)) {
			res = append(res, d)
		}
	}
	return
}
