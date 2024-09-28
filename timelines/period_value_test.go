package core

import (
	"testing"
)

func TestNewPeriodValue(t *testing.T) {
	september1, err := Day(2024, 9, 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	periodValue := NewPeriodValue(*september1, 45)
	if periodValue.Value != 45 {
		t.Errorf("period value: got %v, want %v", periodValue.Value, 45)
	}
}
