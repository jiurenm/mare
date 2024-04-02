package randx

import (
	"testing"
)

func TestRandString(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{name: "happy path", input: 10},
		{name: "edge case: n = 0", input: 0},
		{name: "edge case: n = 1", input: 1},
		{name: "edge case: n = 1000", input: 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandString(tt.input)
			if len(got) != tt.input {
				t.Errorf("RandString(%d) = %q; want length %d", tt.input, got, tt.input)
			}
		})
	}
}

func TestRandNumber(t *testing.T) {
	tests := []struct {
		name  string
		input int
	}{
		{name: "happy path", input: 10},
		{name: "edge case: n = 0", input: 0},
		{name: "edge case: n = 1", input: 1},
		{name: "edge case: n = 1000", input: 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandNumber(tt.input)
			if len(got) != tt.input {
				t.Errorf("RandNumber(%d) = %q; want length %d", tt.input, got, tt.input)
			}
		})
	}
}
