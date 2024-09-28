package core

import (
	"errors"
	"sort"
)

// Timeline represents a list of PeriodValue objects
type Timeline[T any] struct {
	Items []PeriodValue[T]
}

// NewTimeline creates and returns an empty Timeline.
func NewTimeline[T any]() Timeline[T] {
	return Timeline[T]{
		Items: []PeriodValue[T]{}, // Initialize an empty slice
	}
}

// SortTimelineByPeriodStart sorts the Timeline items by the Start date of their Periods
func (t *Timeline[T]) SortTimelineByPeriodStart() {
	sort.Slice(t.Items, func(i, j int) bool {
		return t.Items[i].Period.Start.Before(t.Items[j].Period.Start)
	})
}

func (t *Timeline[T]) FindIntersects(period Period) []PeriodValue[T] {
	var items []PeriodValue[T]

	for _, current := range t.Items {

		// items are ordered, so if we are after period then we finished scan
		if current.Period.Start.After(period.End) {
			break
		}

		if current.Period.Intersects(period) {
			items = append(items, current)
		}
	}

	return items
}

// Add allows adding a new PeriodValue to the Timeline
func (t *Timeline[T]) Add(newPeriod Period, newValue T) {
	// Update the Timeline items with the new list
	t.Items = append(t.Items, PeriodValue[T]{
		Period: newPeriod,
		Value:  newValue,
	})

	// Ensure the items are sorted by the Start date after adding
	t.SortTimelineByPeriodStart()
}

// GetAll returns all PeriodValues in the Timeline
func (t *Timeline[T]) GetAll() []PeriodValue[T] {
	return t.Items
}

func computeValuesOnSamePeriods[T any](buffer []PeriodValue[T], f func(p Period, a T, b T) T) []PeriodValue[T] {
	var items []PeriodValue[T]
	periods := SplitAllPeriods(buffer)

	for _, period := range periods {
		var currentValue T

		for _, candidate := range buffer {
			if candidate.Period.Intersects(period) {
				currentValue = f(period, candidate.Value, currentValue)
			}
		}

		items = append(items, NewPeriodValue(period, currentValue))
	}

	return items
}

// ResolveConflicts returns another Timeline having all values with same period aggregated, slicing them if necessary.
func (t *Timeline[T]) ResolveConflicts(f func(p Period, a T, b T) T) (Timeline[T], error) {
	var items []PeriodValue[T]
	var buffer []PeriodValue[T]
	var currentPeriod Period

	for i, next := range t.Items {
		if i == 0 {
			currentPeriod = next.Period
			buffer = append(buffer, next)
			continue
		}

		// We assume that periods are chronologically sorted
		if next.Period.Before(currentPeriod) {
			return Timeline[T]{}, errors.New("timeline should have sorted periods")
		}

		if next.Period.After(currentPeriod) {
			computed := computeValuesOnSamePeriods(buffer, f)
			items = append(items, computed...)
			currentPeriod = next.Period

			buffer = ClampPeriods(buffer, currentPeriod)
			buffer = append(buffer, next)

			continue
		}

		period, err := NewPeriod(currentPeriod.Start, maxTime(next.Period.Start, currentPeriod.End))
		if err != nil {
			return Timeline[T]{}, err
		}
		currentPeriod = *period
		buffer = append(buffer, next)
	}

	computed := computeValuesOnSamePeriods(buffer, f)
	items = append(items, computed...)
	buffer = make([]PeriodValue[T], 0)

	return Timeline[T]{Items: items}, nil
}

// Optimize merges all contiguous periods having same value
func (t *Timeline[T]) Optimize(equalityComparer func(a T, b T) bool) Timeline[T] {
	var previous PeriodValue[T]
	var items []PeriodValue[T]

	for i, current := range t.Items {
		if i == 0 {
			previous = current
			continue
		}

		if current.Period.IsContiguous(previous.Period) && equalityComparer(previous.Value, current.Value) {
			previous = PeriodValue[T]{Value: previous.Value, Period: Period{Start: previous.Period.Start, End: current.Period.End}}
			continue
		}

		items = append(items, previous)
		previous = current
	}

	if !previous.IsEmpty() {
		items = append(items, previous)
	}

	return Timeline[T]{Items: items}
}

// Aggregate two timelines and return another timeline.
func (t *Timeline[T]) Aggregate(other *Timeline[T], f func(period Period, a T, b T) T) (Timeline[T], error) {
	c1 := len(t.Items)
	c2 := len(other.Items)

	if c1 == 0 {
		return *other, nil
	}
	if c2 == 0 {
		return *t, nil
	}

	concat := Timeline[T]{
		Items: append(t.Items, other.Items...),
	}
	concat.SortTimelineByPeriodStart()
	result, err := concat.ResolveConflicts(f)
	if err != nil {
		return Timeline[T]{}, err
	}

	return result, nil
}
