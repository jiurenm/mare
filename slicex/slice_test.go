package slicex

import (
	"reflect"
	"testing"
)

func TestMergeSlices(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice2 := []int{4, 5, 6}
	expectedResult := []int{1, 2, 3, 4, 5, 6}

	result := MergeSlices(slice1, slice2)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v, but got %v", expectedResult, result)
	}
}

func TestPadSlice(t *testing.T) {
	slice := []int{1, 2, 3}
	expectedLength := 5
	expectedResult := []int{1, 2, 3, 0, 0}
	result := PadSlice(slice, expectedLength)

	if len(result) != expectedLength || !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Expected %v with length %v, but got %v with length %v", expectedResult, expectedLength, result, len(result))
	}

	defaultValue := 10
	expectedResultWithDefault := []int{1, 2, 3, 10, 10}
	resultWithDefault := PadSlice(slice, expectedLength, defaultValue)

	if len(resultWithDefault) != expectedLength || !reflect.DeepEqual(resultWithDefault, expectedResultWithDefault) {
		t.Errorf("Expected %v with length %v, but got %v with length %v", expectedResultWithDefault, expectedLength, resultWithDefault, len(resultWithDefault))
	}
}
