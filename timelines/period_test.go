package core

import (
	"testing"
	"time"
)

func TestNewPeriod(t *testing.T) {
	start := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 8, 31, 23, 59, 59, 0, time.UTC)

	period, err := NewPeriod(start, end)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if period.Start != start || period.End != end {
		t.Errorf("expected start %v and end %v, got start %v and end %v", start, end, period.Start, period.End)
	}
}

func TestNewPeriod_Invalid(t *testing.T) {
	start := time.Date(2024, 8, 31, 23, 59, 59, 0, time.UTC)
	end := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)

	_, err := NewPeriod(start, end)
	if err == nil {
		t.Error("expected an error when end date is before start date")
	}
}

func TestPeriod_Contains(t *testing.T) {
	start := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 8, 31, 23, 59, 59, 0, time.UTC)
	period, _ := NewPeriod(start, end)

	testDate := time.Date(2024, 8, 15, 0, 0, 0, 0, time.UTC)
	if !period.Contains(testDate) {
		t.Errorf("expected date %v to be contained in the period", testDate)
	}
}

func TestPeriod_Duration(t *testing.T) {
	start := time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 8, 31, 23, 59, 59, 0, time.UTC)
	period, _ := NewPeriod(start, end)

	expectedDuration := end.Sub(start)
	if period.Duration() != expectedDuration {
		t.Errorf("expected duration %v, got %v", expectedDuration, period.Duration())
	}
}

func TestMonth(t *testing.T) {
	period, err := Month(2024, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	if !period.Start.Equal(expectedStart) {
		t.Errorf("Expected start %v, got %v", expectedStart, period.Start)
	}

	if !period.End.Equal(expectedEnd) {
		t.Errorf("Expected end %v, got %v", expectedEnd, period.End)
	}
}

func TestYear(t *testing.T) {
	period, err := Year(2024)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	if !period.Start.Equal(expectedStart) {
		t.Errorf("Expected start %v, got %v", expectedStart, period.Start)
	}

	if !period.End.Equal(expectedEnd) {
		t.Errorf("Expected end %v, got %v", expectedEnd, period.End)
	}
}

func TestPeriod_SplitJanuaryByDays(t *testing.T) {
	period, err := Month(2024, 1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	days := period.SplitByDays()
	count := 0
	current := *period

	for d := range days {
		count++
		current = d
	}

	if !current.Start.Equal(time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("expected start %v, got %v", time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC), current.Start)
	}

	if count != 31 {
		t.Errorf("expected 31 days, got %v", count)
	}
}

func TestPeriod_SplitFebruaryByDays(t *testing.T) {
	period, err := Month(2024, 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	days := period.SplitByDays()
	count := 0
	current := *period

	for d := range days {
		count++
		current = d
	}

	if !current.Start.Equal(time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("expected start %v, got %v", time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC), current.Start)
	}

	if count != 29 {
		t.Errorf("expected 29 days, got %v", count)
	}
}

func TestPeriod_SplitYear2024ByDays(t *testing.T) {
	period, err := Year(2024)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	days := period.SplitByDays()
	count := 0
	current := *period

	for d := range days {
		count++
		current = d
	}

	if !current.Start.Equal(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("expected start %v, got %v", time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), current.Start)
	}

	if count != 366 {
		t.Errorf("expected 31 days, got %v", count)
	}
}

func TestPeriod_SplitYear2024ByMonths(t *testing.T) {
	period, err := Year(2024)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	days := period.SplitByMonths()
	count := 0
	var current Period

	for d := range days {
		count++
		current = d
	}

	if !current.Start.Equal(time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("expected start %v, got %v", time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC), current.Start)
	}

	if !current.End.Equal(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("expected end %v, got %v", time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), current.End)
	}

	if count != 12 {
		t.Errorf("expected 12 months, got %v", count)
	}
}

func TestSplitFromPeriod_NoIntersection(t *testing.T) {
	p, _ := Month(2024, 1)

	nonIntersectingPeriod, _ := Month(2024, 3)

	splitChan := p.SplitFromPeriod(*nonIntersectingPeriod)

	_, ok := <-splitChan
	if ok {
		t.Errorf("Expected no periods, but got a value")
	}
}

