package intx

import "testing"

func TestMin(t *testing.T) {
	tests := []struct {
		x      int
		y      int
		result int
	}{
		{3, 5, 3},
		{10, 4, 4},
		{-1, 1, -1},
	}

	for _, test := range tests {
		got := Min(test.x, test.y)
		if got != test.result {
			t.Errorf("Min(%d, %d) = %d; want %d", test.x, test.y, got, test.result)
		}
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		x      int
		y      int
		result int
	}{
		{3, 5, 5},
		{10, 4, 10},
		{-1, 1, 1},
	}

	for _, test := range tests {
		got := Max(test.x, test.y)
		if got != test.result {
			t.Errorf("Max(%d, %d) = %d; want %d", test.x, test.y, got, test.result)
		}
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		input  int
		result int
	}{
		{3, 3},
		{-5, 5},
		{0, 0},
	}

	for _, test := range tests {
		got := Abs(test.input)
		if got != test.result {
			t.Errorf("Abs(%d) = %d; want %d", test.input, got, test.result)
		}
	}
}
