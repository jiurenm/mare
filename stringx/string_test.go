package stringx

import (
	"reflect"
	"testing"
)

func TestStringToBytes(t *testing.T) {
	// Happy path
	input := "hello"
	expected := []byte{104, 101, 108, 108, 111}
	result := StringToBytes(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("StringToBytes(%s) = %v; want %v", input, result, expected)
	}

	// Edge case: Empty string
	input = ""
	var expected2 []byte
	result = StringToBytes(input)
	if !reflect.DeepEqual(result, expected2) {
		t.Errorf("StringToBytes(%s) = %v; want %v", input, result, expected2)
	}

	// Edge case: String with special characters
	input = "abc123!@#"
	expected = []byte{97, 98, 99, 49, 50, 51, 33, 64, 35}
	result = StringToBytes(input)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("StringToBytes(%s) = %v; want %v", input, result, expected)
	}
}

func TestBytesToString(t *testing.T) {
	// Happy path
	expected := "hello"
	input := []byte("hello")
	result := BytesToString(input)
	if result != expected {
		t.Errorf("BytesToString(%v) = %v, expected %v", input, result, expected)
	}

	// Edge case: empty byte slice
	expected = ""
	input = []byte{}
	result = BytesToString(input)
	if result != expected {
		t.Errorf("BytesToString(%v) = %v, expected %v", input, result, expected)
	}

	// Edge case: non-empty byte slice
	expected = "test"
	input = []byte{116, 101, 115, 116}
	result = BytesToString(input)
	if result != expected {
		t.Errorf("BytesToString(%v) = %v, expected %v", input, result, expected)
	}
}

func TestCopy(t *testing.T) {
	input := "Hello, World!"
	expected := "Hello, World!"

	result := Copy(input)

	if result != expected {
		t.Errorf("Copy function returned unexpected result: got %s, want %s", result, expected)
	}
}

func TestPadStart(t *testing.T) {
	// Happy path
	t.Run("pad with default value", func(t *testing.T) {
		result := PadStart("test", 8)
		expected := "    test"
		if result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})

	t.Run("pad with custom value", func(t *testing.T) {
		result := PadStart("test", 8, "x")
		expected := "xxxxtest"
		if result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})

	// Edge cases
	t.Run("pad when string is already equal to size", func(t *testing.T) {
		result := PadStart("test", 4)
		expected := "test"
		if result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})

	t.Run("pad with negative size", func(t *testing.T) {
		result := PadStart("test", -4)
		expected := "test"
		if result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})

	t.Run("pad with empty string", func(t *testing.T) {
		result := PadStart("", 5)
		expected := "     "
		if result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})
}

func TestPadEnd(t *testing.T) {
	// Happy path
	happyResult := PadEnd("hello", 8, " world")
	if happyResult != "hello wo" {
		t.Errorf("PadEnd('hello', 8, ' world') = %s; want 'hello wo'", happyResult)
	}

	// Edge cases
	// Test when input string is already equal to the size
	sizeEqualResult := PadEnd("test", 4)
	if sizeEqualResult != "test" {
		t.Errorf("PadEnd('test', 4) = %s; want 'test'", sizeEqualResult)
	}

	// Test when input string is larger than the size
	largerSizeResult := PadEnd("largerstring", 5)
	if largerSizeResult != "largerstring" {
		t.Errorf("PadEnd('largerstring', 5) = %s; want 'largerstring'", largerSizeResult)
	}

	// Test when using default value to pad
	defaultValueResult := PadEnd("default", 10, " value")
	if defaultValueResult != "default va" {
		t.Errorf("PadEnd('default', 10, ' value') = %s; want 'default va'", defaultValueResult)
	}

	// Test when input string is empty
	emptyInputResult := PadEnd("", 5)
	if emptyInputResult != "     " {
		t.Errorf("PadEnd('', 5) = %s; want '     '", emptyInputResult)
	}
}

func TestIsNumeric(t *testing.T) {
	// Happy path
	t.Run("Valid numeric string", func(t *testing.T) {
		result := IsNumeric("12345")
		if !result {
			t.Error(`IsNumeric("12345") = false; want true`)
		}
	})

	// Edge cases
	t.Run("Empty string", func(t *testing.T) {
		result := IsNumeric("")
		if result {
			t.Error(`IsNumeric("") = true; want false`)
		}
	})

	t.Run("Non-numeric string", func(t *testing.T) {
		result := IsNumeric("abc")
		if result {
			t.Error(`IsNumeric("abc") = true; want false`)
		}
	})
}