func TestSplitFromPeriod_ShouldReturn3Periods(t *testing.T) {
	feb2024, _ := Month(2024, 2)
	intersectingPeriod, _ := Day(2024, 2, 15)

	splitChan := feb2024.SplitFromPeriod(*intersectingPeriod)

	var results []Period
	for p := range splitChan {
		results = append(results, p)
	}

	if len(results) != 3 {
		t.Errorf("Expected 3 periods, got %v", len(results))
	}

	expectedPeriod1, _ := NewPeriod(
		time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
	)
	expectedPeriod2, _ := NewPeriod(
		time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 16, 0, 0, 0, 0, time.UTC),
	)
	expectedPeriod3, _ := NewPeriod(
		time.Date(2024, 2, 16, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
	)

	if !results[0].Equal(*expectedPeriod1) {
		t.Errorf("Expected period1 to be %v, got %v", expectedPeriod1, results[0])
	}

	if !results[1].Equal(*expectedPeriod2) {
		t.Errorf("Expected period2 to be %v, got %v", expectedPeriod2, results[1])
	}

	if !results[2].Equal(*expectedPeriod3) {
		t.Errorf("Expected period3 to be %v, got %v", expectedPeriod3, results[2])
	}
}

func TestPeriod_Clamp(t *testing.T) {
	tests := []struct {
		name        string
		p           Period
		limit       Period
		expected    Period
		expectError bool
	}{
		{
			name: "Period entirely inside limit",
			p: Period{
				Start: DateOnly(2024, 1, 16),
				End:   DateOnly(2024, 1, 18),
			},
			limit: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 20),
			},
			expected: Period{
				Start: DateOnly(2024, 1, 16),
				End:   DateOnly(2024, 1, 18),
			},
			expectError: false,
		},
		{
			name: "Period starting before limit",
			p: Period{
				Start: DateOnly(2024, 1, 10),
				End:   DateOnly(2024, 1, 18),
			},
			limit: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 20),
			},
			expected: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 18),
			},
			expectError: false,
		},
		{
			name: "Periode endind after limit",
			p: Period{
				Start: DateOnly(2024, 1, 17),
				End:   DateOnly(2024, 1, 25),
			},
			limit: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 20),
			},
			expected: Period{
				Start: DateOnly(2024, 1, 17),
				End:   DateOnly(2024, 1, 20),
			},
			expectError: false,
		},
		{
			name: "Period totally before limit",
			p: Period{
				Start: DateOnly(2024, 1, 1),
				End:   DateOnly(2024, 1, 10),
			},
			limit: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 20),
			},
			expected:    Period{},
			expectError: true,
		},
		{
			name: "Period totally after limit",
			p: Period{
				Start: DateOnly(2024, 1, 21),
				End:   DateOnly(2024, 1, 25),
			},
			limit: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 20),
			},
			expected:    Period{},
			expectError: true,
		},
		{
			name: "Period equals limit",
			p: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 20),
			},
			limit: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 20),
			},
			expected: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 20),
			},
			expectError: false,
		},
		{
			name: "Period greater than limit",
			p: Period{
				Start: DateOnly(2024, 1, 10),
				End:   DateOnly(2024, 1, 25),
			},
			limit: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 20),
			},
			expected: Period{
				Start: DateOnly(2024, 1, 15),
				End:   DateOnly(2024, 1, 20),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &tt.p
			clamped, err := p.Clamp(tt.limit)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !clamped.Start.Equal(tt.expected.Start) || !clamped.End.Equal(tt.expected.End) {
				t.Errorf("Clamped period mismatch.\nExpected: %v - %v\nGot: %v - %v",
					tt.expected.Start.Format("2006-01-02"),
					tt.expected.End.Format("2006-01-02"),
					clamped.Start.Format("2006-01-02"),
					clamped.End.Format("2006-01-02"))
			}
		})
	}
}
