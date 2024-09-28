package core

import (
	"errors"
	"time"
)

// TimeLineBuilder permet de construire une Timeline de manière fluide.
type TimeLineBuilder[T comparable] struct {
	items []PeriodValue[T]
	err   error
}

// NewTimeLineBuilder crée une nouvelle instance de TimeLineBuilder.
func NewTimeLineBuilder[T comparable]() *TimeLineBuilder[T] {
	return &TimeLineBuilder[T]{items: []PeriodValue[T]{}}
}

// AddPeriod ajoute une période avec une valeur à la timeline.
func (b *TimeLineBuilder[T]) AddPeriod(start, end time.Time, value T) *TimeLineBuilder[T] {
	if b.err != nil {
		return b
	}

	p, err := NewPeriod(start, end)
	if err != nil {
		b.err = err
		return b
	}
	pv := NewPeriodValue(*p, value)

	b.items = append(b.items, pv)
	return b
}

// AddPeriodValue ajoute un PeriodValue directement à la timeline.
func (b *TimeLineBuilder[T]) AddPeriodValue(pv PeriodValue[T]) *TimeLineBuilder[T] {
	if b.err != nil {
		return b
	}

	if !pv.Period.Start.Before(pv.Period.End) {
		b.err = errors.New("end date must be after start date")
		return b
	}

	b.items = append(b.items, pv)
	return b
}

// AddMonth adds a period corresponding to a given month with a value.
func (b *TimeLineBuilder[T]) AddMonth(year int, month int, value T) *TimeLineBuilder[T] {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	return b.AddPeriod(start, end, value)
}

// AddDay adds a period corresponding to a given day with a value.
func (b *TimeLineBuilder[T]) AddDay(year int, month int, day int, value T) *TimeLineBuilder[T] {
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 1)

	return b.AddPeriod(start, end, value)
}

// Build builds the Timeline by sorting the periods in chronological order.
func (b *TimeLineBuilder[T]) Build() (Timeline[T], error) {
	if b.err != nil {
		return Timeline[T]{}, b.err
	}

	t := Timeline[T]{Items: b.items}
	t.SortTimelineByPeriodStart()

	return t, nil
}
