package intx

import "testing"

func TestPagination(t *testing.T) {
	t.Run("Happy path", func(t *testing.T) {
		start, end := Pagination(2, 10, 100)
		if start != 10 || end != 20 {
			t.Errorf("Expected start: 10 and end: 20, but got start: %d and end: %d", start, end)
		}
	})

	t.Run("Edge case: page = 0", func(t *testing.T) {
		start, end := Pagination(0, 10, 100)
		if start != 0 || end != 0 {
			t.Errorf("Expected start: 0 and end: 0, but got start: %d and end: %d", start, end)
		}
	})

	t.Run("Edge case: pageSize = 0", func(t *testing.T) {
		start, end := Pagination(1, 0, 100)
		if start != 0 || end != 0 {
			t.Errorf("Expected start: 0 and end: 0, but got start: %d and end: %d", start, end)
		}
	})

	t.Run("Edge case: page = -1", func(t *testing.T) {
		start, end := Pagination(-1, 10, 100)
		if start != 0 || end != 00 {
			t.Errorf("Expected start: 0 and end: 0, but got start: %d and end: %d", start, end)
		}
	})

	t.Run("Edge case: page = 1 and pageSize > length", func(t *testing.T) {
		start, end := Pagination(1, 200, 100)
		if start != 0 || end != 100 {
			t.Errorf("Expected start: 0 and end: 100, but got start: %d and end: %d", start, end)
		}
	})

	t.Run("Edge case: page > totalPage", func(t *testing.T) {
		start, end := Pagination(5, 10, 35)
		if start != 0 || end != 0 {
			t.Errorf("Expected start: 0 and end: 0, but got start: %d and end: %d", start, end)
		}
	})
}
