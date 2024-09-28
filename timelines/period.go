package core

import (
	"errors"
	"time"
)

type Period struct {
	Start time.Time
	End   time.Time
}

func NewPeriod(start, end time.Time) (*Period, error) {
	if !end.After(start) {
		return nil, errors.New("end date must be after start date")
	}
	return &Period{Start: start, End: end}, nil
}

func Empty() Period {
	return Period{Start: time.Time{}, End: time.Time{}}
}

// Day returns a Period for the given year, month and day.
// The start is given day, and the end is the next day (exclusive).
func Day(year int, month int, day int) (*Period, error) {
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	nextDay := start.AddDate(0, 0, 1)

	return NewPeriod(start, nextDay)
}

func DateOnly(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// Month returns a Period for the given year and month.
// The start is the first day of the month, and the end is the first day of the next month (exclusive).
func Month(year int, month int) (*Period, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	nextMonth := start.AddDate(0, 1, 0)

	return NewPeriod(start, nextMonth)
}

// Year returns a Period for the given year.
// The start is the first day of the year, and the end is the first day of the next year (exclusive).
func Year(year int) (*Period, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	nextYear := start.AddDate(1, 0, 0)

	return NewPeriod(start, nextYear)
}

// Equal compares two periods.
func (p Period) Equal(other Period) bool {
	return p.Start.Equal(other.Start) && p.End.Equal(other.End)
}

// Duration returns duration of given period.
func (p *Period) Duration() time.Duration {
	return p.End.Sub(p.Start)
}

// Contains check if a periods contains another one.
func (p *Period) Contains(t time.Time) bool {
	return (t.After(p.Start) || t.Equal(p.Start)) && (t.Before(p.End) || t.Equal(p.End))
}

// ContainsPeriod checks if the current period fully contains another period.
func (p *Period) ContainsPeriod(other Period) bool {
	return (p.Start.Before(other.Start) || p.Start.Equal(other.Start)) &&
		(p.End.After(other.End) || p.End.Equal(other.End))
}

// Intersects checks if two periods overlap.
func (p *Period) Intersects(other Period) bool {
	return p.Start.Before(other.End) && p.End.After(other.Start)
}

// Split a period using given function.
func (p *Period) Split(f func(current time.Time) time.Time) <-chan Period {
	ch := make(chan Period)

	go func() {
		defer close(ch)

		current := p.Start
		for current.Before(p.End) {
			next := f(current)
			ch <- Period{Start: current, End: next}
			current = next
		}
	}()

	return ch
}

// SplitByDays returns periods for each days in given period.
func (p *Period) SplitByDays() <-chan Period {
	return p.Split(func(current time.Time) time.Time { return current.AddDate(0, 0, 1) })
}

// SplitByMonths returns periods for each months in given period.
func (p *Period) SplitByMonths() <-chan Period {
	return p.Split(func(current time.Time) time.Time { return current.AddDate(0, 1, 0) })
}

func (p *Period) Before(other Period) bool {
	return p.End.Before(other.Start) || p.End.Equal(other.Start)
}

func (p *Period) After(other Period) bool {
	return p.Start.After(other.End) || p.Start.Equal(other.End)
}

// Helper function to find the minimum of two times
func minTime(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t1
	}
	return t2
}

// Helper function to find the maximum of two times
func maxTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

// SplitFromPeriod returns a split of periods intersecting with given period
func (p *Period) SplitFromPeriod(period Period) <-chan Period {
	ch := make(chan Period)

	go func() {
		defer close(ch)

		if !p.Intersects(period) {
			return
		}

		// before intersecting part
		if p.Start.Before(period.Start) {
			ch <- Period{Start: p.Start, End: period.Start}
		}

		// intersecting part
		overlapStart := maxTime(p.Start, period.Start)
		overlapEnd := minTime(p.End, period.End)
		ch <- Period{
			Start: overlapStart,
			End:   overlapEnd,
		}

		// after intersecting part
		if p.End.After(period.End) {
			ch <- Period{
				Start: period.End,
				End:   p.End,
			}
		}
	}()

	return ch
}

// IsEmpty checks if period is empty
func (p *Period) IsEmpty() bool {
	return !p.Start.Before(p.End)
}

func (p *Period) Clamp(limit Period) (Period, error) {
	if p.Intersects(limit) {
		start := maxTime(limit.Start, p.Start)
		end := minTime(limit.End, p.End)
		period, err := NewPeriod(start, end)
		if err != nil {
			return Empty(), err
		}
		return *period, nil
	}

	return Empty(), errors.New("limit is outside")
}

// IsContiguous checks if the other Period is contiguous
func (p *Period) IsContiguous(other Period) bool {
	return p.End.Equal(other.Start) || p.Start.Equal(other.End)
}
