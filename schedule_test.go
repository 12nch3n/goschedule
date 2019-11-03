package goschedule

import (
	"testing"
	"time"
)

func TestDaily(t *testing.T) {
	t.Log("test init function")
	c := GetChecker()
	s, _ := time.Parse("2006-01-02 15:04:05 MST", "2019-11-03 00:00:00 UTC")
	t.Logf("res: %v", s)
	e, _ := time.Parse("2006-01-02 15:04:05 MST", "2019-11-09 00:00:00 UTC")
	check, _ := time.Parse("2006-01-02 15:04:05 MST", "2019-11-05 00:00:00 UTC")
	tSchedule, _ := InitSchedule(s, e, Daily, 0, 2)
	t.Logf("res: %v", tSchedule)
	res := c.IsScheduledDay(check, tSchedule)
	dates, _ := c.LocScheduledDays(tSchedule)
	t.Logf("res: %v\ndates:%v", res, dates)
}

func TestWeekly(t *testing.T) {
	t.Log("test init function")
	c := GetChecker()
	s, _ := time.Parse("2006-01-02 15:04:05 MST", "2019-11-01 00:00:00 UTC")
	t.Logf("res: %v", s)
	e, _ := time.Parse("2006-01-02 15:04:05 MST", "2019-11-20 00:00:00 UTC")
	check, _ := time.Parse("2006-01-02 15:04:05 MST", "2019-11-04 00:00:00 UTC")
	tSchedule, _ := InitSchedule(s, e, Weekly, 1, 1)
	t.Logf("res: %v", tSchedule)
	res := c.IsScheduledDay(check, tSchedule)
	dates, _ := c.LocScheduledDays(tSchedule)
	t.Logf("res: %v\ndates:%v", res, dates)
}

func TestMonthly(t *testing.T) {
	t.Log("test init function")
	c := GetChecker()
	s, _ := time.Parse("2006-01-02 15:04:05 MST", "2019-10-31 00:00:00 UTC")
	t.Logf("res: %v", s)
	e, _ := time.Parse("2006-01-02 15:04:05 MST", "2019-11-20 00:00:00 UTC")
	check, _ := time.Parse("2006-01-02 15:04:05 MST", "2019-11-02 00:00:00 UTC")
	tSchedule, _ := InitSchedule(s, e, Monthly, 1, 4)
	t.Logf("res: %v", tSchedule)
	res := c.IsScheduledDay(check, tSchedule)
	dates, _ := c.LocScheduledDays(tSchedule)
	t.Logf("res: %v\ndates:%v", res, dates)
}
