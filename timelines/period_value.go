package core

import (
	"sort"
	"time"
)

// PeriodValue associate a value to given period
type PeriodValue[T any] struct {
	Period Period
	Value  T
}

// NewPeriodValue create new PeriodValue
func NewPeriodValue[T any](period Period, value T) PeriodValue[T] {
	return PeriodValue[T]{
		Period: period,
		Value:  value,
	}
}

// NewPeriodValueFromTimes create new PeriodValue from times
func NewPeriodValueFromTimes[T any](start time.Time, end time.Time, value T) (*PeriodValue[T], error) {
	period, err := NewPeriod(start, end)
	if err != nil {
		return nil, err
	}

	pv := PeriodValue[T]{
		Period: *period,
		Value:  value,
	}
	return &pv, nil
}

// IsEmpty checks if Period is empty
func (p *PeriodValue[T]) IsEmpty() bool {
	return p.Period.IsEmpty()
}

// Clamp limits PeriodValue within specified Period
func (p *PeriodValue[T]) Clamp(limit Period) (PeriodValue[T], error) {
	clamp, err := p.Period.Clamp(limit)
	if err != nil {
		return PeriodValue[T]{}, err
	}
	return PeriodValue[T]{Period: clamp, Value: p.Value}, nil
}

// SplitAllPeriods get all periods of PeriodValue list
func SplitAllPeriods[T any](periodValues []PeriodValue[T]) []Period {
	var result []Period
	timeMap := make(map[time.Time]struct{}, len(periodValues)*2)

	for _, pv := range periodValues {
		timeMap[pv.Period.Start] = struct{}{}
		timeMap[pv.Period.End] = struct{}{}
	}

	// Extract keys (times) and ensure order
	allTimes := make([]time.Time, 0, len(timeMap))
	for t := range timeMap {
		allTimes = append(allTimes, t)
	}
	sort.Slice(allTimes, func(i, j int) bool {
		return allTimes[i].Before(allTimes[j])
	})

	for i := 1; i < len(allTimes); i++ {
		result = append(result, Period{
			Start: allTimes[i-1],
			End:   allTimes[i],
		})
	}

	return result
}

// ClampPeriods clamp each PeriodValue with limit
func ClampPeriods[T any](periodValues []PeriodValue[T], limit Period) []PeriodValue[T] {
	var results []PeriodValue[T]

	for _, pv := range periodValues {
		clamp, err := pv.Clamp(limit)
		if err == nil && !clamp.Period.IsEmpty() {
			results = append(results, clamp)
		}
	}

	return results
}
